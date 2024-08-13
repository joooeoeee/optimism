package script

import (
	"github.com/ethereum/go-ethereum/core/tracing"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/holiman/uint256"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// Warp implements https://book.getfoundry.sh/cheatcodes/warp
func (c *CheatCodesPrecompile) Warp(timestamp *big.Int) {
	c.h.env.Context.Time = timestamp.Uint64()
}

// Roll implements https://book.getfoundry.sh/cheatcodes/roll
func (c *CheatCodesPrecompile) Roll(num *big.Int) {
	c.h.env.Context.BlockNumber = num
}

// Fee implements https://book.getfoundry.sh/cheatcodes/fee
func (c *CheatCodesPrecompile) Fee(fee *big.Int) {
	c.h.env.Context.BaseFee = fee
}

// GetBlockTimestamp implements https://book.getfoundry.sh/cheatcodes/get-block-timestamp
func (c *CheatCodesPrecompile) GetBlockTimestamp() *big.Int {
	return new(big.Int).SetUint64(c.h.env.Context.Time)
}

// GetBlockNumber implements https://book.getfoundry.sh/cheatcodes/get-block-number
func (c *CheatCodesPrecompile) GetBlockNumber() *big.Int {
	return c.h.env.Context.BlockNumber
}

// Difficulty implements https://book.getfoundry.sh/cheatcodes/difficulty
func (c *CheatCodesPrecompile) Difficulty(_ *big.Int) error {
	return vm.ErrExecutionReverted // only post-merge is supported
}

// Prevrandao implements https://book.getfoundry.sh/cheatcodes/prevrandao
func (c *CheatCodesPrecompile) Prevrandao(v [32]byte) {
	c.h.env.Context.Random = (*common.Hash)(&v)
}

// ChainId implements https://book.getfoundry.sh/cheatcodes/chain-id
func (c *CheatCodesPrecompile) ChainId(id *big.Int) {
	c.h.env.ChainConfig().ChainID = id
	c.h.chainCfg.ChainID = id
	// c.h.env.rules.ChainID is unused, but should maybe also be modified
}

// Store implements https://book.getfoundry.sh/cheatcodes/store
func (c *CheatCodesPrecompile) Store(account common.Address, slot [32]byte, value [32]byte) {
	c.h.state.SetState(account, slot, value)
}

// Load implements https://book.getfoundry.sh/cheatcodes/load
func (c *CheatCodesPrecompile) Load(account common.Address, slot [32]byte) [32]byte {
	return c.h.state.GetState(account, slot)
}

// Etch implements https://book.getfoundry.sh/cheatcodes/etch
func (c *CheatCodesPrecompile) Etch(who common.Address, code []byte) {
	c.h.state.SetCode(who, code)
}

// Deal implements https://book.getfoundry.sh/cheatcodes/deal
func (c *CheatCodesPrecompile) Deal(who common.Address, newBalance *big.Int) {
	c.h.state.SetBalance(who, uint256.MustFromBig(newBalance), tracing.BalanceChangeUnspecified)
}

// Prank_ca669fa7 implements https://book.getfoundry.sh/cheatcodes/prank
func (c *CheatCodesPrecompile) Prank_ca669fa7(sender common.Address) {
	// TODO
}

// Prank_47e50cce implements https://book.getfoundry.sh/cheatcodes/prank
func (c *CheatCodesPrecompile) Prank_47e50cce(sender common.Address, origin common.Address) {
	c.Prank_ca669fa7(sender)
	c.h.env.Origin = origin
}

// StartPrank_06447d56 implements https://book.getfoundry.sh/cheatcodes/start-prank
func (c *CheatCodesPrecompile) StartPrank_06447d56(sender common.Address) {
	// TODO
}

// StartPrank_45b56078 implements https://book.getfoundry.sh/cheatcodes/start-prank
func (c *CheatCodesPrecompile) StartPrank_45b56078(sender common.Address, origin common.Address) {
	// TODO
}

// StopPrank implements https://book.getfoundry.sh/cheatcodes/stop-prank
func (c *CheatCodesPrecompile) StopPrank() {
	// TODO
}

type CallerMode uint64 // TODO

// ReadCallers implements https://book.getfoundry.sh/cheatcodes/read-callers
func (c *CheatCodesPrecompile) ReadCallers() (callerMode CallerMode, msgSender common.Address, txOrigin common.Address) {
	return 0, c.h.CurrentCall().Sender, c.h.env.TxContext.Origin
}

// Record implements https://book.getfoundry.sh/cheatcodes/record
func (c *CheatCodesPrecompile) Record() {

}

// Accesses implements https://book.getfoundry.sh/cheatcodes/accesses
func (c *CheatCodesPrecompile) Accesses() (reads [][32]byte, writes [][32]byte) {
	// TODO
	return nil, nil
}

// RecordLogs implements https://book.getfoundry.sh/cheatcodes/record-logs
func (c *CheatCodesPrecompile) RecordLogs() {

}

type Log struct {
	Topics  [][32]byte
	Data    []byte
	Emitter common.Address
}

// GetRecordedLogs implements https://book.getfoundry.sh/cheatcodes/get-recorded-logs
func (c *CheatCodesPrecompile) GetRecordedLogs() []Log {
	return nil // TODO
}

// SetNonce implements https://book.getfoundry.sh/cheatcodes/set-nonce
func (c *CheatCodesPrecompile) SetNonce(account common.Address, nonce uint64) {
	c.h.state.SetNonce(account, nonce)
}

// GetNonce implements https://book.getfoundry.sh/cheatcodes/get-nonce
func (c *CheatCodesPrecompile) GetNonce(addr common.Address) uint64 {
	return c.h.state.GetNonce(addr)
}

// MockCall_b96213e4 implements https://book.getfoundry.sh/cheatcodes/mock-call
func (c *CheatCodesPrecompile) MockCall_b96213e4(where common.Address, data []byte, retdata []byte) error {
	return vm.ErrExecutionReverted // TODO
}

// MockCall_81409b91 implements https://book.getfoundry.sh/cheatcodes/mock-call
func (c *CheatCodesPrecompile) MockCall_81409b91(where common.Address, value *big.Int, data []byte, retdata []byte) error {
	return vm.ErrExecutionReverted // TODO
}

// MockCallRevert_dbaad147 implements https://book.getfoundry.sh/cheatcodes/mock-call-revert
func (c *CheatCodesPrecompile) MockCallRevert_dbaad147(where common.Address, data []byte, retdata []byte) error {
	return vm.ErrExecutionReverted // TODO
}

// MockCallRevert_d23cd037 implements https://book.getfoundry.sh/cheatcodes/mock-call-revert
func (c *CheatCodesPrecompile) MockCallRevert_d23cd037(where common.Address, value *big.Int, data []byte, retdata []byte) error {
	return vm.ErrExecutionReverted // TODO
}

// ClearMockedCalls implements https://book.getfoundry.sh/cheatcodes/clear-mocked-calls
func (c *CheatCodesPrecompile) ClearMockedCalls() {
	// TODO
}

// Coinbase implements https://book.getfoundry.sh/cheatcodes/coinbase
func (c *CheatCodesPrecompile) Coinbase(addr common.Address) {
	c.h.env.Context.Coinbase = addr
}

// Broadcast_afc98040 implements https://book.getfoundry.sh/cheatcodes/broadcast
func (c *CheatCodesPrecompile) Broadcast_afc98040() {
	c.h.log.Info("broadcasting next call")
	// TODO
}

// Broadcast_e6962cdb implements https://book.getfoundry.sh/cheatcodes/broadcast
func (c *CheatCodesPrecompile) Broadcast_e6962cdb(who common.Address) {
	c.h.log.Info("broadcasting next call", "who", who)
	// TODO
}

// StartBroadcast_7fb5297f implements https://book.getfoundry.sh/cheatcodes/start-broadcast
func (c *CheatCodesPrecompile) StartBroadcast_7fb5297f() {
	c.h.log.Info("starting repeat-broadcast")
	// TODO
}

// StartBroadcast_7fec2a8d implements https://book.getfoundry.sh/cheatcodes/start-broadcast
func (c *CheatCodesPrecompile) StartBroadcast_7fec2a8d(who common.Address) {
	c.h.log.Info("starting repeat-broadcast", "who", who)
	// TODO
}

// StopBroadcast implements https://book.getfoundry.sh/cheatcodes/stop-broadcast
func (c *CheatCodesPrecompile) StopBroadcast() {
	c.h.log.Info("stopping repeat-broadcast")
	// TODO
}

// PauseGasMetering implements https://book.getfoundry.sh/cheatcodes/pause-gas-metering
func (c *CheatCodesPrecompile) PauseGasMetering() error {
	return vm.ErrExecutionReverted
}

// ResumeGasMetering implements https://book.getfoundry.sh/cheatcodes/resume-gas-metering
func (c *CheatCodesPrecompile) ResumeGasMetering() {
	// no-op, since PauseGasMetering is unsupported
}

// TxGasPrice implements https://book.getfoundry.sh/cheatcodes/tx-gas-price
func (c *CheatCodesPrecompile) TxGasPrice(newGasPrice *big.Int) {
	c.h.env.TxContext.GasPrice = newGasPrice
}

// StartStateDiffRecording implements https://book.getfoundry.sh/cheatcodes/start-state-diff-recording
func (c *CheatCodesPrecompile) StartStateDiffRecording() error {
	return vm.ErrExecutionReverted // TODO not supported yet
}

// StopAndReturnStateDiff implements https://book.getfoundry.sh/cheatcodes/stop-and-return-state-diff
func (c *CheatCodesPrecompile) StopAndReturnStateDiff() error {
	return vm.ErrExecutionReverted // TODO not supported yet
}
