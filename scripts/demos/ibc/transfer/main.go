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
	GetMsgAndArgs() (sdk.Msg, Args)
	GetName() string
	Assert(*sdk.TxResponse)
}

type Args struct {
	Network          Network         `json:"network"`
	ChainID          string          `json:"chain_id,omitempty"`
	GasPrice         string          `json:"gas_price,omitempty"`
	GasAdjustment    float64         `json:"gas_adjustment,omitempty"`
	Keybase          keyring.Keyring `json:"keybase,omitempty"`
	ChannelId        string          `json:"channel_id,omitempty"`
	Sender           sdk.AccAddress  `json:"sender,omitempty"`
	SenderName       string          `json:"sender_name"`
	SifchainReceiver sdk.AccAddress  `json:"receiver,omitempty"`
	CosmosReceiver   string          `json:"cosmos_receiver"`
	Amount           sdk.Coins       `json:""`
	TimeoutTimestamp uint64          `json:"timeout_timestamp,omitempty"`
	Fees             string          `json:"fees"`
}

func main() {
	SetConfig(true)
	tests := []TestCase{IbcSentTx{}, SentTx{}}

	for _, test := range tests {
		msg, args := test.GetMsgAndArgs()
		txf, clientCtx := getClientAndFactory(args)
		res := BroadCast(txf, clientCtx, msg)
		test.Assert(res)
	}
}
