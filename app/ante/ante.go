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
	return sdk.ChainAnteDecorators(
		ante.NewSetUpContextDecorator(), // outermost AnteDecorator. SetUpContext must be called first
		NewReduceGasPriceDecorator(),    // Custom decorator to reduce gas price for specific msg types
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

// ReduceGasPriceDecorator is a custom decorator to reduce fee prices .
type ReduceGasPriceDecorator struct {
}

// NewReduceGasPriceDecorator create a new instance of ReduceGasPriceDecorator
func NewReduceGasPriceDecorator() ReduceGasPriceDecorator {
	return ReduceGasPriceDecorator{}
}

// AnteHandle reduces the gas price to a lower value which is hardcoded.The ReduceGasPriceDecorator should only be used for specific transaction types to lower the fee cost.
func (r ReduceGasPriceDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	msgs := tx.GetMsgs()

	var found bool
	for i := range msgs {
		if msgs[i].Type() == disptypes.MsgTypeCreateDistribution || msgs[i].Type() == disptypes.MsgTypeRunDistribution {
			found = true
		}
	}

	// Pass earlier if not a dispensation tx.
	if !found {
		return next(ctx, tx, simulate)
	}

	// If number of messages is greater than one, the other messages will be able to get away with a lower tx fee
	if len(msgs) != 1 {
		return ctx, errors.New("transaction for dispensation create / run must have exactly one message")
	}

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
