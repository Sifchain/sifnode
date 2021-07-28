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
	return func(
		ctx sdk.Context, tx sdk.Tx, sim bool,
	) (newCtx sdk.Context, err error) {
		var anteHandler sdk.AnteHandler
		msgs := tx.GetMsgs()
		// TODO change to iteration over msgs
		switch msgs[0].Type() {
		case disptypes.MsgTypeCreateDistribution:
			anteHandler = sdk.ChainAnteDecorators(
				ante.NewSetUpContextDecorator(), // outermost AnteDecorator. SetUpContext must be called first
				NewReduceGasPriceDecorator(),
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
		default:
			anteHandler = sdk.ChainAnteDecorators(
				ante.NewSetUpContextDecorator(), // outermost AnteDecorator. SetUpContext must be called first
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
				ante.NewSigGasConsumeDecorator(ak, ante.DefaultSigVerificationGasConsumer),
				ante.NewSigVerificationDecorator(ak, signModeHandler),
				ante.NewIncrementSequenceDecorator(ak),
			)
		}
		return anteHandler(ctx, tx, sim)
	}
}

type ReduceGasPriceDecorator struct {
}

func NewReduceGasPriceDecorator() ReduceGasPriceDecorator {
	return ReduceGasPriceDecorator{}
}

func (r ReduceGasPriceDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	loweredGasPrice := sdk.DecCoin{
		Denom:  "rowan",
		Amount: sdk.MustNewDecFromStr("0.00000005"),
	}
	if !loweredGasPrice.IsValid() {
		return ctx, sdkerrors.Wrap(sdkerrors.ErrLogic, "unable to lower gas price")
	}
	ctx = ctx.WithMinGasPrices(sdk.NewDecCoins(loweredGasPrice))
	return next(ctx, tx, simulate)
}
