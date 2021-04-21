package test

import (
	"bytes"
	"fmt"
	sifapp "github.com/Sifchain/sifnode/app"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/tendermint/tendermint/crypto"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	"strconv"
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

func GenerateInputList(rowanamount string) []types.Input {
	addressList := []string{"A58856F0FD53BF058B4909A21AEC019107BA6", "A58856F0FD53BF058B4909A21AEC019107BA7"}
	accAddrList := GenerateAddressList(addressList)
	rowan, ok := sdk.NewIntFromString(rowanamount)
	if !ok {
		panic(fmt.Sprintf("Err in getting amount : %s", rowanamount))
	}
	rowanAmount := sdk.Coins{sdk.NewCoin("rowan", rowan)}
	var res []types.Input
	for _, address := range accAddrList {
		in := types.NewInput(address, rowanAmount)
		res = append(res, in)
	}
	return res
}

func GenerateOutputList(rowanamount string) []types.Output {
	addressList := []string{"A58856F0FD53BF058B4909A21AEC019107BA3", "A58856F0FD53BF058B4909A21AEC019107BA4", "A58856F0FD53BF058B4909A21AEC019107BA5"}
	accAddrList := GenerateAddressList(addressList)
	rowan, ok := sdk.NewIntFromString(rowanamount)
	if !ok {
		panic(fmt.Sprintf("Err in getting amount : %s", rowanamount))
	}
	rowanAmount := sdk.Coins{sdk.NewCoin("rowan", rowan)}
	var res []types.Output
	for _, address := range accAddrList {
		out := types.NewOutput(address, rowanAmount)
		res = append(res, out)
	}
	return res
}

func GenerateAddressList(addressList []string) []sdk.AccAddress {
	var acclist []sdk.AccAddress
	for _, key := range addressList {
		var buffer bytes.Buffer
		buffer.WriteString(key)
		buffer.WriteString(strconv.Itoa(100))
		res, _ := sdk.AccAddressFromHex(buffer.String())
		bech := res.String()
		addr := buffer.String()
		res, err := sdk.AccAddressFromHex(addr)

		if err != nil {
			panic(err)
		}

		bechexpected := res.String()
		if bech != bechexpected {
			panic("Bech encoding doesn't match reference")
		}

		bechres, err := sdk.AccAddressFromBech32(bech)
		if err != nil {
			panic(err)
		}
		if !bytes.Equal(bechres, res) {
			panic("Bech decode and hex decode don't match")
		}
		acclist = append(acclist, res)
	}
	return acclist
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
