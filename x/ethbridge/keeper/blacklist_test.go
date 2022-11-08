package keeper_test

import (
	"testing"

	admintypes "github.com/Sifchain/sifnode/x/admin/types"

	"github.com/Sifchain/sifnode/x/ethbridge/test"
	"github.com/Sifchain/sifnode/x/ethbridge/types"
	oracletypes "github.com/Sifchain/sifnode/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestIsBlacklisted(t *testing.T) {
	tt := []struct {
		name      string
		addresses []string
		check     string
		expected  bool
	}{
		{"basic true",
			[]string{
				"0x782D10cC8c352D0524a1639eD261d29F47023922",
				"0x782D10cC8c352D0524a1639eD261d29F47023923",
			},
			"0x782D10cC8c352D0524a1639eD261d29F47023922",
			true,
		},
		{"basic false",
			[]string{
				"0x782D10cC8c352D0524a1639eD261d29F47023922",
				"0x782D10cC8c352D0524a1639eD261d29F47023923",
			},
			"0x782D10cC8c352D0524a1639eD261d29F47023924",
			false,
		},
		{"empty list", []string{}, "0x782D10cC8c352D0524a1639eD261d29F47023922", false},
		{"empty check", []string{}, "", false},
	}

	adminAddress, err := sdk.AccAddressFromBech32(types.TestAddress)
	require.NoError(t, err)
	for _, tc := range tt {
		tc := tc
		app, ctx := test.CreateTestApp(false)
		admin := admintypes.AdminAccount{
			AdminType:    admintypes.AdminType_ETHBRIDGE,
			AdminAddress: adminAddress.String(),
		}
		app.AdminKeeper.SetAdminAccount(ctx, &admin)
		err := app.EthbridgeKeeper.SetBlacklist(ctx, &types.MsgSetBlacklist{
			From:      adminAddress.String(),
			Addresses: tc.addresses,
		})
		require.NoError(t, err)
		got := app.EthbridgeKeeper.IsBlacklisted(ctx, tc.check)
		require.Equal(t, tc.expected, got)
	}
}

func TestSetBlacklist(t *testing.T) {
	tt := []struct {
		name        string
		addresses   []string
		updated     []string
		expectFalse []string
		expectTrue  []string
	}{
		{"replace all",
			[]string{
				"0x782D10cC8c352D0524a1639eD261d29F47023922",
				"0x782D10cC8c352D0524a1639eD261d29F47023923",
			},
			[]string{
				"0x782D10cC8c352D0524a1639eD261d29F47023924",
				"0x782D10cC8c352D0524a1639eD261d29F47023925",
			},
			[]string{"0x782D10cC8c352D0524a1639eD261d29F47023922", "0x782D10cC8c352D0524a1639eD261d29F47023923"},
			[]string{"0x782D10cC8c352D0524a1639eD261d29F47023924", "0x782D10cC8c352D0524a1639eD261d29F47023925"},
		},
		{"replace one",
			[]string{
				"0x782D10cC8c352D0524a1639eD261d29F47023922",
				"0x782D10cC8c352D0524a1639eD261d29F47023923",
			},
			[]string{
				"0x782D10cC8c352D0524a1639eD261d29F47023924",
				"0x782D10cC8c352D0524a1639eD261d29F47023922",
			},
			[]string{"0x782D10cC8c352D0524a1639eD261d29F47023923"},
			[]string{"0x782D10cC8c352D0524a1639eD261d29F47023924", "0x782D10cC8c352D0524a1639eD261d29F47023922"},
		},
		{"remove all",
			[]string{
				"0x782D10cC8c352D0524a1639eD261d29F47023922",
				"0x782D10cC8c352D0524a1639eD261d29F47023923",
			},
			[]string{},
			[]string{"0x782D10cC8c352D0524a1639eD261d29F47023922", "0x782D10cC8c352D0524a1639eD261d29F47023923"},
			[]string{},
		},
	}

	adminAddress, err := sdk.AccAddressFromBech32(types.TestAddress)
	require.NoError(t, err)

	for _, tc := range tt {
		tc := tc
		app, ctx := test.CreateTestApp(false)
		admin := admintypes.AdminAccount{
			AdminType:    admintypes.AdminType_ETHBRIDGE,
			AdminAddress: adminAddress.String(),
		}
		app.AdminKeeper.SetAdminAccount(ctx, &admin)
		err := app.EthbridgeKeeper.SetBlacklist(ctx, &types.MsgSetBlacklist{
			From:      adminAddress.String(),
			Addresses: tc.addresses,
		})
		require.NoError(t, err)
		err = app.EthbridgeKeeper.SetBlacklist(ctx, &types.MsgSetBlacklist{
			From:      adminAddress.String(),
			Addresses: tc.updated,
		})
		require.NoError(t, err)

		//list := app.EthbridgeKeeper.GetBlacklist(ctx)
		//fmt.Println(list)
		for _, address := range tc.expectTrue {
			require.True(t, app.EthbridgeKeeper.IsBlacklisted(ctx, address))
		}
		for _, address := range tc.expectFalse {
			require.False(t, app.EthbridgeKeeper.IsBlacklisted(ctx, address))
		}

	}
}

func TestKeeper_SetBlacklist_Nonadmin(t *testing.T) {
	app, ctx := test.CreateTestApp(false)
	err := app.EthbridgeKeeper.SetBlacklist(ctx, &types.MsgSetBlacklist{
		From: types.TestAddress,
	})
	require.ErrorIs(t, err, oracletypes.ErrNotAdminAccount)
}
