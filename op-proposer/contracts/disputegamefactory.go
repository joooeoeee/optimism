package contracts

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum-optimism/optimism/op-service/sources/batching"
	"github.com/ethereum-optimism/optimism/op-service/sources/batching/rpcblock"
	"github.com/ethereum-optimism/optimism/op-service/txmgr"
	"github.com/ethereum-optimism/optimism/packages/contracts-bedrock/snapshots"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

const (
	methodGameCount   = "gameCount"
	methodGameAtIndex = "gameAtIndex"
	methodInitBonds   = "initBonds"
	methodCreateGame  = "create"
	methodVersion     = "version"

	methodClaim = "claimData"
)

type gameMetadata struct {
	GameType  uint32
	Timestamp time.Time
	Proxy     common.Address
	Proposer  common.Address
}

type DisputeGameFactory struct {
	caller   *batching.MultiCaller
	contract *batching.BoundContract
	gameABI  *abi.ABI
}

func NewDisputeGameFactory(addr common.Address, caller *batching.MultiCaller) *DisputeGameFactory {
	factoryABI := snapshots.LoadDisputeGameFactoryABI()
	gameABI := snapshots.LoadFaultDisputeGameABI()
	return &DisputeGameFactory{
		caller:   caller,
		contract: batching.NewBoundContract(factoryABI, addr),
		gameABI:  gameABI,
	}
}

func (f *DisputeGameFactory) Version(ctx context.Context) (string, error) {
	result, err := f.caller.SingleCall(ctx, rpcblock.Latest, f.contract.Call(methodVersion))
	if err != nil {
		return "", fmt.Errorf("failed to get version: %w", err)
	}
	return result.GetString(0), nil
}

func (f *DisputeGameFactory) HasProposedSince(ctx context.Context, proposer common.Address, cutoff time.Time, gameType uint32) (uint64, bool, time.Time, error) {
	gameCount, err := f.gameCount(ctx)
	if err != nil {
		return 0, false, time.Time{}, fmt.Errorf("failed to get dispute game count: %w", err)
	}
	if gameCount == 0 {
		return 0, false, time.Time{}, nil
	}
	for idx := gameCount - 1; ; idx-- {
		game, err := f.gameAtIndex(ctx, idx)
		if err != nil {
			return 0, false, time.Time{}, fmt.Errorf("failed to get dispute game %d: %w", idx, err)
		}
		if game.Timestamp.Before(cutoff) {
			// Reached a game that is before the expected cutoff, so we haven't found a suitable proposal
			return gameCount, false, time.Time{}, nil
		}
		if game.GameType == gameType && game.Proposer == proposer {
			// Found a matching proposal
			return gameCount, true, game.Timestamp, nil
		}
		if idx == 0 { // Need to check here rather than in the for condition to avoid underflow
			// Checked every game and didn't find a match
			return gameCount, false, time.Time{}, nil
		}
	}
}

func (f *DisputeGameFactory) ProposalTx(ctx context.Context, gameType uint32, outputRoot common.Hash, l2BlockNum uint64) (txmgr.TxCandidate, error) {
	result, err := f.caller.SingleCall(ctx, rpcblock.Latest, f.contract.Call(methodInitBonds, gameType))
	if err != nil {
		return txmgr.TxCandidate{}, fmt.Errorf("failed to fetch init bond: %w", err)
	}
	initBond := result.GetBigInt(0)
	call := f.contract.Call(methodCreateGame, gameType, outputRoot, common.BigToHash(big.NewInt(int64(l2BlockNum))).Bytes())
	candidate, err := call.ToTxCandidate()
	if err != nil {
		return txmgr.TxCandidate{}, err
	}
	candidate.Value = initBond
	return candidate, err
}

func (f *DisputeGameFactory) gameCount(ctx context.Context) (uint64, error) {
	result, err := f.caller.SingleCall(ctx, rpcblock.Latest, f.contract.Call(methodGameCount))
	if err != nil {
		return 0, fmt.Errorf("failed to load game count: %w", err)
	}
	return result.GetBigInt(0).Uint64(), nil
}

func (f *DisputeGameFactory) gameAtIndex(ctx context.Context, idx uint64) (gameMetadata, error) {
	result, err := f.caller.SingleCall(ctx, rpcblock.Latest, f.contract.Call(methodGameAtIndex, new(big.Int).SetUint64(idx)))
	if err != nil {
		return gameMetadata{}, fmt.Errorf("failed to load game %v: %w", idx, err)
	}
	gameType := result.GetUint32(0)
	timestamp := result.GetUint64(1)
	proxy := result.GetAddress(2)

	gameContract := batching.NewBoundContract(f.gameABI, proxy)
	result, err = f.caller.SingleCall(ctx, rpcblock.Latest, gameContract.Call(methodClaim, big.NewInt(0)))
	if err != nil {
		return gameMetadata{}, fmt.Errorf("failed to load root claim of game %v: %w", idx, err)
	}
	// We don't need most of the claim data, only the claimant which is the game proposer
	claimant := result.GetAddress(2)

	return gameMetadata{
		GameType:  gameType,
		Timestamp: time.Unix(int64(timestamp), 0),
		Proxy:     proxy,
		Proposer:  claimant,
	}, nil
}
