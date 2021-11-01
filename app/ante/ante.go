package ante

import (
	disptypes "github.com/Sifchain/sifnode/x/dispensation/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	"github.com/cosmos/cosmos-sdk/x/auth/signing"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
)

func NewAnteHandler(ak authkeeper.AccountKeeper, bk bankkeeper.Keeper, gasConsumer ante.SignatureVerificationGasConsumer, signModeHandler signing.SignModeHandler) sdk.AnteHandler {
	return sdk.ChainAnteDecorators(
		ante.NewSetUpContextDecorator(), // outermost AnteDecorator. SetUpContext must be called first
		NewAdjustGasPriceDecorator(),    // Custom decorator to adjust gas price for specific msg types
		ante.NewRejectExtensionOptionsDecorator(),
		ante.NewMempoolFeeDecorator(),
		ante.NewValidateBasicDecorator(),
		ante.TxTimeoutHeightDecorator{},
		ante.NewValidateMemoDecorator(ak),
		ante.NewConsumeGasForTxSizeDecorator(ak),
		ante.NewRejectFeeGranterDecorator(),
		ante.NewSetPubKeyDecorator(ak), // SetPubKeyDecorator must be called before all signature verification decorators
		ante.NewValidateSigCountDecorator(ak),
		ante.NewDeductFeeDecorator(ak, bk),
		ante.NewSigGasConsumeDecorator(ak, gasConsumer),
		ante.NewSigVerificationDecorator(ak, signModeHandler),
		ante.NewIncrementSequenceDecorator(ak),
	)
}

// AdjustGasPriceDecorator is a custom decorator to reduce fee prices .
type AdjustGasPriceDecorator struct {
}

// NewAdjustGasPriceDecorator create a new instance of AdjustGasPriceDecorator
func NewAdjustGasPriceDecorator() AdjustGasPriceDecorator {
	return AdjustGasPriceDecorator{}
}

// AnteHandle adjusts the gas price based on the tx type.
func (r AdjustGasPriceDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	msgs := tx.GetMsgs()
	minGasPrice := sdk.DecCoin{
		Denom:  "rowan",
		Amount: sdk.MustNewDecFromStr("100000000000000000"),
	}
	if len(msgs) == 1 && (msgs[0].Type() == disptypes.MsgTypeCreateDistribution || msgs[0].Type() == disptypes.MsgTypeRunDistribution) {
		minGasPrice = sdk.DecCoin{
			Denom:  "rowan",
			Amount: sdk.MustNewDecFromStr("0.00000005"),
		}
	}
	if !minGasPrice.IsValid() {
		return ctx, sdkerrors.Wrap(sdkerrors.ErrLogic, "invalid gas price")
	}
	ctx = ctx.WithMinGasPrices(sdk.NewDecCoins(minGasPrice))
	return next(ctx, tx, simulate)
}
