package ante

import (
	"strings"

	adminkeeper "github.com/Sifchain/sifnode/x/admin/keeper"
	clptypes "github.com/Sifchain/sifnode/x/clp/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
)

const (
	// min gas price for rowan
	MIN_GAS_PRICE_UROWAN = "100000000000000"       // 0.0001rowan
	LOW_MIN_FEE_UROWAN   = "20000000000000000000"  // 20rowan
	HIGH_MIN_FEE_UROWAN  = "200000000000000000000" // 200rowan
)

func NewAnteHandler(options HandlerOptions) (sdk.AnteHandler, error) {
	if options.AccountKeeper == nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrLogic, "account keeper is required for ante builder")
	}
	if options.SignModeHandler == nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrLogic, "sign mode handler is required for ante builder")
	}
	var sigGasConsumer = options.SigGasConsumer
	if sigGasConsumer == nil {
		sigGasConsumer = ante.DefaultSigVerificationGasConsumer
	}
	return sdk.ChainAnteDecorators(
		ante.NewSetUpContextDecorator(),                 // outermost AnteDecorator. SetUpContext must be called first
		NewAdjustGasPriceDecorator(options.AdminKeeper), // Custom decorator to adjust gas price for specific msg types
		ante.NewRejectExtensionOptionsDecorator(),
		ante.NewMempoolFeeDecorator(),
		ante.NewValidateBasicDecorator(),
		ante.NewTxTimeoutHeightDecorator(),
		ante.NewValidateMemoDecorator(options.AccountKeeper),
		ante.NewConsumeGasForTxSizeDecorator(options.AccountKeeper),
		ante.NewDeductFeeDecorator(options.AccountKeeper, options.BankKeeper, options.FeegrantKeeper),
		NewValidateMinCommissionDecorator(options.StakingKeeper, options.BankKeeper), // Custom decorator to ensure the minimum commission rate of validators
		ante.NewSetPubKeyDecorator(options.AccountKeeper),                            // SetPubKeyDecorator must be called before all signature verification decorators
		ante.NewValidateSigCountDecorator(options.AccountKeeper),
		ante.NewSigGasConsumeDecorator(options.AccountKeeper, sigGasConsumer),
		ante.NewSigVerificationDecorator(options.AccountKeeper, options.SignModeHandler),
		ante.NewIncrementSequenceDecorator(options.AccountKeeper),
	), nil

}

// AdjustGasPriceDecorator is a custom decorator to reduce fee prices .
type AdjustGasPriceDecorator struct {
	adminKeeper adminkeeper.Keeper
}

// NewAdjustGasPriceDecorator create a new instance of AdjustGasPriceDecorator
func NewAdjustGasPriceDecorator(adminKeeper adminkeeper.Keeper) AdjustGasPriceDecorator {
	return AdjustGasPriceDecorator{adminKeeper: adminKeeper}
}

// AnteHandle adjusts the gas price based on the tx type.
func (r AdjustGasPriceDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	adminParams := r.adminKeeper.GetParams(ctx)
	submitProposalFee := adminParams.SubmitProposalFee

	// Get the symbol of the settlement asset
	settlementAssetSymbol := clptypes.GetSettlementAsset().Symbol

	if !ctx.MinGasPrices().IsValid() {
		return ctx, sdkerrors.Wrap(sdkerrors.ErrLogic, "invalid gas price")
	}

	// Define the global minimum gas price
	minGasPrice := sdk.DecCoin{
		Denom:  settlementAssetSymbol,
		Amount: sdk.MustNewDecFromStr(MIN_GAS_PRICE_UROWAN),
	}
	if !minGasPrice.IsValid() {
		return ctx, sdkerrors.Wrap(sdkerrors.ErrLogic, "invalid gas price")
	}

	// Get current minimum gas prices from context
	currentMinGasPrices := ctx.MinGasPrices()

	// Check and update context's minimum gas prices if necessary
	if currentAssetPrice := currentMinGasPrices.AmountOf(settlementAssetSymbol); currentAssetPrice.LT(minGasPrice.Amount) {
		// Replace the current minimum gas price with the new minimum gas price for the asset
		updatedMinGasPrices := currentMinGasPrices.Sub(sdk.NewDecCoins().Add(sdk.NewDecCoinFromDec(settlementAssetSymbol, currentAssetPrice))).Add(minGasPrice)
		ctx = ctx.WithMinGasPrices(updatedMinGasPrices)
	}

	highMinFee, ok := sdk.NewIntFromString(HIGH_MIN_FEE_UROWAN)
	if !ok {
		return ctx, sdkerrors.Wrap(sdkerrors.ErrLogic, "invalid high fee amount")
	}
	lowMinFee, ok := sdk.NewIntFromString(LOW_MIN_FEE_UROWAN)
	if !ok {
		return ctx, sdkerrors.Wrap(sdkerrors.ErrLogic, "invalid low fee amount")
	}

	msgs := tx.GetMsgs()
	minFee := sdk.ZeroInt()
	for i := range msgs {
		msgTypeURLLower := strings.ToLower(sdk.MsgTypeURL(msgs[i]))
		if strings.Contains(msgTypeURLLower, strings.ToLower(banktypes.TypeMsgSend)) ||
			strings.Contains(msgTypeURLLower, strings.ToLower(banktypes.TypeMsgMultiSend)) ||
			strings.Contains(msgTypeURLLower, "createuserclaim") ||
			strings.Contains(msgTypeURLLower, "swap") ||
			strings.Contains(msgTypeURLLower, "removeliquidity") ||
			strings.Contains(msgTypeURLLower, "removeliquidityunits") ||
			strings.Contains(msgTypeURLLower, "addliquidity") {
			minFee = sdk.MaxInt(minFee, highMinFee)
		} else if strings.Contains(msgTypeURLLower, "transfer") && minFee.LTE(sdk.NewInt(10000000000000000)) {
			minFee = sdk.MaxInt(minFee, lowMinFee)
		} else if strings.Contains(msgTypeURLLower, "submitproposal") || strings.Contains(msgTypeURLLower, govtypes.TypeMsgSubmitProposal) {
			minFee = sdk.MaxInt(minFee, sdk.NewIntFromBigInt(submitProposalFee.BigInt()))
		}
	}
	if minFee.Equal(sdk.ZeroInt()) {
		return next(ctx, tx, simulate)
	}
	feeTx, ok := tx.(sdk.FeeTx)
	if !ok {
		return ctx, sdkerrors.Wrap(sdkerrors.ErrTxDecode, "tx must be a FeeTx")
	}
	fees := feeTx.GetFee()
	rowanFee := fees.AmountOf(settlementAssetSymbol)
	if rowanFee.LTE(sdk.ZeroInt()) {
		return ctx, sdkerrors.Wrap(sdkerrors.ErrLogic, "unsupported fee asset")
	}
	if rowanFee.LT(minFee) {
		return ctx, sdkerrors.Wrap(sdkerrors.ErrLogic, "tx fee is too low")
	}
	return next(ctx, tx, simulate)
}
