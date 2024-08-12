package script

import (
	"testing"

	"github.com/holiman/uint256"
	"github.com/stretchr/testify/require"

	"github.com/ethereum/go-ethereum/log"

	"github.com/ethereum-optimism/optimism/op-chain-ops/foundry"
	"github.com/ethereum-optimism/optimism/op-service/testlog"
)

func TestScript(t *testing.T) {
	logger := testlog.Logger(t, log.LevelInfo)
	af := foundry.OpenArtifactsDir("../../packages/contracts-bedrock/forge-artifacts")

	scriptContext := DefaultContext
	h := NewHost(logger, af, scriptContext)
	addr, err := h.LoadContract("Experiment.s", "Experiment")
	require.NoError(t, err)

	require.NoError(t, h.EnableCheats())

	input := bytes4("doThing()")
	returnData, _, err := h.Call(scriptContext.sender, addr, input[:], DefaultFoundryGasLimit, uint256.NewInt(0))
	require.NoError(t, err, "call failed: %x", string(returnData))
	t.Logf("call succeeded: %x", string(returnData))
}
