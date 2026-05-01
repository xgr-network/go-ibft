package core

import (
	"math/big"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_CalculateQuorum(t *testing.T) {
	t.Parallel()

	vm := &ValidatorManager{
		vpLock: &sync.RWMutex{},
	}

	cases := []struct {
		validatorsVotingPower map[string]*big.Int
		signers               map[string]struct{}
		hasQuorum             bool
	}{
		{
			// case total voting power 4
			validatorsVotingPower: map[string]*big.Int{
				"A": big.NewInt(1),
				"B": big.NewInt(1),
				"C": big.NewInt(1),
				"D": big.NewInt(1),
			},
			// all 4 signed, has quorum (quorum is 3)
			signers: map[string]struct{}{
				"A": {},
				"B": {},
				"C": {},
				"D": {},
			},
			hasQuorum: true,
		},
		{
			// case total voting power 4
			validatorsVotingPower: map[string]*big.Int{
				"A": big.NewInt(1),
				"B": big.NewInt(1),
				"C": big.NewInt(1),
				"D": big.NewInt(1),
			},
			// only two signed (quorum is 3)
			signers: map[string]struct{}{
				"A": {},
				"B": {},
			},
			hasQuorum: false,
		},
		{
			// case total voting power of 6
			validatorsVotingPower: map[string]*big.Int{
				"A": big.NewInt(1),
				"B": big.NewInt(1),
				"C": big.NewInt(1),
				"D": big.NewInt(1),
				"E": big.NewInt(1),
				"F": big.NewInt(1),
			},
			// 4 signed (quorum is ceil(12/3)=4)
			signers: map[string]struct{}{
				"A": {},
				"B": {},
				"C": {},
				"D": {},
			},
			hasQuorum: true,
		},
		{
			// case total voting power of 6
			validatorsVotingPower: map[string]*big.Int{
				"A": big.NewInt(1),
				"B": big.NewInt(1),
				"C": big.NewInt(1),
				"D": big.NewInt(1),
				"E": big.NewInt(1),
				"F": big.NewInt(1),
			},
			// only 3 signed (quorum is 4)
			signers: map[string]struct{}{
				"A": {},
				"B": {},
				"C": {},
			},
			hasQuorum: false,
		},
		{
			// case total voting power of 9
			validatorsVotingPower: map[string]*big.Int{
				"A": big.NewInt(2),
				"B": big.NewInt(2),
				"C": big.NewInt(2),
				"D": big.NewInt(3),
			},
			// 3 signed with voting power of 6 (quorum is 6)
			signers: map[string]struct{}{
				"A": {},
				"C": {},
				"D": {},
			},
			hasQuorum: true,
		},
		{
			// case total voting power of 9
			validatorsVotingPower: map[string]*big.Int{
				"A": big.NewInt(2),
				"B": big.NewInt(2),
				"C": big.NewInt(2),
				"D": big.NewInt(3),
			},
			// only 2 signed with voting power of 5 (quorum is 6)
			signers: map[string]struct{}{
				"A": {},
				"D": {},
			},
			hasQuorum: false,
		},
		{
			// case total voting power of 10
			validatorsVotingPower: map[string]*big.Int{
				"A": big.NewInt(2),
				"B": big.NewInt(2),
				"C": big.NewInt(3),
				"D": big.NewInt(3),
			},
			// 3 signed with voting power of 7 (quorum is 7)
			signers: map[string]struct{}{
				"A": {},
				"B": {},
				"D": {},
			},
			hasQuorum: true,
		},
		{
			// case total voting power of 10
			validatorsVotingPower: map[string]*big.Int{
				"A": big.NewInt(2),
				"B": big.NewInt(2),
				"C": big.NewInt(3),
				"D": big.NewInt(3),
			},
			// only 2 signed with voting power of 5 (quorum is 7)
			signers: map[string]struct{}{
				"A": {},
				"D": {},
			},
			hasQuorum: false,
		},
		{
			// case total voting power of 21
			validatorsVotingPower: map[string]*big.Int{
				"A": big.NewInt(2),
				"B": big.NewInt(7),
				"C": big.NewInt(7),
				"D": big.NewInt(5),
			},
			// 3 signed with voting power of 16 (quorum is 14)
			signers: map[string]struct{}{
				"A": {},
				"B": {},
				"C": {},
			},
			hasQuorum: true,
		},
		{
			// case total voting power of 21
			validatorsVotingPower: map[string]*big.Int{
				"A": big.NewInt(2),
				"B": big.NewInt(7),
				"C": big.NewInt(7),
				"D": big.NewInt(5),
			},
			// only 2 signed with voting power of 12 (quorum is 14)
			signers: map[string]struct{}{
				"C": {},
				"D": {},
			},
			hasQuorum: false,
		},
		{
			// decay case: 4 validators, one decayed to 0, total voting power 30
			validatorsVotingPower: map[string]*big.Int{
				"A": big.NewInt(0),
				"B": big.NewInt(10),
				"C": big.NewInt(10),
				"D": big.NewInt(10),
			},
			// two live validators have exactly ceil(2/3*30)=20 power
			signers: map[string]struct{}{
				"B": {},
				"C": {},
			},
			hasQuorum: true,
		},
	}

	for _, c := range cases {
		require.NoError(t, vm.setCurrentVotingPower(c.validatorsVotingPower))
		require.Equal(t, c.hasQuorum, vm.HasQuorum(c.signers))
	}
}
