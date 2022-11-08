package main

import (
	dispensationtypes "github.com/Sifchain/sifnode/x/dispensation/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/tendermint/tendermint/crypto"
)

type CreateDispensationMultiTx struct{} //nolint

func (CreateDispensationMultiTx) GetMsgAndArgs(_ CommonArgs) (sdk.Msg, Args) {
	args := getDispensationTxArgs()
	amount, ok := sdk.NewIntFromString("10000000000000000000")
	if !ok {
		panic("Failed to get amount")
	}
	address := sdk.AccAddress(crypto.AddressHash([]byte("Output")))
	coinRowan := sdk.NewCoins(sdk.NewCoin("rowan", amount))
	coinCeth := sdk.NewCoins(sdk.NewCoin("ceth", amount))
	output := []types.Output{types.NewOutput(address, coinCeth), types.NewOutput(address, coinRowan)}

	createDispensation := dispensationtypes.NewMsgCreateDistribution(args.Sender, dispensationtypes.DistributionType_DISTRIBUTION_TYPE_VALIDATOR_SUBSIDY, output, args.Sender.String())
	return &createDispensation, args
}
func (CreateDispensationMultiTx) GetName() string {
	return "CREATE-DISPENSATION"
}

func (s CreateDispensationMultiTx) Assert(response *sdk.TxResponse, _ *CommonArgs) {
	defaultAssert(response, s.GetName())
}
