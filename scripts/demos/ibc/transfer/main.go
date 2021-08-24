package main

import (
	sdkclient "github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"
	rpchttp "github.com/tendermint/tendermint/rpc/client/http"

	transfertypes "github.com/cosmos/cosmos-sdk/x/ibc/applications/transfer/types"

	clienttypes "github.com/cosmos/cosmos-sdk/x/ibc/core/02-client/types"
)

type Args struct {
	chainID          string
	gasPrice         string
	gasAdjustment    float64
	keybase          keyring.Keyring
	channelId        string
	from             string
	sender           sdk.AccAddress
	receiver         string
	amount           sdk.Coin
	timeoutTimestamp uint64
}

func main() {

	uri := "http://rpc-testnet.sifchain.finance:80"
	client, err := newClient(uri)
	if err != nil {
		panic(err)
	}

	// TODO: Specify and setup keyring
	args := Args{}

	txf, clientCtx := newClientContext(uri, client, args)

	transferReq := transfertypes.NewMsgTransfer("transfer", args.channelId, args.amount, args.sender, args.receiver, clienttypes.NewHeight(0, 0), args.timeoutTimestamp)

	tx.BroadcastTx(clientCtx, txf, transferReq)

}

func debug(txHash string) {
	// Query to find out what the error / status is, log along the way.
}

func newClient(uri string) (*rpchttp.HTTP, error) {
	rpcClient, err := rpchttp.New(uri, "/websocket")
	if err != nil {
		return nil, err
	}

	return rpcClient, nil
}

func newClientContext(uri string, client *rpchttp.HTTP, args Args) (tx.Factory, sdkclient.Context) {
	txf := tx.Factory{}.
		WithChainID(args.chainID).
		WithGasPrices(args.gasPrice).
		WithGasAdjustment(args.gasAdjustment).
		WithKeybase(args.keybase)

	clientCtx := sdkclient.Context{}.
		WithNodeURI(uri).
		WithClient(client).
		WithFrom("")

	return txf, clientCtx
}
