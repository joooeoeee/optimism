package script

import (
	"encoding/binary"
	"fmt"
	"math/big"

	"github.com/holiman/uint256"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/tracing"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/params"

	"github.com/ethereum-optimism/optimism/op-chain-ops/foundry"
)

type CallFrame struct {
	Depth  int
	Opener vm.OpCode
	Sender common.Address
	Addr   common.Address
	Value  *big.Int
	Gas    uint64
}

type Host struct {
	log      log.Logger
	af       *foundry.ArtifactsFS
	chainCfg *params.ChainConfig
	env      *vm.EVM
	state    *state.StateDB
	stateDB  state.Database
	rawDB    ethdb.Database

	cheatcodes *Precompile[*CheatCodesPrecompile]
	console    *Precompile[*ConsolePrecompile]

	callStack []CallFrame
}

func NewHost(logger log.Logger, fs *foundry.ArtifactsFS, executionContext Context) *Host {
	h := &Host{
		log: logger,
		af:  fs,
	}

	h.chainCfg = &params.ChainConfig{
		ChainID: executionContext.chainID,
		// Ethereum forks in proof-of-work era.
		HomesteadBlock:      big.NewInt(0),
		EIP150Block:         big.NewInt(0),
		EIP155Block:         big.NewInt(0),
		EIP158Block:         big.NewInt(0),
		ByzantiumBlock:      big.NewInt(0),
		ConstantinopleBlock: big.NewInt(0),
		PetersburgBlock:     big.NewInt(0),
		IstanbulBlock:       big.NewInt(0),
		MuirGlacierBlock:    big.NewInt(0),
		BerlinBlock:         big.NewInt(0),
		LondonBlock:         big.NewInt(0),
		ArrowGlacierBlock:   big.NewInt(0),
		GrayGlacierBlock:    big.NewInt(0),
		MergeNetsplitBlock:  big.NewInt(0),
		// Ethereum forks in proof-of-stake era.
		TerminalTotalDifficulty:       big.NewInt(1),
		TerminalTotalDifficultyPassed: true,
		ShanghaiTime:                  new(uint64),
		CancunTime:                    new(uint64),
		PragueTime:                    nil,
		VerkleTime:                    nil,
		// OP-Stack forks are disabled, since we use this for L1.
		BedrockBlock: nil,
		RegolithTime: nil,
		CanyonTime:   nil,
		EcotoneTime:  nil,
		FjordTime:    nil,
		GraniteTime:  nil,
		InteropTime:  nil,
		Optimism:     nil,
	}

	h.rawDB = rawdb.NewMemoryDatabase()
	h.stateDB = state.NewDatabase(h.rawDB)
	var err error
	h.state, err = state.New(types.EmptyRootHash, h.stateDB, nil)
	if err != nil {
		panic(fmt.Errorf("failed to create memory state db: %w", err))
	}

	blockContext := vm.BlockContext{
		CanTransfer: core.CanTransfer,
		Transfer:    core.Transfer,
		GetHash: func(n uint64) (out common.Hash) {
			// mock a hash. // TODO: maybe warn/error, since we don't want scripts to be blockhash dependent?
			binary.BigEndian.PutUint64(out[:8], n)
			return crypto.Keccak256Hash(out[:])
		},
		L1CostFunc:  nil,
		Coinbase:    executionContext.feeRecipient,
		GasLimit:    executionContext.gasLimit,
		BlockNumber: new(big.Int).SetUint64(executionContext.blockNum),
		Time:        executionContext.timestamp,
		Difficulty:  nil, // not used anymore post-merge
		BaseFee:     big.NewInt(0),
		BlobBaseFee: big.NewInt(0),
		Random:      &executionContext.prevRandao,
	}

	txContext := vm.TxContext{
		Origin:       executionContext.origin,
		GasPrice:     big.NewInt(0),
		BlobHashes:   executionContext.blobHashes,
		BlobFeeCap:   big.NewInt(0),
		AccessEvents: state.NewAccessEvents(h.stateDB.PointCache()),
	}

	// TODO: attach to Host, and log each significant step
	trHooks := &tracing.Hooks{
		OnTxStart:         nil,
		OnTxEnd:           nil,
		OnEnter:           h.onEnter,
		OnExit:            h.onExit,
		OnOpcode:          nil,
		OnFault:           h.onFault,
		OnGasChange:       nil,
		OnBlockchainInit:  nil,
		OnClose:           nil,
		OnBlockStart:      nil,
		OnBlockEnd:        nil,
		OnSkippedBlock:    nil,
		OnGenesisBlock:    nil,
		OnSystemCallStart: nil,
		OnSystemCallEnd:   nil,
		OnBalanceChange:   nil,
		OnNonceChange:     nil,
		OnCodeChange:      nil,
		OnStorageChange:   nil,
		OnLog:             nil,
	}

	vmCfg := vm.Config{
		NoBaseFee: true,
		Tracer:    trHooks,
		// Override the precompiles, so we can insert things like the console, cheatcodes, and config contracts.
		PrecompileOverrides: h.getPrecompile,
	}

	h.env = vm.NewEVM(blockContext, txContext, h.state, h.chainCfg, vmCfg)

	return h
}

func (h *Host) EnableCheats() error {
	vmPrecompile, err := NewPrecompile[*CheatCodesPrecompile](&CheatCodesPrecompile{h: h})
	if err != nil {
		return fmt.Errorf("failed to init VM cheatcodes precompile: %w", err)
	}
	h.cheatcodes = vmPrecompile

	consolePrecompile, err := NewPrecompile[*ConsolePrecompile](&ConsolePrecompile{
		logger: h.log,
		sender: func() common.Address {
			return h.CurrentCall().Sender
		},
	})
	if err != nil {
		return fmt.Errorf("failed to init console precompile: %w", err)
	}
	h.console = consolePrecompile
	return nil
}

func (h *Host) prelude(from common.Address, to *common.Address) {
	rules := h.chainCfg.Rules(h.env.Context.BlockNumber, true, h.env.Context.Time)
	activePrecompiles := vm.ActivePrecompiles(rules)
	h.env.StateDB.Prepare(rules, from, h.env.Context.Coinbase, to, activePrecompiles, nil)
}

func (h *Host) Call(from common.Address, to common.Address, input []byte, gas uint64, value *uint256.Int) (returnData []byte, leftOverGas uint64, err error) {
	h.prelude(from, &to)
	return h.env.Call(vm.AccountRef(from), to, input, gas, value)
}

func (h *Host) LoadContract(artifactName, contractName string) (common.Address, error) {
	artifact, err := h.af.ReadArtifact(artifactName, contractName)
	if err != nil {
		return common.Address{}, fmt.Errorf("failed to load %s / %s: %w", artifactName, contractName, err)
	}
	deployedBytecode := artifact.DeployedBytecode.Object
	nonce := h.state.GetNonce(DefaultSenderAddr)
	addr := crypto.CreateAddress(DefaultSenderAddr, nonce+1)
	h.env.StateDB.SetCode(addr, deployedBytecode)
	h.state.SetNonce(DefaultSenderAddr, nonce+1)

	// TODO: srcmap.ParseSourceMap, register the sourcemap at this address, use for debugging
	return addr, nil
}

func (h *Host) getPrecompile(rules params.Rules, original vm.PrecompiledContract, addr common.Address) vm.PrecompiledContract {
	switch addr {
	case VMAddr:
		return h.cheatcodes // nil if cheats are not enabled
	case ConsoleAddr:
		return h.console // nil if cheats are not enabled
	// TODO: we can attach configurations this way, and directly provide reads for
	// TODO: we can override Artifacts.s.sol, to remember deployments
	default:
		return original
	}
}

func (h *Host) onEnter(depth int, typ byte, from common.Address, to common.Address, input []byte, gas uint64, value *big.Int) {
	h.callStack = append(h.callStack, CallFrame{
		Depth:  depth,
		Opener: vm.OpCode(typ),
		Sender: from,
		Addr:   to,
		Value:  value,
		Gas:    gas,
	})
}

func (h *Host) onExit(depth int, output []byte, gasUsed uint64, err error, reverted bool) {
	var addr common.Address
	// pop callStack (more if necessary)
	for len(h.callStack) > 0 && h.callStack[len(h.callStack)-1].Depth >= depth {
		addr = h.callStack[len(h.callStack)-1].Addr
		h.callStack = h.callStack[:len(h.callStack)-1]
	}
	if reverted {
		if msg, revertInspectErr := abi.UnpackRevert(output); revertInspectErr == nil {
			h.log.Warn("Revert", "addr", addr, "err", err, "revertMsg", msg)
		} else {
			h.log.Warn("Revert", "addr", addr, "err", err, "revertData", hexutil.Bytes(output))
		}
	}
}

func (h *Host) onFault(pc uint64, op byte, gas, cost uint64, scope tracing.OpContext, depth int, err error) {
	h.log.Warn("Fault", "addr", scope.Address(), "err", err)
}

func (h *Host) CurrentCall() CallFrame {
	if len(h.callStack) == 0 {
		return CallFrame{}
	}
	return h.callStack[len(h.callStack)-1]
}
