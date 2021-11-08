package main

import (
	sdkclient "github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func BroadCast(txf tx.Factory, clientCtx sdkclient.Context, msg sdk.Msg) *sdk.TxResponse {
	preparedTfx, err := tx.PrepareFactory(clientCtx, txf)
	if err != nil {
		panic(err)
	}
	unsignedTx, err := tx.BuildUnsignedTx(preparedTfx, msg)
	if err != nil {
		panic(err)
	}
	err = tx.Sign(preparedTfx, clientCtx.GetFromName(), unsignedTx, true)
	if err != nil {
		panic(err)
	}

	txBytes, err := clientCtx.TxConfig.TxEncoder()(unsignedTx.GetTx())
	if err != nil {
		panic(err)
	}
	res, err := clientCtx.BroadcastTx(txBytes)
	if err != nil {
		panic(err)
	}
	return res
}
