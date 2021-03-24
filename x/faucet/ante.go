package faucet

import (
	"github.com/Sifchain/sifnode/x/clp"
	ftypes "github.com/Sifchain/sifnode/x/faucet/types"
	"github.com/cosmos/cosmos-sdk/types"
)

type RemoveFaucetFeeDecorator struct{}

//NewRemoveFacuetFeeDecorator returns a new RemoveFaucetFeeDecorator
func NewRemoveFacuetFeeDecorator() RemoveFaucetFeeDecorator {
	return RemoveFaucetFeeDecorator{}
}

var _ types.AnteDecorator = RemoveFaucetFeeDecorator{}

//AnteHandle for the RemoveFaucetFeeDecorator removes the MinGasPrice for a request token transaction
//This Function only works ith the profile is set as TestNet and the chain-ID is not sifchain
func (r RemoveFaucetFeeDecorator) AnteHandle(ctx types.Context, tx types.Tx, simulate bool, next types.AnteHandler) (newCtx types.Context, err error) {
	if ctx.ChainID() != "sifchain" && ftypes.Profile == ftypes.TESTNET {
		newDecCoin := types.NewDecCoinFromDec(clp.GetSettlementAsset().Symbol, types.NewDec(0))
		ctx = ctx.WithMinGasPrices(types.NewDecCoins(newDecCoin))
	}
	return next(ctx, tx, simulate)
}
