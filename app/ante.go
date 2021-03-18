package app

import (
	"github.com/Sifchain/sifnode/x/clp"
	"github.com/Sifchain/sifnode/x/faucet"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
)

func NewAnteHandler(ak auth.AccountKeeper, sk types.SupplyKeeper, ck clp.Keeper) sdk.AnteHandler {
	return func(ctx sdk.Context, tx sdk.Tx, sim bool) (newCtx sdk.Context, err error) {
		var anteHandler sdk.AnteHandler
		for _, msg := range tx.GetMsgs() {
			switch msg.Type() {
			case faucet.RequestCoinsType:
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
			case clp.SwapType:
				anteHandler = sdk.ChainAnteDecorators(
					clp.NewSwapFeeChangeDecorator(ck),
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
			}
		}
		return anteHandler(ctx, tx, sim)
	}
}
