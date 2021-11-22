package keeper_test

import (
	"testing"

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
		var ctx, keeper, _, _, oracleKeeper, _, _ = test.CreateTestKeepers(t, 0.7, []int64{3, 3}, "")
		oracleKeeper.SetAdminAccount(ctx, adminAddress)
		err := keeper.SetBlacklist(ctx, &types.MsgSetBlacklist{
			From:      adminAddress.String(),
			Addresses: tc.addresses,
		})
		require.NoError(t, err)
		got := keeper.IsBlacklisted(ctx, tc.check)
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
		var ctx, keeper, _, _, oracleKeeper, _, _ = test.CreateTestKeepers(t, 0.7, []int64{3, 3}, "")
		oracleKeeper.SetAdminAccount(ctx, adminAddress)
		err := keeper.SetBlacklist(ctx, &types.MsgSetBlacklist{
			From:      adminAddress.String(),
			Addresses: tc.addresses,
		})
		require.NoError(t, err)
		err = keeper.SetBlacklist(ctx, &types.MsgSetBlacklist{
			From:      adminAddress.String(),
			Addresses: tc.updated,
		})
		require.NoError(t, err)
		for _, address := range tc.expectTrue {
			require.True(t, keeper.IsBlacklisted(ctx, address))
		}
		for _, address := range tc.expectFalse {
			require.False(t, keeper.IsBlacklisted(ctx, address))
		}
	}
}

func TestKeeper_SetBlacklist_Nonadmin(t *testing.T) {
	var ctx, keeper, _, _, _, _, _ = test.CreateTestKeepers(t, 0.7, []int64{3, 3}, "")
	err := keeper.SetBlacklist(ctx, &types.MsgSetBlacklist{
		From: types.TestAddress,
	})
	require.ErrorIs(t, err, oracletypes.ErrNotAdminAccount)
}
