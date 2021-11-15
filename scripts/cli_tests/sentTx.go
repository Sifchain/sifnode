package main

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	bank "github.com/cosmos/cosmos-sdk/x/bank/types"
)

type SentTx struct{}

func (SentTx) GetMsgAndArgs(_ CommonArgs) (sdk.Msg, Args) {
	args := getSendTxArgs()
	sendReq := bank.NewMsgSend(args.Sender, args.SifchainReceiver, args.Amount)
	return sendReq, args
}
func (SentTx) GetName() string {
	return "SEND"
}

func (s SentTx) Assert(response *sdk.TxResponse, _ *CommonArgs) {
	defaultAssert(response, s.GetName())
}

func getSendTxArgs() Args {
	defaultArgs := getDefaultArgs()
	setNetwork(&defaultArgs, LocalNet)
	return defaultArgs
}
