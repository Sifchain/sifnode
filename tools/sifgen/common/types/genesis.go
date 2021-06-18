package types

import (
	"encoding/json"
	"time"
)

const (
	StakeTokenDenom = "rowan"
)

type Block struct {
	MaxBytes   string `json:"max_bytes"`
	MaxGas     string `json:"max_gas"`
	TimeIotaMs string `json:"time_iota_ms"`
}

type Evidence struct {
	MaxAgeNumBlocks string `json:"max_age_num_blocks"`
	MaxAgeDuration  string `json:"max_age_duration"`
	MaxBytes        string `json:"max_bytes"`
}

type Validator struct {
	PubKeyTypes []string `json:"pub_key_types"`
}

type Version struct{}

type ConsensusParams struct {
	Version   Version   `json:"version"`
	Block     Block     `json:"block"`
	Evidence  Evidence  `json:"evidence"`
	Validator Validator `json:"validator"`
}

type Params struct {
	MaxMemoCharacters      string `json:"max_memo_characters"`
	TxSigLimit             string `json:"tx_sig_limit"`
	TxSizeCostPerByte      string `json:"tx_size_cost_per_byte"`
	SigVerifyCostEd25519   string `json:"sig_verify_cost_ed25519"`
	SigVerifyCostSecp256K1 string `json:"sig_verify_cost_secp256k1"`
}

type Accounts []struct {
	Type          string      `json:"@type"`
	Address       string      `json:"address"`
	PubKey        interface{} `json:"pub_key"`
	AccountNumber string      `json:"account_number"`
	Sequence      string      `json:"sequence"`
}

type Auth struct {
	Params   Params   `json:"params"`
	Accounts Accounts `json:"accounts"`
}

type BankParams struct {
	SendEnabled        []interface{} `json:"send_enabled"`
	DefaultSendEnabled bool          `json:"default_send_enabled"`
}

type Coins []struct {
	Denom  string `json:"denom"`
	Amount string `json:"amount"`
}

type Balances []struct {
	Address string `json:"address"`
	Coins   Coins  `json:"coins"`
}

type Supply []struct {
	Denom  string `json:"denom"`
	Amount string `json:"amount"`
}

type Bank struct {
	Params        BankParams    `json:"params"`
	Balances      Balances      `json:"balances"`
	Supply        Supply        `json:"supply"`
	DenomMetadata []interface{} `json:"denom_metadata"`
}

type Capability struct {
	Index  string        `json:"index"`
	Owners []interface{} `json:"owners"`
}

type CLPParams struct {
	MinCreatePoolThreshold string `json:"min_create_pool_threshold"`
}

type CLP struct {
	Params             CLPParams     `json:"params"`
	AddressWhitelist   []string      `json:"address_whitelist"`
	PoolList           []interface{} `json:"pool_list"`
	LiquidityProviders []interface{} `json:"liquidity_providers"`
}

type Dispensation struct {
	DistributionRecords interface{} `json:"distribution_records"`
	Distributions       interface{} `json:"distributions"`
}

type DistributionParams struct {
	CommunityTax        string `json:"community_tax"`
	BaseProposerReward  string `json:"base_proposer_reward"`
	BonusProposerReward string `json:"bonus_proposer_reward"`
	WithdrawAddrEnabled bool   `json:"withdraw_addr_enabled"`
}

type FeePool struct {
	CommunityPool []interface{} `json:"community_pool"`
}

type Distribution struct {
	Params                          DistributionParams `json:"params"`
	FeePool                         FeePool            `json:"fee_pool"`
	DelegatorWithdrawInfos          []interface{}      `json:"delegator_withdraw_infos"`
	PreviousProposer                string             `json:"previous_proposer"`
	OutstandingRewards              []interface{}      `json:"outstanding_rewards"`
	ValidatorAccumulatedCommissions []interface{}      `json:"validator_accumulated_commissions"`
	ValidatorHistoricalRewards      []interface{}      `json:"validator_historical_rewards"`
	ValidatorCurrentRewards         []interface{}      `json:"validator_current_rewards"`
	DelegatorStartingInfos          []interface{}      `json:"delegator_starting_infos"`
	ValidatorSlashEvents            []interface{}      `json:"validator_slash_events"`
}

type EvidenceState struct {
	Evidence []interface{} `json:"evidence"`
}

type Description struct {
	Moniker         string `json:"moniker"`
	Identity        string `json:"identity"`
	Website         string `json:"website"`
	SecurityContact string `json:"security_contact"`
	Details         string `json:"details"`
}

type Commission struct {
	Rate          string `json:"rate"`
	MaxRate       string `json:"max_rate"`
	MaxChangeRate string `json:"max_change_rate"`
}

type Pubkey struct {
	Type string `json:"@type"`
	Key  string `json:"key"`
}

type Value struct {
	Denom  string `json:"denom"`
	Amount string `json:"amount"`
}

type GenTxsBodyMessages []struct {
	Type              string      `json:"@type"`
	Description       Description `json:"description"`
	Commission        Commission  `json:"commission"`
	MinSelfDelegation string      `json:"min_self_delegation"`
	DelegatorAddress  string      `json:"delegator_address"`
	ValidatorAddress  string      `json:"validator_address"`
	Pubkey            Pubkey      `json:"pubkey"`
	Value             Value       `json:"value"`
}

type GenTxsBody struct {
	Messages                    GenTxsBodyMessages `json:"messages"`
	Memo                        string             `json:"memo"`
	TimeoutHeight               string             `json:"timeout_height"`
	ExtensionOptions            []interface{}      `json:"extension_options"`
	NonCriticalExtensionOptions []interface{}      `json:"non_critical_extension_options"`
}

type PublicKey struct {
	Type string `json:"@type"`
	Key  string `json:"key"`
}

type ModeInfo struct {
	Single Single `json:"single"`
}

type Single struct {
	Mode string `json:"mode"`
}

type SignerInfos []struct {
	PublicKey PublicKey `json:"public_key"`
	ModeInfo  ModeInfo  `json:"mode_info"`
	Sequence  string    `json:"sequence"`
}

type Fee struct {
	Amount   []interface{} `json:"amount"`
	GasLimit string        `json:"gas_limit"`
	Payer    string        `json:"payer"`
	Granter  string        `json:"granter"`
}

type GenTxsAuthInfo struct {
	SignerInfos SignerInfos `json:"signer_infos"`
	Fee         Fee         `json:"fee"`
}

type GenTxs []struct {
	Body       GenTxsBody     `json:"body"`
	AuthInfo   GenTxsAuthInfo `json:"auth_info"`
	Signatures []string       `json:"signatures"`
}

type Genutil struct {
	GenTxs GenTxs `json:"gen_txs"`
}

type GovMinDeposit []struct {
	Denom  string `json:"denom"`
	Amount string `json:"amount"`
}

type GovDepositParams struct {
	MinDeposit       GovMinDeposit `json:"min_deposit"`
	MaxDepositPeriod string        `json:"max_deposit_period"`
}

type GovVotingParams struct {
	VotingPeriod string `json:"voting_period"`
}

type GovTallyParams struct {
	Quorum        string `json:"quorum"`
	Threshold     string `json:"threshold"`
	VetoThreshold string `json:"veto_threshold"`
}

type Gov struct {
	StartingProposalID string           `json:"starting_proposal_id"`
	Deposits           []interface{}    `json:"deposits"`
	Votes              []interface{}    `json:"votes"`
	Proposals          []interface{}    `json:"proposals"`
	DepositParams      GovDepositParams `json:"deposit_params"`
	VotingParams       GovVotingParams  `json:"voting_params"`
	TallyParams        GovTallyParams   `json:"tally_params"`
}

type ClientGenesisParams struct {
	AllowedClients []string `json:"allowed_clients"`
}

type ClientGenesis struct {
	Clients            []interface{}       `json:"clients"`
	ClientsConsensus   []interface{}       `json:"clients_consensus"`
	ClientsMetadata    []interface{}       `json:"clients_metadata"`
	Params             ClientGenesisParams `json:"params"`
	CreateLocalhost    bool                `json:"create_localhost"`
	NextClientSequence string              `json:"next_client_sequence"`
}

type ConnectionGenesis struct {
	Connections            []interface{} `json:"connections"`
	ClientConnectionPaths  []interface{} `json:"client_connection_paths"`
	NextConnectionSequence string        `json:"next_connection_sequence"`
}

type ChannelGenesis struct {
	Channels            []interface{} `json:"channels"`
	Acknowledgements    []interface{} `json:"acknowledgements"`
	Commitments         []interface{} `json:"commitments"`
	Receipts            []interface{} `json:"receipts"`
	SendSequences       []interface{} `json:"send_sequences"`
	RecvSequences       []interface{} `json:"recv_sequences"`
	AckSequences        []interface{} `json:"ack_sequences"`
	NextChannelSequence string        `json:"next_channel_sequence"`
}

type Ibc struct {
	ClientGenesis     ClientGenesis     `json:"client_genesis"`
	ConnectionGenesis ConnectionGenesis `json:"connection_genesis"`
	ChannelGenesis    ChannelGenesis    `json:"channel_genesis"`
}

type Oracle struct {
	AddressWhitelist []interface{} `json:"address_whitelist"`
	AdminAddress     string        `json:"admin_address"`
}

type SlashingParams struct {
	SignedBlocksWindow      string `json:"signed_blocks_window"`
	MinSignedPerWindow      string `json:"min_signed_per_window"`
	DowntimeJailDuration    string `json:"downtime_jail_duration"`
	SlashFractionDoubleSign string `json:"slash_fraction_double_sign"`
	SlashFractionDowntime   string `json:"slash_fraction_downtime"`
}

type Slashing struct {
	Params       SlashingParams `json:"params"`
	SigningInfos []interface{}  `json:"signing_infos"`
	MissedBlocks []interface{}  `json:"missed_blocks"`
}

type StakingParams struct {
	UnbondingTime     string      `json:"unbonding_time"`
	MaxValidators     json.Number `json:"max_validators"`
	MaxEntries        json.Number `json:"max_entries"`
	HistoricalEntries json.Number `json:"historical_entries"`
	BondDenom         string      `json:"bond_denom"`
}

type Staking struct {
	Params               StakingParams `json:"params"`
	LastTotalPower       string        `json:"last_total_power"`
	LastValidatorPowers  []interface{} `json:"last_validator_powers"`
	Validators           []interface{} `json:"validators"`
	Delegations          []interface{} `json:"delegations"`
	UnbondingDelegations []interface{} `json:"unbonding_delegations"`
	Redelegations        []interface{} `json:"redelegations"`
	Exported             bool          `json:"exported"`
}

type TransferParams struct {
	SendEnabled    bool `json:"send_enabled"`
	ReceiveEnabled bool `json:"receive_enabled"`
}

type Transfer struct {
	PortID      string         `json:"port_id"`
	DenomTraces []interface{}  `json:"denom_traces"`
	Params      TransferParams `json:"params"`
}

type AppState struct {
	Upgrade      struct{}      `json:"upgrade"`
	Ethbridge    interface{}   `json:"ethbridge"`
	Params       interface{}   `json:"params"`
	Ibc          Ibc           `json:"ibc"`
	Distribution Distribution  `json:"distribution"`
	Staking      Staking       `json:"staking"`
	Gov          Gov           `json:"gov"`
	Slashing     Slashing      `json:"slashing"`
	Auth         Auth          `json:"auth"`
	Bank         Bank          `json:"bank"`
	CLP          CLP           `json:"clp"`
	Transfer     Transfer      `json:"transfer"`
	Capability   Capability    `json:"capability"`
	Dispensation Dispensation  `json:"dispensation"`
	Oracle       Oracle        `json:"oracle"`
	Evidence     EvidenceState `json:"evidence"`
	Genutil      Genutil       `json:"genutil"`
}

type Genesis struct {
	GenesisTime     time.Time       `json:"genesis_time"`
	ChainID         string          `json:"chain_id"`
	InitialHeight   string          `json:"initial_height"`
	ConsensusParams ConsensusParams `json:"consensus_params"`
	AppHash         string          `json:"app_hash"`
	AppState        AppState        `json:"app_state"`
}
