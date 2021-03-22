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
		anteHandler := GetDefaultAnteHandler(ak, sk)
		for _, msg := range tx.GetMsgs() {
			switch msg.Type() {
			case faucet.RequestCoinsType:
				anteHandler = GetFaucetAnteHandler(ak, sk)
			case clp.SwapType:
				anteHandler = GetSwapAnteHandler(ak, sk, ck)
			default:
				anteHandler = GetDefaultAnteHandler(ak, sk)
			}
		}
		return anteHandler(ctx, tx, sim)
	}
}

// GetDefaultAnteHandler returns the default antehandle for all transactions on sifchain
func GetDefaultAnteHandler(ak auth.AccountKeeper, sk types.SupplyKeeper) sdk.AnteHandler {
	return sdk.ChainAnteDecorators(
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

//GetFaucetAnteHandler adds a new decorator NewRemoveFacuetFeeDecorator to the default antehandler
func GetFaucetAnteHandler(ak auth.AccountKeeper, sk types.SupplyKeeper) sdk.AnteHandler {
	return sdk.ChainAnteDecorators(
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
}

//GetSwapAnteHandler adds a new decorator NewSwapFeeChangeDecorator to the default antehandler
func GetSwapAnteHandler(ak auth.AccountKeeper, sk types.SupplyKeeper, ck clp.Keeper) sdk.AnteHandler {
	return sdk.ChainAnteDecorators(
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
}
