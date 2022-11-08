package keeper_test

import (
	"testing"

	adminTypes "github.com/Sifchain/sifnode/x/admin/types"
	ethbriddgeKeeper "github.com/Sifchain/sifnode/x/ethbridge/keeper"
	"github.com/Sifchain/sifnode/x/ethbridge/test"
	"github.com/Sifchain/sifnode/x/ethbridge/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestMsgServer_Lock_No_Pause_Set(t *testing.T) {
	ctx, app := test.CreateSimulatorApp(false)
	addresses, _ := test.CreateTestAddrs(2)
	admin := addresses[0]
	// nonAdmin := addresses[1]
	msg := types.NewMsgLock(1, admin, ethereumSender, amount, "stake", amount)
	coins := sdk.NewCoins(sdk.NewCoin("stake", amount), sdk.NewCoin(types.CethSymbol, amount))
	_ = app.BankKeeper.MintCoins(ctx, types.ModuleName, coins)
	_ = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, admin, coins)
	app.AdminKeeper.SetAdminAccount(ctx, &adminTypes.AdminAccount{
		AdminType:    adminTypes.AdminType_ETHBRIDGE,
		AdminAddress: admin.String(),
	})
	msgServer := ethbriddgeKeeper.NewMsgServerImpl(app.EthbridgeKeeper)

	_, err := msgServer.Lock(sdk.WrapSDKContext(ctx), &msg)
	require.NoError(t, err)
}

func TestMsgServer_Lock(t *testing.T) {
	ctx, app := test.CreateSimulatorApp(false)
	addresses, _ := test.CreateTestAddrs(2)
	admin := addresses[0]
	nonAdmin := addresses[1]
	msg := types.NewMsgLock(1, admin, ethereumSender, amount, "stake", amount)
	coins := sdk.NewCoins(sdk.NewCoin("stake", amount), sdk.NewCoin(types.CethSymbol, amount))
	_ = app.BankKeeper.MintCoins(ctx, types.ModuleName, coins)
	_ = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, admin, coins)
	app.AdminKeeper.SetAdminAccount(ctx, &adminTypes.AdminAccount{
		AdminType:    adminTypes.AdminType_ETHBRIDGE,
		AdminAddress: admin.String(),
	})
	msgPauseNonAdmin := types.MsgPause{
		Signer:   nonAdmin.String(),
		IsPaused: true,
	}
	msgPause := types.MsgPause{
		Signer:   admin.String(),
		IsPaused: true,
	}
	msgUnPause := types.MsgPause{
		Signer:   admin.String(),
		IsPaused: false,
	}
	msgServer := ethbriddgeKeeper.NewMsgServerImpl(app.EthbridgeKeeper)
	// Pause with Non Admin Account
	_, err := msgServer.SetPause(sdk.WrapSDKContext(ctx), &msgPauseNonAdmin)
	require.Error(t, err)

	// Pause Transactions
	_, err = msgServer.SetPause(sdk.WrapSDKContext(ctx), &msgPause)
	require.NoError(t, err)

	// Fail Lock
	_, err = msgServer.Lock(sdk.WrapSDKContext(ctx), &msg)
	require.Error(t, err)

	// Unpause Transactions
	_, err = msgServer.SetPause(sdk.WrapSDKContext(ctx), &msgUnPause)
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
	app.AdminKeeper.SetAdminAccount(ctx, &adminTypes.AdminAccount{
		AdminType:    adminTypes.AdminType_ETHBRIDGE,
		AdminAddress: admin.String(),
	})
	app.EthbridgeKeeper.AddPeggyToken(ctx, "stake")
	msg := types.NewMsgBurn(1, admin, ethereumSender, amount, "stake", amount)
	msgPause := types.MsgPause{
		Signer:   admin.String(),
		IsPaused: true,
	}
	msgUnPause := types.MsgPause{
		Signer:   admin.String(),
		IsPaused: false,
	}
	msgServer := ethbriddgeKeeper.NewMsgServerImpl(app.EthbridgeKeeper)

	// Pause Transactions
	_, err := msgServer.SetPause(sdk.WrapSDKContext(ctx), &msgPause)
	require.NoError(t, err)

	// Fail Burn
	_, err = msgServer.Burn(sdk.WrapSDKContext(ctx), &msg)
	require.Error(t, err)

	// Unpause Transactions
	_, err = msgServer.SetPause(sdk.WrapSDKContext(ctx), &msgUnPause)
	require.NoError(t, err)

	// Burn Success
	_, err = msgServer.Burn(sdk.WrapSDKContext(ctx), &msg)
	require.NoError(t, err)

}
