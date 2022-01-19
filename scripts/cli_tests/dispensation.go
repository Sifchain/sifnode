package main

import (
	"github.com/Sifchain/sifnode/x/dispensation/test"
	dispensationtypes "github.com/Sifchain/sifnode/x/dispensation/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type CreateDispensationTx struct{}

func (CreateDispensationTx) GetMsgAndArgs(_ CommonArgs) (sdk.Msg, Args) {
	args := getDispensationTxArgs()
	output := test.CreatOutputList(2, "10000000000000000000")
	createDispensation := dispensationtypes.NewMsgCreateDistribution(args.Sender, dispensationtypes.DistributionType_DISTRIBUTION_TYPE_VALIDATOR_SUBSIDY, output, args.Sender.String())
	return &createDispensation, args
}
func (CreateDispensationTx) GetName() string {
	return "CREATE-DISPENSATION"
}

func (s CreateDispensationTx) Assert(response *sdk.TxResponse, _ *CommonArgs) {
	defaultAssert(response, s.GetName())
}

func getDispensationTxArgs() Args {
	defaultArgs := getDefaultArgs()
	setNetwork(&defaultArgs, LocalNet)
	return defaultArgs
}
