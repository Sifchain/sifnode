package main

import (
	dispensationtypes "github.com/Sifchain/sifnode/x/dispensation/types"
	dispensationutils "github.com/Sifchain/sifnode/x/dispensation/utils"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type CreateDispensationTx struct{}

func (CreateDispensationTx) GetMsgAndArgs(_ InterTestArgs) (sdk.Msg, Args) {
	args := getDispensationTxArgs()
	output, err := dispensationutils.ParseOutput("output.json")
	if err != nil {
		panic(err)
	}
	createDispensation := dispensationtypes.NewMsgCreateDistribution(args.Sender, dispensationtypes.DistributionType_DISTRIBUTION_TYPE_VALIDATOR_SUBSIDY, output, args.Sender.String())
	return &createDispensation, args
}
func (CreateDispensationTx) GetName() string {
	return "CREATE-DISPENSATION"
}

func (s CreateDispensationTx) Assert(response *sdk.TxResponse, _ *InterTestArgs) {
	commonAssert(response, s.GetName())
}

func getDispensationTxArgs() Args {
	commonArgs := getCommonArgs()
	setNetwork(&commonArgs, LocalNet)
	return commonArgs
}
