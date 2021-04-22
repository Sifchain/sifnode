package types

import (
	"time"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	disttypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	tmtypes "github.com/tendermint/tendermint/proto/tendermint/types"

	clptypes "github.com/Sifchain/sifnode/x/clp/types"
)

const (
	StakeTokenDenom = "rowan"
)

type GentxValueSignaturePubKey struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type GentxValueSignature struct {
	PubKey    GentxValueSignaturePubKey `json:"pub_key"`
	Signature string                    `json:"signature"`
}

type GentxValueFee struct {
	Amount []interface{} `json:"amount"`
	Gas    string        `json:"gas"`
}

type GentxValueMsgValueValue struct {
	Denom  string `json:"denom"`
	Amount string `json:"amount"`
}

type GentxValueMsgValueCommission struct {
	Rate          string `json:"rate"`
	MaxRate       string `json:"max_rate"`
	MaxChangeRate string `json:"max_change_rate"`
}

type GentxValueMsgValueDescription struct {
	Moniker         string `json:"moniker"`
	Identity        string `json:"identity"`
	Website         string `json:"website"`
	SecurityContact string `json:"security_contact"`
	Details         string `json:"details"`
}

type GentxValueMsgValue struct {
	Description       GentxValueMsgValueDescription `json:"description"`
	Commission        GentxValueMsgValueCommission  `json:"commission"`
	MinSelfDelegation string                        `json:"min_self_delegation"`
	DelegatorAddress  string                        `json:"delegator_address"`
	ValidatorAddress  string                        `json:"validator_address"`
	Pubkey            string                        `json:"pubkey"`
	Value             GentxValueMsgValueValue       `json:"value"`
}

type GentxValueMsg struct {
	Type  string             `json:"type"`
	Value GentxValueMsgValue `json:"value"`
}

type GentxValue struct {
	Msg        []GentxValueMsg       `json:"msg"`
	Fee        GentxValueFee         `json:"fee"`
	Signatures []GentxValueSignature `json:"signatures"`
	Memo       string                `json:"memo"`
}

type Gentx struct {
	Type  string     `json:"type"`
	Value GentxValue `json:"value"`
}

type Genutil struct {
	Gentxs []Gentx `json:"gentxs"`
}

type Upgrade struct{}

type AppState struct {
	Auth         authtypes.GenesisState     `json:"auth"`
	Bank         banktypes.GenesisState     `json:"bank"`
	Staking      stakingtypes.GenesisState  `json:"staking"`
	Params       interface{}                `json:"params"`
	Ethbridge    interface{}                `json:"ethbridge"`
	CLP          clptypes.GenesisState      `json:"clp"`
	Oracle       interface{}                `json:"oracle"`
	Genutil      genutiltypes.GenesisState  `json:"genutil"`
	Gov          govtypes.GenesisState      `json:"gov"`
	Slashing     slashingtypes.GenesisState `json:"slashing"`
	Distribution disttypes.GenesisState     `json:"distribution"`
	Dispensation interface{}                `json:"dispensation"`
}

type Genesis struct {
	GenesisTime     time.Time               `json:"genesis_time"`
	ChainID         string                  `json:"chain_id"`
	ConsensusParams tmtypes.ConsensusParams `json:"consensus_params"`
	AppHash         string                  `json:"app_hash"`
	AppState        AppState                `json:"app_state"`
}
