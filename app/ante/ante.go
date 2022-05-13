package ante

import (
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	"github.com/Sifchain/sifnode/tools/slicex"
	disptypes "github.com/Sifchain/sifnode/x/dispensation/types"
)

// Predefined errors
var (
	ErrAccountMissing         = sdkerrors.Wrap(sdkerrors.ErrLogic, "account keeper is required for ante builder")
	ErrBankKeeperMissing      = sdkerrors.Wrap(sdkerrors.ErrLogic, "bank keeper is required for ante builder")
	ErrSignModeHandlerMissing = sdkerrors.Wrap(sdkerrors.ErrLogic, "sign mode handler is required for ante builder")
	ErrInvalidGasPrice        = sdkerrors.Wrap(sdkerrors.ErrLogic, "invalid gas price")
	ErrUnsupportedAsset       = sdkerrors.Wrap(sdkerrors.ErrLogic, "unsupported fee asset")
	ErrLowFee                 = sdkerrors.Wrap(sdkerrors.ErrLogic, "tx fee is too low")
)

var (
	distributionMessageTypes = []string{
		strings.ToLower(disptypes.MsgTypeCreateDistribution),
		strings.ToLower(disptypes.MsgTypeRunDistribution),
	}

	regularMessageTypes = []string{
		strings.ToLower(banktypes.TypeMsgSend),
		strings.ToLower(banktypes.TypeMsgMultiSend),
		"createuserclaim",
		"swap",
		"removeliquidity",
		"removeliquidityunits",
		"addliquidity",
	}
)

// NewAnteHandler is the constructor of sdk.AnteHandler.
func NewAnteHandler(options ante.HandlerOptions) (sdk.AnteHandler, error) {
	if options.AccountKeeper == nil {
		return nil, ErrAccountMissing
	}

	if options.BankKeeper == nil {
		return nil, ErrBankKeeperMissing
	}

	if options.SignModeHandler == nil {
		return nil, ErrSignModeHandlerMissing
	}

	sigGasConsumer := options.SigGasConsumer
	if sigGasConsumer == nil {
		sigGasConsumer = ante.DefaultSigVerificationGasConsumer
	}

	return sdk.ChainAnteDecorators(
		ante.NewSetUpContextDecorator(), // outermost AnteDecorator. SetUpContext must be called first
		NewAdjustGasPriceDecorator(),    // Custom decorator to adjust gas price for specific msg types
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
	if len(msgs) == 1 && isDistributedMsg(msgs[0]) {
		minGasPrice := sdk.DecCoin{
			Denom:  "rowan",
			Amount: sdk.MustNewDecFromStr("0.00000005"),
		}

		if !minGasPrice.IsValid() {
			return ctx, ErrInvalidGasPrice
		}

		ctx = ctx.WithMinGasPrices(sdk.NewDecCoins(minGasPrice))
		return next(ctx, tx, simulate)
	}

	minFee := sdk.ZeroInt()
	for i := range msgs {
		if isRegularMessage(msgs[i]) {
			minFee = sdk.NewInt(100000000000000000) // 0.1
		} else if isTransferMessage(msgs[i]) && minFee.LTE(sdk.NewInt(10000000000000000)) {
			minFee = sdk.NewInt(10000000000000000) // 0.01
		}
	}

	if minFee.Equal(sdk.ZeroInt()) {
		return next(ctx, tx, simulate)
	}

	feeTx, ok := tx.(sdk.FeeTx)
	if !ok {
		return ctx, sdkerrors.Wrapf(sdkerrors.ErrTxDecode, "tx must be a FeeTx, not %T", tx)
	}

	fees := feeTx.GetFee()
	rowanFee := sdk.ZeroInt()
	for j := range fees {
		if strings.EqualFold("rowan", fees[j].Denom) {
			rowanFee = fees[j].Amount
		}
	}

	if rowanFee.LTE(sdk.ZeroInt()) {
		return ctx, ErrUnsupportedAsset
	}

	if rowanFee.LT(minFee) {
		return ctx, ErrLowFee
	}

	return next(ctx, tx, simulate)
}

func isDistributedMsg(msg sdk.Msg) bool {
	return slicex.StringsContain(strings.ToLower(sdk.MsgTypeURL(msg)), distributionMessageTypes)
}

func isRegularMessage(msg sdk.Msg) bool {
	return slicex.StringsContain(strings.ToLower(sdk.MsgTypeURL(msg)), regularMessageTypes)
}

func isTransferMessage(msg sdk.Msg) bool {
	return strings.Contains(strings.ToLower(sdk.MsgTypeURL(msg)), "transfer")
}
