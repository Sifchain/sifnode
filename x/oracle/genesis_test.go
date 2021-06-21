package oracle_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/Sifchain/sifnode/x/ethbridge/test"
	"github.com/Sifchain/sifnode/x/oracle"

	"github.com/Sifchain/sifnode/x/oracle/types"
)

func TestInitGenesis(t *testing.T) {
	tt, _ := testGenesisData(t)

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ctx, _, _, _, keeper, _, _ := test.CreateTestKeepers(t, 1, []int64{1}, "")
			_ = oracle.InitGenesis(ctx, keeper, tc.genesis)

			if len(tc.genesis.AdminAddress) <= 0 {
				require.Nil(t, keeper.GetAdminAccount(ctx))
			} else {
				require.Equal(t, tc.genesis.AdminAddress, keeper.GetAdminAccount(ctx).String())
			}

			wl := keeper.GetOracleWhiteList(ctx)
			require.Equal(t, len(tc.genesis.AddressWhitelist), len(wl))
			for i, addr := range tc.genesis.AddressWhitelist {
				require.Equal(t, addr, wl[i].String())
			}

			prophecies := keeper.GetProphecies(ctx)
			require.Equal(t, len(tc.genesis.Prophecies), len(prophecies))
			for i, p := range tc.genesis.Prophecies {
				serialised, err := prophecies[i].SerializeForDB()
				require.NoError(t, err)
				require.Equal(t, p, &serialised)
			}
		})
	}
}

func TestExportGenesis(t *testing.T) {
	tt, _ := testGenesisData(t)

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ctx, _, _, _, keeper, _, _ := test.CreateTestKeepers(t, 1, []int64{1}, "")
			_ = oracle.InitGenesis(ctx, keeper, tc.genesis)
			genesis := oracle.ExportGenesis(ctx, keeper)

			require.Equal(t, tc.genesis.AdminAddress, genesis.AdminAddress)

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
	tt, prophecies := testGenesisData(t)

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ctx, _, _, _, keeper, encCfg, _ := test.CreateTestKeepers(t, 1, []int64{1}, "")
			_ = oracle.InitGenesis(ctx, keeper, tc.genesis)
			genesis := oracle.ExportGenesis(ctx, keeper)

			genesisData := encCfg.Marshaler.MustMarshalJSON(genesis)

			var genesisState types.GenesisState
			encCfg.Marshaler.MustUnmarshalJSON(genesisData, &genesisState)

			ctx, _, _, _, keeper, _, _ = test.CreateTestKeepers(t, 1, []int64{1}, "")
			_ = oracle.InitGenesis(ctx, keeper, genesisState)

			require.Equal(t, tc.genesis.AdminAddress, genesis.AdminAddress)

			wl := genesis.AddressWhitelist
			require.Equal(t, len(tc.genesis.AddressWhitelist), len(wl))
			for i, addr := range tc.genesis.AddressWhitelist {
				require.Equal(t, addr, wl[i])
			}

			dbProphecies := genesis.Prophecies
			require.Equal(t, len(tc.genesis.Prophecies), len(dbProphecies))
			for i, p := range tc.genesis.Prophecies {
				require.Equal(t, p, dbProphecies[i])
				unserialised, err := p.DeserializeFromDB()
				require.NoError(t, err)
				require.Equal(t, prophecies[i], unserialised)
			}
		})
	}
}

type testCase struct {
	name    string
	genesis types.GenesisState
}

func testGenesisData(t *testing.T) ([]testCase, []types.Prophecy) {
	addrs, valAddrs := test.CreateTestAddrs(2)

	var whitelist []string
	for _, addr := range valAddrs {
		whitelist = append(whitelist, addr.String())
	}

	prophecy := types.Prophecy{
		ID: "asd",
		Status: types.Status{
			Text:       types.StatusText_STATUS_TEXT_PENDING,
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
			genesis: *types.DefaultGenesisState(),
		},
		{
			name:    "Nil genesis",
			genesis: types.GenesisState{},
		},
		{
			name: "Prophecy",
			genesis: types.GenesisState{
				AddressWhitelist: whitelist,
				AdminAddress:     addrs[0].String(),
				Prophecies: []*types.DBProphecy{
					&dbProphecy,
				},
			},
		},
	}, []types.Prophecy{prophecy}
}
