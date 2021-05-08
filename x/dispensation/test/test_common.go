package test

import (
	"strconv"
	"time"

	sifapp "github.com/Sifchain/sifnode/app"
	dispensation "github.com/Sifchain/sifnode/x/dispensation/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/tendermint/tendermint/crypto"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
)

func CreateTestApp(isCheckTx bool) (*sifapp.SifchainApp, sdk.Context) {
	app := sifapp.Setup(isCheckTx)
	ctx := app.BaseApp.NewContext(isCheckTx, tmproto.Header{})
	app.AccountKeeper.SetParams(ctx, authtypes.DefaultParams())
	initTokens := sdk.TokensFromConsensusPower(1000)
	app.BankKeeper.SetSupply(ctx, types.NewSupply(sdk.Coins{}))
	_ = sifapp.AddTestAddrs(app, ctx, 6, initTokens)
	return app, ctx
}

func CreatOutputList(count int, rowanAmount string) []types.Output {
	outputList := make([]types.Output, count)
	amount, ok := sdk.NewIntFromString(rowanAmount)
	if !ok {
		panic("Unable to generate rowan amount")
	}
	coin := sdk.Coins{sdk.NewCoin("rowan", amount)}
	for i := 0; i < count; i++ {
		address := sdk.AccAddress(crypto.AddressHash([]byte("Output1" + strconv.Itoa(i))))
		out := types.NewOutput(address, coin)
		outputList[i] = out
	}
	return outputList
}

func CreatInputList(count int, rowanAmount string) []types.Input {
	list := make([]types.Input, count)
	amount, ok := sdk.NewIntFromString(rowanAmount)
	if !ok {
		panic("Unable to generate rowan amount")
	}
	coin := sdk.Coins{sdk.NewCoin("rowan", amount)}
	for i := 0; i < count; i++ {
		address := sdk.AccAddress(crypto.AddressHash([]byte("Output1" + strconv.Itoa(i))))
		out := types.NewInput(address, coin)
		list[i] = out
	}
	return list
}

func CreateClaimsList(count int, claimType dispensation.DistributionType) []dispensation.UserClaim {
	list := make([]dispensation.UserClaim, count)
	for i := 0; i < count; i++ {
		address := sdk.AccAddress(crypto.AddressHash([]byte("User" + strconv.Itoa(i))))
		claim := dispensation.NewUserClaim(address.String(), claimType, time.Now().String())
		list[i] = claim
	}
	return list
}
