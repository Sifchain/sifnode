package main

import (
	"fmt"
	"github.com/Sifchain/sifnode/app"
	sdkclient "github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/simapp/params"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bank "github.com/cosmos/cosmos-sdk/x/bank/types"
	rpchttp "github.com/tendermint/tendermint/rpc/client/http"
	"os"
)

type Network int8

const (
	AccountAddressPrefix         = "sif"
	Devnet               Network = 1
	TestNet              Network = 2
	MainNet              Network = 3
)

var (
	AccountPubKeyPrefix    = AccountAddressPrefix + "pub"
	ValidatorAddressPrefix = AccountAddressPrefix + "valoper"
	ValidatorPubKeyPrefix  = AccountAddressPrefix + "valoperpub"
	ConsNodeAddressPrefix  = AccountAddressPrefix + "valcons"
	ConsNodePubKeyPrefix   = AccountAddressPrefix + "valconspub"
)

func SetConfig(seal bool) {
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(AccountAddressPrefix, AccountPubKeyPrefix)
	config.SetBech32PrefixForValidator(ValidatorAddressPrefix, ValidatorPubKeyPrefix)
	config.SetBech32PrefixForConsensusNode(ConsNodeAddressPrefix, ConsNodePubKeyPrefix)
	if seal {
		config.Seal()
	}
}

type Args struct {
	ChainID          string          `json:"chain_id,omitempty"`
	GasPrice         string          `json:"gas_price,omitempty"`
	GasAdjustment    float64         `json:"gas_adjustment,omitempty"`
	Keybase          keyring.Keyring `json:"keybase,omitempty"`
	ChannelId        string          `json:"channel_id,omitempty"`
	Sender           sdk.AccAddress  `json:"sender,omitempty"`
	Receiver         sdk.AccAddress  `json:"receiver,omitempty"`
	Amount           sdk.Coins       `json:"amount"`
	TimeoutTimestamp uint64          `json:"timeout_timestamp,omitempty"`
	Fees             string          `json:"fees"`
}

func main() {
	SetConfig(true)
	args := getDevnetNetArgs()
	txf, clientCtx := getClientAndFactory(Devnet, *args)

	//transferReq := transfertypes.NewMsgTransfer("transfer", args.ChannelId, args.Amount, args.Sender, args.Receiver, clienttypes.NewHeight(0, 0), args.TimeoutTimestamp)
	sendReq := bank.NewMsgSend(args.Sender, args.Receiver, args.Amount)
	BroadCastAndCheck(txf, clientCtx, sendReq)

}

func BroadCastAndCheck(txf tx.Factory, clientCtx sdkclient.Context, sendReq *bank.MsgSend) {
	preparedTfx, err := tx.PrepareFactory(clientCtx, txf)
	unsignedTx, err := tx.BuildUnsignedTx(preparedTfx, sendReq)
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
	debug(res)
}

func debug(res *sdk.TxResponse) {
	// Works only in block
	if res.Code != 0 {
		panic("Transaction Failed")
	}
	fmt.Println(res.Logs)
	fmt.Println(res.TxHash)
}

func getClientAndFactory(network Network, args Args) (tx.Factory, sdkclient.Context) {
	switch network {
	case 1:
		{
			uri := "https://rpc-devnet.sifchain.finance:443"
			client, err := newClient(uri)
			if err != nil {
				panic(err)
			}
			encConfig := app.MakeTestEncodingConfig()
			return newClientContext(uri, client, args, encConfig)

		}
	default:
		panic("nocase")
	}
}

func newClientContext(uri string, client *rpchttp.HTTP, args Args, config params.EncodingConfig) (tx.Factory, sdkclient.Context) {
	txf := tx.Factory{}.
		WithChainID(args.ChainID).
		WithFees(args.Fees).
		WithKeybase(args.Keybase).
		WithAccountRetriever(authtypes.AccountRetriever{}).
		WithTxConfig(config.TxConfig).
		WithGas(1000000)

	clientCtx := sdkclient.Context{}.
		WithNodeURI(uri).
		WithClient(client).
		WithFrom(args.Sender.String()).
		WithFromAddress(args.Sender).
		WithTxConfig(config.TxConfig).
		WithInterfaceRegistry(config.InterfaceRegistry).
		WithSkipConfirmation(true).
		WithFromName("sif").
		WithBroadcastMode("block").
		WithOutputFormat("json").
		WithJSONMarshaler(config.Marshaler)

	return txf, clientCtx
}

func newClient(uri string) (*rpchttp.HTTP, error) {
	rpcClient, err := rpchttp.New(uri, "/websocket")
	if err != nil {
		return nil, err
	}
	return rpcClient, nil
}

// TODO : Replace this function to read from file
func getDevnetNetArgs() *Args {
	amount, ok := sdk.NewIntFromString("100000000000000000000000")
	if !ok {
		panic("Cannot parse amount")
	}
	path := hd.CreateHDPath(118, 0, 0).String()

	kr, err := keyring.New("sifchain", "test", os.TempDir(), nil)
	if err != nil {
		panic(err)
	}
	mnemonic := "race draft rival universe maid cheese steel logic crowd fork comic easy truth drift tomorrow eye buddy head time cash swing swift midnight borrow"

	// TODO improve this logic
	accInfo, err := kr.NewAccount("sif", mnemonic, "", path, hd.Secp256k1)
	if err != nil {
		accInfo, err = kr.Key("sif")
		if err != nil {
			panic(err)
		}
	}

	toAddr, err := sdk.AccAddressFromBech32("sif1l7hypmqk2yc334vc6vmdwzp5sdefygj2ad93p5")
	if err != nil {
		panic(toAddr)
	}

	return &Args{
		ChainID:          "sifchain-devnet-1",
		GasPrice:         "",
		GasAdjustment:    0,
		Keybase:          kr,
		ChannelId:        "",
		Sender:           accInfo.GetAddress(),
		Receiver:         toAddr,
		Amount:           sdk.NewCoins(sdk.NewCoin("rowan", amount)),
		TimeoutTimestamp: 0,
		Fees:             "1000000rowan",
	}
}
