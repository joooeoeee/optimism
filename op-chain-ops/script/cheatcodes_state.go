package script

import (
	"github.com/ethereum-optimism/optimism/op-chain-ops/foundry"
	"github.com/ethereum/go-ethereum/core/state"
)

func (c *CheatCodesPrecompile) LoadAllocs(pathToAllocsJson string) {
	c.h.log.Info("loading state", "target", pathToAllocsJson)
}

func (c *CheatCodesPrecompile) DumpState(pathToStateJson string) {
	c.h.log.Info("dumping state", "target", pathToStateJson)
	var allocs foundry.ForgeAllocs
	c.h.state.DumpToCollector(&allocs, &state.DumpConfig{
		OnlyWithAddresses: true,
	})
	_ = allocs
	// TODO
}
