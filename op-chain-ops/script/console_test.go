package script

import (
	"math/rand" // nosemgrep
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/log"

	"github.com/ethereum-optimism/optimism/op-service/testlog"
	"github.com/ethereum-optimism/optimism/op-service/testutils"
)

func TestConsole(t *testing.T) {
	logger, captLog := testlog.CaptureLogger(t, log.LevelDebug)

	rng := rand.New(rand.NewSource(123))
	sender := testutils.RandomAddress(rng)
	alice := testutils.RandomAddress(rng)
	bob := testutils.RandomAddress(rng)
	c := &ConsolePrecompile{
		logger: logger,
		sender: func() common.Address { return sender },
	}
	p, err := NewPrecompile[*ConsolePrecompile](c)
	require.NoError(t, err)

	// test Log_daf0d4aa
	input := make([]byte, 0, 4+32+32)
	input = append(input, hexutil.MustDecode("0xdaf0d4aa")...)
	input = append(input, pad32(alice[:])...)
	input = append(input, pad32(bob[:])...)
	_, err = p.Run(input)
	require.NoError(t, err)

	captLog.FindLog(testlog.NewMessageFilter("console"))
	captLog.FindLog(testlog.NewAttributesFilter("sender", sender.String()))
	captLog.FindLog(testlog.NewAttributesFilter("p0", alice.String()))
	captLog.FindLog(testlog.NewAttributesFilter("p1", bob.String()))
}
