package main

import (
	"github.com/Sifchain/sifnode/x/dispensation/test"
	dispensationtypes "github.com/Sifchain/sifnode/x/dispensation/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/tendermint/tendermint/crypto"
)

type CreateDispensationNegativeTx struct{} //nolint

func (CreateDispensationNegativeTx) GetMsgAndArgs(_ CommonArgs) (sdk.Msg, Args) {
	args := getDispensationTxArgs()
	output := test.CreatOutputList(9, "10000000000000000000")
	amount, ok := sdk.NewIntFromString("10000000000000000000")
	if !ok {
		panic("Unable to create amount")
	}
	unknownCoin := sdk.NewCoins(sdk.NewCoin("unknown", amount))
	doesntMatterAddress := sdk.AccAddress(crypto.AddressHash([]byte("Random")))
	output = append(output, types.Output{
		Address: doesntMatterAddress.String(),
		Coins:   unknownCoin,
	})
	createDispensation := dispensationtypes.NewMsgCreateDistribution(args.Sender, dispensationtypes.DistributionType_DISTRIBUTION_TYPE_VALIDATOR_SUBSIDY, output, args.Sender.String())
	return &createDispensation, args
}
func (CreateDispensationNegativeTx) GetName() string {
	return "CREATE-DISPENSATION"
}

func (s CreateDispensationNegativeTx) Assert(response *sdk.TxResponse, _ *CommonArgs) {
	if response.Code == 0 {
		panic("Test Failed , Transaction successfully submitted ")
	}
}
