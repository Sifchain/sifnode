//go:build !FEATURE_TOGGLE_MARGIN_CLI_ALPHA
// +build !FEATURE_TOGGLE_MARGIN_CLI_ALPHA

package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/Sifchain/sifnode/x/clp/types"
	tokenregistrytypes "github.com/Sifchain/sifnode/x/tokenregistry/types"
)

func FEATURE_TOGGLE_MARGIN_CLI_ALPHA_ProcessRemovelQueue(ctx sdk.Context, k msgServer, msg *types.MsgAddLiquidity, newPoolUnits sdk.Uint) {
}

//  ensure requested removal amount is less than available - what is already on the queue
func FEATURE_TOGGLE_MARGIN_CLI_ALPHA_VerifyEnoughWithdrawUnitsAvailableForLP(ctx sdk.Context, k msgServer, msg *types.MsgRemoveLiquidityUnits, lp types.LiquidityProvider) error {
	if msg.WithdrawUnits.GT(lp.LiquidityProviderUnits) {
		return sdkerrors.Wrap(types.ErrUnableToRemoveLiquidity, fmt.Sprintf("WithdrawUnits %s greater than total LP units %s", msg.WithdrawUnits, lp.LiquidityProviderUnits))
	}
	return nil
}

//  ensure requested removal amount is less than available - what is already on the queue
func FEATURE_TOGGLE_MARGIN_CLI_ALPHA_VerifyEnoughWBasisPointsAvailableForLP(ctx sdk.Context, k msgServer, msg *types.MsgRemoveLiquidity, lp types.LiquidityProvider) error {
	return nil
}

func FEATURE_TOGGLE_MARGIN_CLI_ALPHA_QueueRemovalWithWithdrawUnits(ctx sdk.Context, k msgServer, msg *types.MsgRemoveLiquidityUnits, lp types.LiquidityProvider, pool types.Pool, withdrawNativeAssetAmount, withdrawExternalAssetAmount sdk.Uint, eAsset *tokenregistrytypes.RegistryEntry, pmtpCurrentRunningRate sdk.Dec) error {
	return nil
}

func FEATURE_TOGGLE_MARGIN_CLI_ALPHA_QueueRemovalWithWBasisPoints(ctx sdk.Context, k msgServer, msg *types.MsgRemoveLiquidity, lp types.LiquidityProvider, pool types.Pool, withdrawNativeAssetAmount, withdrawExternalAssetAmount sdk.Uint, eAsset *tokenregistrytypes.RegistryEntry, pmtpCurrentRunningRate sdk.Dec) error {
	return nil
}

func FEATURE_TOGGLE_MARGIN_CLI_ALPHA_SwapOne(ctx sdk.Context,
	k msgServer,
	sentAsset types.Asset,
	sentAmount sdk.Uint,
	nativeAsset types.Asset,
	inPool types.Pool,
	pmtpCurrentRunningRate sdk.Dec) (sdk.Uint, sdk.Uint, sdk.Uint, types.Pool, error) {
	return SwapOne(sentAsset, sentAmount, nativeAsset, inPool, pmtpCurrentRunningRate)
}

func FEATURE_TOGGLE_MARGIN_CLI_ALPHA_GetSwapFee(ctx sdk.Context,
	k msgServer,
	ReceivedAsset *types.Asset,
	liquidityFeeNative sdk.Uint,
	outPool types.Pool,
	pmtpCurrentRunningRate sdk.Dec) sdk.Uint {
	return GetSwapFee(liquidityFeeNative, *ReceivedAsset, outPool, pmtpCurrentRunningRate)
}
