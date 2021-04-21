package test

import (
	"github.com/Sifchain/sifnode/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/supply"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
	"strconv"
)

func CreateTestApp(isCheckTx bool) (*simapp.SimApp, sdk.Context) {
	app := simapp.Setup(isCheckTx)
	ctx := app.BaseApp.NewContext(isCheckTx, abci.Header{})
	app.AccountKeeper.SetParams(ctx, authtypes.DefaultParams())
	initTokens := sdk.TokensFromConsensusPower(1000)
	app.SupplyKeeper.SetSupply(ctx, supply.NewSupply(sdk.Coins{}))
	_ = simapp.AddTestAddrs(app, ctx, 6, initTokens)
	return app, ctx
}

func CreatOutputList(count int, rowanAmount string) []bank.Output {
	outputList := make([]bank.Output, count)
	amount, ok := sdk.NewIntFromString(rowanAmount)
	if !ok {
		panic("Unable to generate rowan amount")
	}
	coin := sdk.Coins{sdk.NewCoin("rowan", amount)}
	for i := 0; i < count; i++ {
		address := sdk.AccAddress(crypto.AddressHash([]byte("Output1" + strconv.Itoa(i))))
		out := bank.NewOutput(address, coin)
		outputList[i] = out
	}
	return outputList
}

func CreatInputList(count int, rowanAmount string) []bank.Input {
	list := make([]bank.Input, count)
	amount, ok := sdk.NewIntFromString(rowanAmount)
	if !ok {
		panic("Unable to generate rowan amount")
	}
	coin := sdk.Coins{sdk.NewCoin("rowan", amount)}
	for i := 0; i < count; i++ {
		address := sdk.AccAddress(crypto.AddressHash([]byte("Input1" + strconv.Itoa(i))))
		out := bank.NewInput(address, coin)
		list[i] = out
	}
	return list
}
