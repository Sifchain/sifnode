package ante

import (
	"strings"

	disptypes "github.com/Sifchain/sifnode/x/dispensation/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	"github.com/pkg/errors"
)

func NewAnteHandler(options ante.HandlerOptions) (sdk.AnteHandler, error) {

	if options.AccountKeeper == nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrLogic, "account keeper is required for ante builder")
	}

	if options.BankKeeper == nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrLogic, "bank keeper is required for ante builder")
	}

	if options.SignModeHandler == nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrLogic, "sign mode handler is required for ante builder")
	}

	var sigGasConsumer = options.SigGasConsumer
	if sigGasConsumer == nil {
		sigGasConsumer = ante.DefaultSigVerificationGasConsumer
	}
	return sdk.ChainAnteDecorators(
		ante.NewSetUpContextDecorator(), // outermost AnteDecorator. SetUpContext must be called first
		NewReduceGasPriceDecorator(),    // Custom decorator to reduce gas price for specific msg types
		ante.NewRejectExtensionOptionsDecorator(),
		ante.NewMempoolFeeDecorator(),
		ante.NewValidateBasicDecorator(),
		ante.NewTxTimeoutHeightDecorator(),
		ante.NewValidateMemoDecorator(options.AccountKeeper),
		ante.NewConsumeGasForTxSizeDecorator(options.AccountKeeper),
		ante.NewDeductFeeDecorator(options.AccountKeeper, options.BankKeeper, options.FeegrantKeeper),
		ante.NewSetPubKeyDecorator(options.AccountKeeper), // SetPubKeyDecorator must be called before all signature verification decorators
		ante.NewValidateSigCountDecorator(options.AccountKeeper),
		ante.NewSigGasConsumeDecorator(options.AccountKeeper, sigGasConsumer),
		ante.NewSigVerificationDecorator(options.AccountKeeper, options.SignModeHandler),
		ante.NewIncrementSequenceDecorator(options.AccountKeeper),
	), nil

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
		msgTypeURLLower := strings.ToLower(sdk.MsgTypeURL(msgs[i]))
		if strings.Contains(msgTypeURLLower, strings.ToLower(disptypes.MsgTypeCreateDistribution)) ||
			strings.Contains(msgTypeURLLower, strings.ToLower(disptypes.MsgTypeRunDistribution)) {
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
