package main

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	transfertypes "github.com/cosmos/cosmos-sdk/x/ibc/applications/transfer/types"
	clienttypes "github.com/cosmos/cosmos-sdk/x/ibc/core/02-client/types"
)

type IbcSentTx struct {
}

func (IbcSentTx) GetMsgAndArgs() (sdk.Msg, Args) {
	args := getIbcSentTxArgs()
	transferReq := transfertypes.NewMsgTransfer("transfer", args.ChannelId, args.Amount[0], args.Sender, args.CosmosReceiver, clienttypes.NewHeight(0, 18446744073709551615), args.TimeoutTimestamp)
	return transferReq, args
}

func (IbcSentTx) GetName() string {
	return "IBC-SEND"
}

func (i IbcSentTx) Assert(response *sdk.TxResponse) {
	commonAssert(response, i.GetName())
}

func getIbcSentTxArgs() Args {
	commonArgs := getCommonArgs()
	commonArgs.ChannelId = "channel-101"
	commonArgs.CosmosReceiver = "cosmos1syavy2npfyt9tcncdtsdzf7kny9lh777pahuux"
	commonArgs.TimeoutTimestamp = 0
	commonArgs.Amount = sdk.NewCoins(sdk.NewCoin("ibc/C782C1DE5F380BC8A5B7D490684894B439D31847A004B271D7B7BA07751E582A", sdk.OneInt()))
	return commonArgs
}
