package oracle

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/Sifchain/sifnode/x/oracle/types"
)

func TestInitGenesis(t *testing.T) {
	addrs, valAddrs := CreateTestAddrs(2)

	tt := []struct {
		name    string
		genesis types.GenesisState
	}{
		{
			name:    "Default genesis",
			genesis: types.DefaultGenesisState(),
		},
		{
			name: "Nil genesis",
			genesis: GenesisState{
				AddressWhitelist: nil,
				AdminAddress:     nil,
				Prophecies:       nil,
			},
		},
		{
			name: "Prophecy",
			genesis: GenesisState{
				AddressWhitelist: valAddrs,
				AdminAddress:     addrs[0],
				Prophecies: []Prophecy{
					{
						ID: "asd",
						Status: types.Status{
							Text:       PendingStatusText,
							FinalClaim: "abc",
						},
						ClaimValidators: map[string][]sdk.ValAddress{
							"123": valAddrs,
						},
						ValidatorClaims: map[string]string{
							"321": "4321",
						},
					},
				},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ctx, keeper, _, _, _, _, _ := CreateTestKeepers(t, 1, []int64{1}, "")
			_ = InitGenesis(ctx, keeper, tc.genesis)

			if len(tc.genesis.AdminAddress) <= 0 {
				require.Nil(t, keeper.GetAdminAccount(ctx))
			} else {
				require.Equal(t, tc.genesis.AdminAddress, keeper.GetAdminAccount(ctx))
			}

			wl := keeper.GetOracleWhiteList(ctx)
			require.Equal(t, len(tc.genesis.AddressWhitelist), len(wl))
			for i, addr := range tc.genesis.AddressWhitelist {
				require.Equal(t, addr, wl[i])
			}

			prophecies := keeper.GetProphecies(ctx)
			require.Equal(t, len(tc.genesis.Prophecies), len(prophecies))
			for i, p := range tc.genesis.Prophecies {
				require.Equal(t, p, prophecies[i])
			}
		})
	}
}

func TestExportGenesis(t *testing.T) {
	addrs, valAddrs := CreateTestAddrs(2)

	tt := []struct {
		name    string
		genesis types.GenesisState
	}{
		{
			name:    "Default genesis",
			genesis: types.DefaultGenesisState(),
		},
		{
			name: "Nil genesis",
			genesis: GenesisState{
				AddressWhitelist: nil,
				AdminAddress:     nil,
				Prophecies:       nil,
			},
		},
		{
			name: "Prophecy",
			genesis: GenesisState{
				AddressWhitelist: valAddrs,
				AdminAddress:     addrs[0],
				Prophecies: []Prophecy{
					{
						ID: "asd",
						Status: types.Status{
							Text:       PendingStatusText,
							FinalClaim: "abc",
						},
						ClaimValidators: map[string][]sdk.ValAddress{
							"123": valAddrs,
						},
						ValidatorClaims: map[string]string{
							"321": "4321",
						},
					},
				},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ctx, keeper, _, _, _, _, _ := CreateTestKeepers(t, 1, []int64{1}, "")
			_ = InitGenesis(ctx, keeper, tc.genesis)
			genesis := ExportGenesis(ctx, keeper)

			if len(tc.genesis.AdminAddress) <= 0 {
				require.Nil(t, genesis.AdminAddress)
			} else {
				require.Equal(t, tc.genesis.AdminAddress, genesis.AdminAddress)
			}

			wl := genesis.AddressWhitelist
			require.Equal(t, len(tc.genesis.AddressWhitelist), len(wl))
			for i, addr := range tc.genesis.AddressWhitelist {
				require.Equal(t, addr, wl[i])
			}

			prophecies := genesis.Prophecies
			require.Equal(t, len(tc.genesis.Prophecies), len(prophecies))
			for i, p := range tc.genesis.Prophecies {
				require.Equal(t, p, prophecies[i])
			}
		})
	}
}
