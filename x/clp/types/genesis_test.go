package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_NewGenesisState(t *testing.T) {
	genesisState := DefaultGenesisState()
	assert.Equal(t, genesisState.Params.MinCreatePoolThreshold, uint64(100))
}

func Test_DefaultGenesisState(t *testing.T) {
	genesisState := DefaultGenesisState()
	assert.Equal(t, genesisState.Params.MinCreatePoolThreshold, uint64(100))
	assert.Equal(t, genesisState.AddressWhitelist, []string{"cosmos1ny48eeuk4dm9f63dy0lwfgjhnvud9yvt8tcaat"})
}

func TestGenesisState_Validate(t *testing.T) {
	admin := GetDefaultCLPAdmin()

	tests := []struct {
		desc     string
		genState *GenesisState
		valid    bool
	}{
		{
			desc:     "default is valid",
			genState: DefaultGenesisState(),
			valid:    true,
		},
		{
			desc: "valid genesis state",
			genState: &GenesisState{
				Params:           DefaultParams(),
				AddressWhitelist: []string{admin.String()},
				RewardsBucketList: []RewardsBucket{
					{
						Denom: "0",
					},
					{
						Denom: "1",
					},
				},
			},
			valid: true,
		},
		{
			desc: "duplicated rewardsBucket",
			genState: &GenesisState{
				Params:           DefaultParams(),
				AddressWhitelist: []string{admin.String()},
				RewardsBucketList: []RewardsBucket{
					{
						Denom: "0",
					},
					{
						Denom: "0",
					},
				},
			},
			valid: false,
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			err := tc.genState.Validate()
			if tc.valid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}
