package app

import (
	"fmt"
	"github.com/Sifchain/sifnode/x/clp"
	"github.com/Sifchain/sifnode/x/faucet"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
)

func NewAnteHandler(ak auth.AccountKeeper, sk types.SupplyKeeper) sdk.AnteHandler {
	return func(ctx sdk.Context, tx sdk.Tx, sim bool) (newCtx sdk.Context, err error) {
		var anteHandler sdk.AnteHandler
		t, ok := tx.(faucet.MsgRequestCoins)
		if ok {
			fmt.Println(t.Type())
		}
		switch tx.(type) {
		case faucet.MsgRequestCoins:
			fmt.Println("Executing AnteHandler For RequestCoins")
			anteHandler = sdk.ChainAnteDecorators(
				faucet.NewRemoveFacuetFeeDecorator(),
				ante.NewSetUpContextDecorator(), // outermost AnteDecorator. SetUpContext must be called first
				ante.NewMempoolFeeDecorator(),
				ante.NewValidateBasicDecorator(),
				ante.NewValidateMemoDecorator(ak),
				ante.NewConsumeGasForTxSizeDecorator(ak),
				ante.NewSetPubKeyDecorator(ak), // SetPubKeyDecorator must be called before all signature verification decorators
				ante.NewValidateSigCountDecorator(ak),
				ante.NewDeductFeeDecorator(ak, sk),
				ante.NewSigGasConsumeDecorator(ak, auth.DefaultSigVerificationGasConsumer),
				ante.NewSigVerificationDecorator(ak),
				ante.NewIncrementSequenceDecorator(ak),
			)
		case auth.StdTx:
			fmt.Println("Executing AnteHandler For stdTX")
			anteHandler = sdk.ChainAnteDecorators(
				ante.NewSetUpContextDecorator(), // outermost AnteDecorator. SetUpContext must be called first
				ante.NewMempoolFeeDecorator(),
				ante.NewValidateBasicDecorator(),
				ante.NewValidateMemoDecorator(ak),
				ante.NewConsumeGasForTxSizeDecorator(ak),
				ante.NewSetPubKeyDecorator(ak), // SetPubKeyDecorator must be called before all signature verification decorators
				ante.NewValidateSigCountDecorator(ak),
				ante.NewDeductFeeDecorator(ak, sk),
				ante.NewSigGasConsumeDecorator(ak, auth.DefaultSigVerificationGasConsumer),
				ante.NewSigVerificationDecorator(ak),
				ante.NewIncrementSequenceDecorator(ak),
			)
		case clp.MsgSwap:
			fmt.Println("Executing AnteHandler For swap")
			anteHandler = sdk.ChainAnteDecorators(
				ante.NewSetUpContextDecorator(), // outermost AnteDecorator. SetUpContext must be called first
				ante.NewMempoolFeeDecorator(),
				ante.NewValidateBasicDecorator(),
				ante.NewValidateMemoDecorator(ak),
				ante.NewConsumeGasForTxSizeDecorator(ak),
				ante.NewSetPubKeyDecorator(ak), // SetPubKeyDecorator must be called before all signature verification decorators
				ante.NewValidateSigCountDecorator(ak),
				ante.NewDeductFeeDecorator(ak, sk),
				ante.NewSigGasConsumeDecorator(ak, auth.DefaultSigVerificationGasConsumer),
				ante.NewSigVerificationDecorator(ak),
				ante.NewIncrementSequenceDecorator(ak),
			)

		default:
			return ctx, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "invalid transaction type: %T", tx)
		}

		return anteHandler(ctx, tx, sim)
	}
}
