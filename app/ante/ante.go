package ante

import (
	disptypes "github.com/Sifchain/sifnode/x/dispensation/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	"github.com/cosmos/cosmos-sdk/x/auth/signing"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	"github.com/pkg/errors"
)

func NewAnteHandler(ak authkeeper.AccountKeeper, bk bankkeeper.Keeper, gasConsumer ante.SignatureVerificationGasConsumer, signModeHandler signing.SignModeHandler) sdk.AnteHandler {
	return func(
		ctx sdk.Context, tx sdk.Tx, sim bool,
	) (newCtx sdk.Context, err error) {
		var anteHandler sdk.AnteHandler
		msgs := tx.GetMsgs()
		// If number of messages is greater than one ,the second message will be able to get away with a lower tx fee
		if len(msgs) > 1 && msgs[0].Type() == disptypes.MsgTypeCreateDistribution {
			return ctx, errors.New("Create Dispensation cannot be part of a multi message transaction")
		}
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
				ante.NewSigGasConsumeDecorator(ak, gasConsumer),
				ante.NewSigVerificationDecorator(ak, signModeHandler),
				ante.NewIncrementSequenceDecorator(ak),
			)
		}
		return anteHandler(ctx, tx, sim)
	}
}

// ReduceGasPriceDecorator is a custom decorator to reduce fee prices .
type ReduceGasPriceDecorator struct {
}

// NewReduceGasPriceDecorator create a new instance of ReduceGasPriceDecorator
func NewReduceGasPriceDecorator() ReduceGasPriceDecorator {
	return ReduceGasPriceDecorator{}
}

// AnteHandle reduces the gas price to a lower value which is hardcoded.The ReduceGasPriceDecorator should only be used for specific transaction types to lower the fee cost.
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
