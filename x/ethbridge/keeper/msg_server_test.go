package keeper_test

import (
	keeper2 "github.com/Sifchain/sifnode/x/ethbridge/keeper"
	"github.com/Sifchain/sifnode/x/ethbridge/test"
	"github.com/Sifchain/sifnode/x/ethbridge/types"
	types2 "github.com/Sifchain/sifnode/x/tokenregistry/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMsgServer_Lock(t *testing.T) {
	ctx, app := test.CreateSimulatorApp(false)
	addresses, _ := test.CreateTestAddrs(2)
	admin := addresses[0]
	nonAdmin := addresses[1]
	msg := types.NewMsgLock(1, admin, ethereumSender, amount, "stake", amount)
	coins := sdk.NewCoins(sdk.NewCoin("stake", amount), sdk.NewCoin(types.CethSymbol, amount))
	_ = app.BankKeeper.MintCoins(ctx, types.ModuleName, coins)
	_ = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, admin, coins)
	app.TokenRegistryKeeper.SetAdminAccount(ctx, &types2.AdminAccount{
		AdminType:    types2.AdminType_ETHBRIDGE,
		AdminAddress: admin.String(),
	})
	msgPauseNonAdmin := types.MsgPauser{
		Signer:   nonAdmin.String(),
		IsPaused: true,
	}
	msgPause := types.MsgPauser{
		Signer:   admin.String(),
		IsPaused: true,
	}
	msgUnPause := types.MsgPauser{
		Signer:   admin.String(),
		IsPaused: false,
	}
	msgServer := keeper2.NewMsgServerImpl(app.EthbridgeKeeper)
	// Pause with Non Admin Account
	_, err := msgServer.SetPauser(sdk.WrapSDKContext(ctx), &msgPauseNonAdmin)
	require.Error(t, err)

	// Pause Transactions
	_, err = msgServer.SetPauser(sdk.WrapSDKContext(ctx), &msgPause)
	require.NoError(t, err)

	// Fail Lock
	_, err = msgServer.Lock(sdk.WrapSDKContext(ctx), &msg)
	require.Error(t, err)

	// Unpause Transactions
	_, err = msgServer.SetPauser(sdk.WrapSDKContext(ctx), &msgUnPause)
	require.NoError(t, err)

	// Lock Success
	_, err = msgServer.Lock(sdk.WrapSDKContext(ctx), &msg)
	require.NoError(t, err)
}

func TestMsgServer_Burn(t *testing.T) {
	ctx, app := test.CreateSimulatorApp(false)
	addresses, _ := test.CreateTestAddrs(1)
	admin := addresses[0]
	coins := sdk.NewCoins(sdk.NewCoin("stake", amount), sdk.NewCoin(types.CethSymbol, amount))
	_ = app.BankKeeper.MintCoins(ctx, types.ModuleName, coins)
	_ = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, admin, coins)
	app.TokenRegistryKeeper.SetAdminAccount(ctx, &types2.AdminAccount{
		AdminType:    types2.AdminType_ETHBRIDGE,
		AdminAddress: admin.String(),
	})
	app.EthbridgeKeeper.AddPeggyToken(ctx, "stake")
	msg := types.NewMsgBurn(1, admin, ethereumSender, amount, "stake", amount)
	msgPause := types.MsgPauser{
		Signer:   admin.String(),
		IsPaused: true,
	}
	msgUnPause := types.MsgPauser{
		Signer:   admin.String(),
		IsPaused: false,
	}
	msgServer := keeper2.NewMsgServerImpl(app.EthbridgeKeeper)

	// Pause Transactions
	_, err := msgServer.SetPauser(sdk.WrapSDKContext(ctx), &msgPause)
	require.NoError(t, err)

	// Fail Burn
	_, err = msgServer.Burn(sdk.WrapSDKContext(ctx), &msg)
	require.Error(t, err)

	// Unpause Transactions
	_, err = msgServer.SetPauser(sdk.WrapSDKContext(ctx), &msgUnPause)
	require.NoError(t, err)

	// Burn Success
	_, err = msgServer.Burn(sdk.WrapSDKContext(ctx), &msg)
	require.NoError(t, err)

}
