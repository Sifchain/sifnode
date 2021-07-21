package oracle_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Sifchain/sifnode/x/ethbridge/test"
	"github.com/Sifchain/sifnode/x/oracle"

	"github.com/Sifchain/sifnode/x/oracle/types"
)

//nolint:golint
func TestInitGenesis(t *testing.T) {
	networkDescriptor := types.NewNetworkIdentity(types.NetworkDescriptor_NETWORK_DESCRIPTOR_ETHEREUM)

	tt, _ := testGenesisData(t)

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ctx, _, _, _, keeper, _, _, _ := test.CreateTestKeepers(t, 1, []int64{1}, "")
			keeper.RemoveOracleWhiteList(ctx, networkDescriptor)

			_ = oracle.InitGenesis(ctx, keeper, tc.genesis)

			if len(tc.genesis.AdminAddress) <= 0 {
				require.Nil(t, keeper.GetAdminAccount(ctx))
			} else {
				require.Equal(t, tc.genesis.AdminAddress, keeper.GetAdminAccount(ctx).String())
			}

			wl := keeper.GetOracleWhiteList(ctx, networkDescriptor).WhiteList

			whiteList, ok := tc.genesis.AddressWhitelist[uint32(types.NetworkDescriptor_NETWORK_DESCRIPTOR_ETHEREUM)]

			if ok {
				for addr := range whiteList.WhiteList {
					_, ok := wl[addr]
					require.Equal(t, ok, true)
				}
			}

			prophecies := keeper.GetProphecies(ctx)
			require.Equal(t, len(tc.genesis.Prophecies), len(prophecies))
			for i, p := range tc.genesis.Prophecies {
				require.Equal(t, *p, prophecies[i])
			}
		})
	}
}

func TestExportGenesis(t *testing.T) {
	tt, _ := testGenesisData(t)

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ctx, _, _, _, keeper, _, _, _ := test.CreateTestKeepers(t, 1, []int64{1}, "")
			networkDescriptor := types.NetworkIdentity{NetworkDescriptor: types.NetworkDescriptor_NETWORK_DESCRIPTOR_ETHEREUM}

			keeper.RemoveOracleWhiteList(ctx, networkDescriptor)

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
			ctx, _, _, _, keeper, encCfg, _, _ := test.CreateTestKeepers(t, 1, []int64{1}, "")
			networkDescriptor := types.NetworkIdentity{NetworkDescriptor: types.NetworkDescriptor_NETWORK_DESCRIPTOR_ETHEREUM}
			keeper.RemoveOracleWhiteList(ctx, networkDescriptor)

			_ = oracle.InitGenesis(ctx, keeper, tc.genesis)
			genesis := oracle.ExportGenesis(ctx, keeper)

			genesisData := encCfg.Marshaler.MustMarshalJSON(genesis)

			var genesisState types.GenesisState
			encCfg.Marshaler.MustUnmarshalJSON(genesisData, &genesisState)

			ctx, _, _, _, keeper, _, _, _ = test.CreateTestKeepers(t, 1, []int64{1}, "")

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
				require.Equal(t, prophecies[i], *p)
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
	power := uint32(100)

	whiteList := types.ValidatorWhiteList{WhiteList: make(map[string]uint32)}
	whiteList.WhiteList[valAddrs[0].String()] = power
	whiteList.WhiteList[valAddrs[1].String()] = power

	prophecy := types.Prophecy{
		Id:              []byte("asd"),
		Status:          types.StatusText_STATUS_TEXT_PENDING,
		ClaimValidators: []string{valAddrs[0].String()},
	}

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
				AddressWhitelist: map[uint32]*types.ValidatorWhiteList{uint32(types.NetworkDescriptor_NETWORK_DESCRIPTOR_ETHEREUM): &whiteList},
				AdminAddress:     addrs[0].String(),
				Prophecies: []*types.Prophecy{
					&prophecy,
				},
			},
		},
	}, []types.Prophecy{prophecy}
}
