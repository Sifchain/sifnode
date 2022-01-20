package main

import (
	"github.com/Sifchain/sifnode/app"
	sdkclient "github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/simapp/params"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	rpchttp "github.com/tendermint/tendermint/rpc/client/http"
)

func getClientAndFactory(args Args) (tx.Factory, sdkclient.Context) {
	uri := ""
	switch args.Network {
	case Devnet:
		uri = "https://rpc-devnet.sifchain.finance:443"
	case TestNet:
		uri = "https://rpc-testnet.sifchain.finance:443"
	case MainNet:
		uri = "https://rpc.sifchain.finance:443"
	case LocalNet:
		uri = "tcp://127.0.0.1:26657"
	default:
		panic("Network is a required arg")
	}
	client, err := newClient(uri)
	if err != nil {
		panic(err)
	}
	encConfig := app.MakeTestEncodingConfig()
	return newClientContext(uri, client, args, encConfig)
}

func newClientContext(uri string, client *rpchttp.HTTP, args Args, config params.EncodingConfig) (tx.Factory, sdkclient.Context) {
	txf := tx.Factory{}.
		WithChainID(args.ChainID).
		WithFees(args.Fees).
		WithKeybase(args.Keybase).
		WithAccountRetriever(authtypes.AccountRetriever{}).
		WithTxConfig(config.TxConfig).
		WithGas(100000000)

	clientCtx := sdkclient.Context{}.
		WithNodeURI(uri).
		WithClient(client).
		WithFrom(args.Sender.String()).
		WithFromAddress(args.Sender).
		WithTxConfig(config.TxConfig).
		WithInterfaceRegistry(config.InterfaceRegistry).
		WithSkipConfirmation(true).
		WithFromName(args.SenderName).
		WithBroadcastMode("block").
		WithOutputFormat("json")

	return txf, clientCtx
}

func newClient(uri string) (*rpchttp.HTTP, error) {
	rpcClient, err := rpchttp.New(uri, "/websocket")
	if err != nil {
		return nil, err
	}
	return rpcClient, nil
}
