package main

import (
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Network int8

const (
	AccountAddressPrefix         = "sif"
	Devnet               Network = 1
	TestNet              Network = 2
	MainNet              Network = 3
	LocalNet             Network = 4
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

type TestCase interface {
	GetMsgAndArgs(iArgs CommonArgs) (sdk.Msg, Args)
	GetName() string
	Assert(*sdk.TxResponse, *CommonArgs)
}

type Args struct {
	Network          Network         `json:"network"`
	ChainID          string          `json:"chain_id,omitempty"`
	GasPrice         string          `json:"gas_price,omitempty"`
	GasAdjustment    float64         `json:"gas_adjustment,omitempty"`
	Keybase          keyring.Keyring `json:"keybase,omitempty"`
	ChannelID        string          `json:"channel_id,omitempty"`
	Sender           sdk.AccAddress  `json:"sender,omitempty"`
	SenderName       string          `json:"sender_name"`
	SifchainReceiver sdk.AccAddress  `json:"receiver,omitempty"`
	CosmosReceiver   string          `json:"cosmos_receiver"`
	Amount           sdk.Coins       `json:""`
	TimeoutTimestamp uint64          `json:"timeout_timestamp,omitempty"`
	Fees             string          `json:"fees"`
}

type CommonArgs struct {
	DispensationName string
}

func main() {
	SetConfig(true)
	tests := []TestCase{CreateDispensationTx{}}
	cArgs := CommonArgs{}
	for _, test := range tests {
		msg, args := test.GetMsgAndArgs(cArgs)
		txf, clientCtx := getClientAndFactory(args)
		res := BroadCast(txf, clientCtx, msg)
		test.Assert(res, &cArgs)
	}
}
