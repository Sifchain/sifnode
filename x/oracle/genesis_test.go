package oracle

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/Sifchain/sifnode/x/oracle/types"
)

func TestInitGenesis(t *testing.T) {
	tt := getTestGenesisCases(t)

	for c := range tt {
		tc := tt[c]
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
				dbP, err := prophecies[i].SerializeForDB()
				require.NoError(t, err)
				require.Equal(t, p, dbP)
			}
		})
	}
}

func TestExportGenesis(t *testing.T) {
	tt := getTestGenesisCases(t)

	for c := range tt {
		tc := tt[c]
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

func TestGenesisMarshalling(t *testing.T) {
	tt := getTestGenesisCases(t)

	for c := range tt {
		tc := tt[c]
		t.Run(tc.name, func(t *testing.T) {
			ctx, keeper, _, _, _, _, _ := CreateTestKeepers(t, 1, []int64{1}, "")
			_ = InitGenesis(ctx, keeper, tc.genesis)
			genesis := ExportGenesis(ctx, keeper)

			genesisData := keeper.Cdc.MustMarshalJSON(genesis)

			var genesisState GenesisState
			keeper.Cdc.MustUnmarshalJSON(genesisData, &genesisState)

			ctx, keeper, _, _, _, _, _ = CreateTestKeepers(t, 1, []int64{1}, "")
			_ = InitGenesis(ctx, keeper, genesisState)

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

type testCase struct {
	name    string
	genesis types.GenesisState
}

func getTestGenesisCases(t *testing.T) []testCase {
	addrs, valAddrs := CreateTestAddrs(2)

	prophecy := Prophecy{
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
	}

	dbProphecy, err := prophecy.SerializeForDB()
	require.NoError(t, err)

	return []testCase{
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
				Prophecies:       []DBProphecy{dbProphecy},
			},
		},
	}
}
