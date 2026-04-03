package core

import (
	"math/big"
	"sync/atomic"
	"testing"

	"github.com/0xPolygon/go-ibft/messages/proto"
	"github.com/stretchr/testify/require"
)

func TestIBFT_RefreshVotingPowers(t *testing.T) {
	t.Run("updates live validator and quorum state for active height without restarting sequence", func(t *testing.T) {
		var (
			getVotingPowersCalls uint64
			roundStartsCalls     uint64
		)

		i := NewIBFT(
			mockLogger{},
			mockBackend{
				getVotingPowerFn: func(height uint64) (map[string]*big.Int, error) {
					atomic.AddUint64(&getVotingPowersCalls, 1)
					if atomic.LoadUint64(&getVotingPowersCalls) == 1 {
						return map[string]*big.Int{
							"A": big.NewInt(1),
							"B": big.NewInt(1),
							"C": big.NewInt(1),
						}, nil
					}

					return map[string]*big.Int{
						"A": big.NewInt(3),
						"B": big.NewInt(1),
					}, nil
				},
				roundStartsFn: func(_ *proto.View) error {
					atomic.AddUint64(&roundStartsCalls, 1)

					return nil
				},
			},
			mockTransport{},
		)

		const activeHeight uint64 = 77

		require.NoError(t, i.validatorManager.Init(activeHeight))
		require.Equal(t, "3", i.validatorManager.quorumSize.String())

		i.state.setView(&proto.View{Height: activeHeight, Round: 2})
		i.state.setRoundStarted(true)
		i.state.changeState(prepare)

		require.NoError(t, i.RefreshVotingPowers(activeHeight))

		require.Equal(t, uint64(2), atomic.LoadUint64(&getVotingPowersCalls))
		require.Equal(t, "3", i.validatorManager.quorumSize.String())
		require.Equal(t, big.NewInt(3), i.validatorManager.validatorsVotingPower["A"])
		require.Equal(t, uint64(0), atomic.LoadUint64(&roundStartsCalls))
	})

	t.Run("returns error for wrong height", func(t *testing.T) {
		i := NewIBFT(mockLogger{}, mockBackend{}, mockTransport{})

		i.state.setView(&proto.View{Height: 10, Round: 0})
		i.state.setRoundStarted(true)
		i.state.changeState(newRound)

		err := i.RefreshVotingPowers(11)
		require.ErrorIs(t, err, errActiveSequenceHeightMismatch)
	})
}
