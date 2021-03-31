package test

import (
	"bytes"
	"fmt"
	"github.com/Sifchain/sifnode/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/supply"
	abci "github.com/tendermint/tendermint/abci/types"
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

func GenerateInputList(rowanamount string) []bank.Input {
	addressList := []string{"A58856F0FD53BF058B4909A21AEC019107BA6", "A58856F0FD53BF058B4909A21AEC019107BA7"}
	accAddrList := GenerateAddressList(addressList)
	rowan, ok := sdk.NewIntFromString(rowanamount)
	if !ok {
		panic(fmt.Sprintf("Err in getting amount : %s", rowanamount))
	}
	rowanAmount := sdk.Coins{sdk.NewCoin("rowan", rowan)}
	res := []bank.Input{}
	for _, address := range accAddrList {
		in := bank.NewInput(address, rowanAmount)
		res = append(res, in)
	}
	return res
}

func GenerateOutputList(rowanamount string) []bank.Output {
	addressList := []string{"A58856F0FD53BF058B4909A21AEC019107BA3", "A58856F0FD53BF058B4909A21AEC019107BA4", "A58856F0FD53BF058B4909A21AEC019107BA5"}
	accAddrList := GenerateAddressList(addressList)
	rowan, ok := sdk.NewIntFromString(rowanamount)
	if !ok {
		panic(fmt.Sprintf("Err in getting amount : %s", rowanamount))
	}
	rowanAmount := sdk.Coins{sdk.NewCoin("rowan", rowan)}
	res := []bank.Output{}
	for _, address := range accAddrList {
		out := bank.NewOutput(address, rowanAmount)
		res = append(res, out)
	}
	return res
}

func GenerateAddressList(addressList []string) []sdk.AccAddress {
	var acclist [len(addressList)]sdk.AccAddress
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
