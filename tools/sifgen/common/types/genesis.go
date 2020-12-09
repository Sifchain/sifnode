package types

import (
	"fmt"
	"time"
)

const (
	BondAmount      = "5000000000"
	StakeTokenDenom = "rowan"
	StakeFundAmount = "10000000000000000000000"
	SwapFundAmount  = "10000000000000000000000000000000000"
)

var (
	ToBond = fmt.Sprintf("%s%s", BondAmount, StakeTokenDenom)
	ToFund = []string{
		fmt.Sprintf("%s%s", StakeFundAmount, StakeTokenDenom),
		fmt.Sprintf("%s%s", SwapFundAmount, "clink"),
		fmt.Sprintf("%s%s", SwapFundAmount, "chot"),
		fmt.Sprintf("%s%s", SwapFundAmount, "cusdt"),
		fmt.Sprintf("%s%s", SwapFundAmount, "cusdc"),
	}
)

type AuthAccountValueCoin struct {
	Denom  string `json:"denom"`
	Amount string `json:"amount"`
}

type AuthAccountValue struct {
	Address       string                 `json:"address"`
	Coins         []AuthAccountValueCoin `json:"coins"`
	PublicKey     interface{}            `json:"public_key"`
	AccountNumber string                 `json:"account_number"`
	Sequence      string                 `json:"sequence"`
}

type AuthAccount struct {
	Type  string           `json:"type"`
	Value AuthAccountValue `json:"value"`
}

type AuthParams struct {
	MaxMemoCharacters      string `json:"max_memo_characters"`
	TxSigLimit             string `json:"tx_sig_limit"`
	TxSizeCostPerByte      string `json:"tx_size_cost_per_byte"`
	SigVerifyCostEd25519   string `json:"sig_verify_cost_ed25519"`
	SigVerifyCostSecp256K1 string `json:"sig_verify_cost_secp256k1"`
}

type Auth struct {
	Params   AuthParams    `json:"params"`
	Accounts []AuthAccount `json:"accounts"`
}

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

type CLPParams struct {
	MinCreatePoolThreshold string `json:"min_create_pool_threshold"`
}

type CLP struct {
	Params                    CLPParams   `json:"params"`
	WhiteListValidatorAddress interface{} `json:"white_list_validator_address"`
	PoolList                  interface{} `json:"pool_list"`
	LiquidityProviderList     interface{} `json:"liquidity_provider_list"`
}

type Supply struct {
	Supply []interface{} `json:"supply"`
}

type Staking struct {
	Params               StakingParams `json:"params"`
	LastTotalPower       string        `json:"last_total_power"`
	LastValidatorPowers  interface{}   `json:"last_validator_powers"`
	Validators           interface{}   `json:"validators"`
	Delegations          interface{}   `json:"delegations"`
	UnbondingDelegations interface{}   `json:"unbonding_delegations"`
	Redelegations        interface{}   `json:"redelegations"`
	Exported             bool          `json:"exported"`
}

type StakingParams struct {
	UnbondingTime     string `json:"unbonding_time"`
	MaxValidators     int    `json:"max_validators"`
	MaxEntries        int    `json:"max_entries"`
	HistoricalEntries int    `json:"historical_entries"`
	BondDenom         string `json:"bond_denom"`
}

type Bank struct {
	SendEnabled bool `json:"send_enabled"`
}

type AppState struct {
	Bank      Bank        `json:"bank"`
	Staking   Staking     `json:"staking"`
	Params    interface{} `json:"params"`
	Supply    Supply      `json:"supply"`
	Ethbridge interface{} `json:"ethbridge"`
	CLP       CLP         `json:"clp"`
	Oracle    interface{} `json:"oracle"`
	Genutil   Genutil     `json:"genutil"`
	Auth      Auth        `json:"auth"`
}

type Evidence struct {
	MaxAgeNumBlocks string `json:"max_age_num_blocks"`
	MaxAgeDuration  string `json:"max_age_duration"`
}

type Validator struct {
	PubKeyTypes []string `json:"pub_key_types"`
}

type Block struct {
	MaxBytes   string `json:"max_bytes"`
	MaxGas     string `json:"max_gas"`
	TimeIotaMs string `json:"time_iota_ms"`
}

type ConsensusParams struct {
	Block     Block     `json:"block"`
	Evidence  Evidence  `json:"evidence"`
	Validator Validator `json:"validator"`
}

type Genesis struct {
	GenesisTime     time.Time       `json:"genesis_time"`
	ChainID         string          `json:"chain_id"`
	ConsensusParams ConsensusParams `json:"consensus_params"`
	AppHash         string          `json:"app_hash"`
	AppState        AppState        `json:"app_state"`
}
