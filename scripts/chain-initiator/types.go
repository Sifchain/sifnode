package main

import (
	"encoding/json"
	"time"

	admintypes "github.com/Sifchain/sifnode/x/admin/types"
	clptypes "github.com/Sifchain/sifnode/x/clp/types"
	dispensationtypes "github.com/Sifchain/sifnode/x/dispensation/types"
	epochstypes "github.com/Sifchain/sifnode/x/epochs/types"
	ethbridgetypes "github.com/Sifchain/sifnode/x/ethbridge/types"
	margintypes "github.com/Sifchain/sifnode/x/margin/types"
	oracletypes "github.com/Sifchain/sifnode/x/oracle/types"
	tokenregistrytypes "github.com/Sifchain/sifnode/x/tokenregistry/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	authz "github.com/cosmos/cosmos-sdk/x/authz"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	evidencetypes "github.com/cosmos/cosmos-sdk/x/evidence/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	transfertypes "github.com/cosmos/ibc-go/v4/modules/apps/transfer/types"
	ibcclienttypes "github.com/cosmos/ibc-go/v4/modules/core/02-client/types"
	ibcconnectiontypes "github.com/cosmos/ibc-go/v4/modules/core/03-connection/types"
	ibcchanneltypes "github.com/cosmos/ibc-go/v4/modules/core/04-channel/types"
	ibctypes "github.com/cosmos/ibc-go/v4/modules/core/types"
	// genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
)

type Genesis struct {
	GenesisTime     time.Time       `json:"genesis_time"`
	ChainID         string          `json:"chain_id"`
	InitialHeight   string          `json:"initial_height"`
	ConsensusParams ConsensusParams `json:"consensus_params"`
	AppHash         string          `json:"app_hash"`
	AppState        AppState        `json:"app_state"`
	// Include other top-level fields as needed
}

type ConsensusParams struct {
	Version   Version   `json:"version"`
	Block     Block     `json:"block"`
	Evidence  Evidence  `json:"evidence"`
	Validator Validator `json:"validator"`
}

type Version struct{}

type Validator struct {
	PubKeyTypes []string `json:"pub_key_types"`
}

type Evidence struct {
	MaxAgeNumBlocks string `json:"max_age_num_blocks"`
	MaxAgeDuration  string `json:"max_age_duration"`
	MaxBytes        string `json:"max_bytes,omitempty"`
}

type Block struct {
	MaxBytes   string `json:"max_bytes"`
	MaxGas     string `json:"max_gas"`
	TimeIotaMs string `json:"time_iota_ms"`
}

type AppState struct {
	Admin         Admin                       `json:"admin"`
	Auth          Auth                        `json:"auth"`
	AuthZ         authz.GenesisState          `json:"authz"`
	Bank          banktypes.GenesisState      `json:"bank"`
	Capability    Capability                  `json:"capability"`
	CLP           CLP                         `json:"clp"`
	Crisis        crisistypes.GenesisState    `json:"crisis"`
	Dispensation  Dispensation                `json:"dispensation"`
	Distribution  Distribution                `json:"distribution"`
	Epochs        Epochs                      `json:"epochs"`
	Ethbridge     ethbridgetypes.GenesisState `json:"ethbridge"`
	Evidence      EvidenceState               `json:"evidence"`
	Genutil       Genutil                     `json:"genutil"`
	Gov           Gov                         `json:"gov"`
	Ibc           Ibc                         `json:"ibc"`
	Margin        Margin                      `json:"margin"`
	Mint          Mint                        `json:"mint"`
	Oracle        Oracle                      `json:"oracle"`
	Params        interface{}                 `json:"params"`
	Slashing      Slashing                    `json:"slashing"`
	Staking       Staking                     `json:"staking"`
	TokenRegistry TokenRegistry               `json:"tokenregistry"`
	Transfer      transfertypes.GenesisState  `json:"transfer"`
	Upgrade       struct{}                    `json:"upgrade"`
	// Include other fields as needed
}

type Epochs struct {
	epochstypes.GenesisState

	Epochs []interface{} `json:"epochs"`
}

type Genutil struct {
	// genutiltypes.GenesisState

	GenTxs []interface{} `json:"gen_txs"`
}

type Admin struct {
	admintypes.GenesisState

	AdminAccounts []AdminAccount `json:"admin_accounts"`
}

type AdminAccount struct {
	admintypes.AdminAccount

	AdminType string `json:"admin_type"`
}

type TokenRegistry struct {
	tokenregistrytypes.GenesisState

	Registry Registry `json:"registry"`
}

type Registry struct {
	tokenregistrytypes.Registry

	Entries []*RegistryEntry `json:"entries"`
}

type RegistryEntry struct {
	tokenregistrytypes.RegistryEntry

	Decimals    json.Number   `json:"decimals"`
	Permissions []interface{} `json:"permissions"`
}

type EvidenceState struct {
	evidencetypes.GenesisState

	Evidence []interface{} `json:"evidence"`
}

type Oracle struct {
	oracletypes.GenesisState

	AddressWhitelist []interface{} `json:"address_whitelist"`
	Prophecies       []interface{} `json:"prophecies"`
}

type Dispensation struct {
	dispensationtypes.GenesisState

	DistributionRecords interface{} `json:"distribution_records"`
	Distributions       interface{} `json:"distributions"`
	Claims              interface{} `json:"claims"`
}

type Capability struct {
	capabilitytypes.GenesisState

	Index  json.Number   `json:"index"`
	Owners []interface{} `json:"owners"`
}

type Slashing struct {
	slashingtypes.GenesisState

	Params       SlashingParams `json:"params"`
	SigningInfos []interface{}  `json:"signing_infos"`
	MissedBlocks []interface{}  `json:"missed_blocks"`
}

type SlashingParams struct {
	slashingtypes.Params

	SignedBlocksWindow   json.Number `json:"signed_blocks_window"`
	DowntimeJailDuration string      `json:"downtime_jail_duration"`
}

type Mint struct {
	minttypes.GenesisState

	Params MintParams `json:"params"`
}

type MintParams struct {
	minttypes.Params

	BlocksPerYear json.Number `json:"blocks_per_year"`
}

type Gov struct {
	govtypes.GenesisState

	StartingProposalId json.Number      `json:"starting_proposal_id"`
	Deposits           []interface{}    `json:"deposits"`
	Votes              []interface{}    `json:"votes"`
	Proposals          []interface{}    `json:"proposals"`
	DepositParams      GovDepositParams `json:"deposit_params"`
	VotingParams       GovVotingParams  `json:"voting_params"`
}

type GovDepositParams struct {
	govtypes.DepositParams

	MaxDepositPeriod string `json:"max_deposit_period"`
}

type GovVotingParams struct {
	govtypes.VotingParams

	VotingPeriod string `json:"voting_period"`
}

type Staking struct {
	stakingtypes.GenesisState

	Params               StakingParams `json:"params"`
	LastValidatorPowers  []interface{} `json:"last_validator_powers"`
	Validators           []interface{} `json:"validators"`
	Delegations          []interface{} `json:"delegations"`
	UnbondingDelegations []interface{} `json:"unbonding_delegations"`
	Redelegations        []interface{} `json:"redelegations"`
}

type StakingParams struct {
	stakingtypes.Params

	UnbondingTime     string      `json:"unbonding_time"`
	MaxValidators     json.Number `json:"max_validators"`
	MaxEntries        json.Number `json:"max_entries"`
	HistoricalEntries json.Number `json:"historical_entries"`
}

type Distribution struct {
	distributiontypes.GenesisState

	DelegatorWithdrawInfos          []interface{} `json:"delegator_withdraw_infos"`
	OutstandingRewards              []interface{} `json:"outstanding_rewards"`
	ValidatorAccumulatedCommissions []interface{} `json:"validator_accumulated_commissions"`
	ValidatorHistoricalRewards      []interface{} `json:"validator_historical_rewards"`
	ValidatorCurrentRewards         []interface{} `json:"validator_current_rewards"`
	DelegatorStartingInfos          []interface{} `json:"delegator_starting_infos"`
	ValidatorSlashEvents            []interface{} `json:"validator_slash_events"`
}

type Ibc struct {
	ibctypes.GenesisState

	ClientGenesis     ClientGenesis     `json:"client_genesis"`
	ConnectionGenesis ConnectionGenesis `json:"connection_genesis"`
	ChannelGenesis    ChannelGenesis    `json:"channel_genesis"`
}

type ClientGenesis struct {
	ibcclienttypes.GenesisState

	Clients            []interface{}         `json:"clients"`
	ClientsConsensus   []interface{}         `json:"clients_consensus"`
	ClientsMetadata    []interface{}         `json:"clients_metadata"`
	Params             ibcclienttypes.Params `json:"params"`
	NextClientSequence json.Number           `json:"next_client_sequence"`
}

type ConnectionGenesis struct {
	ibcconnectiontypes.GenesisState

	Connections            []interface{}           `json:"connections"`
	ClientConnectionPaths  []interface{}           `json:"client_connection_paths"`
	NextConnectionSequence json.Number             `json:"next_connection_sequence"`
	Params                 ConnectionGenesisParams `json:"params"`
}

type ConnectionGenesisParams struct {
	ibcconnectiontypes.Params

	MaxExpectedTimePerBlock json.Number `json:"max_expected_time_per_block"`
}

type ChannelGenesis struct {
	ibcchanneltypes.GenesisState

	Channels            []interface{} `json:"channels"`
	Acknowledgements    []interface{} `json:"acknowledgements"`
	Commitments         []interface{} `json:"commitments"`
	Receipts            []interface{} `json:"receipts"`
	SendSequences       []interface{} `json:"send_sequences"`
	RecvSequences       []interface{} `json:"recv_sequences"`
	AckSequences        []interface{} `json:"ack_sequences"`
	NextChannelSequence json.Number   `json:"next_channel_sequence"`
}

type CLP struct {
	clptypes.GenesisState

	Params                        CLPParams                              `json:"params"`
	PoolList                      []interface{}                          `json:"pool_list"`
	LiquidityProviders            []interface{}                          `json:"liquidity_providers"`
	RewardsBucketList             []interface{}                          `json:"rewards_bucket_list"`
	RewardParams                  CLPRewardParams                        `json:"reward_params,omitempty"`
	PmtpParams                    CLPPmtpParams                          `json:"pmtp_params,omitempty"`
	PmtpEpoch                     CLPPmtpEpoch                           `json:"pmtp_epoch,omitempty"`
	PmtpRateParams                clptypes.PmtpRateParams                `json:"pmtp_rate_params,omitempty"`
	LiquidityProtectionParams     CLPLiquidityProtectionParams           `json:"liquidity_protection_params,omitempty"`
	LiquidityProtectionRateParams clptypes.LiquidityProtectionRateParams `json:"liquidity_protection_rate_params,omitempty"`
	SwapFeeParams                 clptypes.SwapFeeParams                 `json:"swap_fee_params,omitempty"`
	ProviderDistributionParams    CLPProviderDistributionParams          `json:"provider_distribution_params,omitempty"`
}

type CLPProviderDistributionParams struct {
	clptypes.ProviderDistributionParams

	DistributionPeriods []interface{} `json:"distribution_periods"`
}

type CLPParams struct {
	clptypes.Params

	MinCreatePoolThreshold json.Number `json:"min_create_pool_threshold"`
}

type CLPRewardParams struct {
	clptypes.RewardParams

	LiquidityRemovalLockPeriod   json.Number   `json:"liquidity_removal_lock_period"`
	LiquidityRemovalCancelPeriod json.Number   `json:"liquidity_removal_cancel_period"`
	RewardPeriods                []interface{} `json:"reward_periods"`
	RewardsLockPeriod            json.Number   `json:"rewards_lock_period"`
}

type CLPPmtpParams struct {
	clptypes.PmtpParams

	PmtpPeriodEpochLength json.Number `json:"pmtp_period_epoch_length"`
	PmtpPeriodStartBlock  json.Number `json:"pmtp_period_start_block"`
	PmtpPeriodEndBlock    json.Number `json:"pmtp_period_end_block"`
}

type CLPPmtpEpoch struct {
	clptypes.PmtpEpoch

	EpochCounter json.Number `json:"epoch_counter"`
	BlockCounter json.Number `json:"block_counter"`
}

type CLPLiquidityProtectionParams struct {
	clptypes.LiquidityProtectionParams

	EpochLength json.Number `json:"epoch_length"`
}

type Margin struct {
	margintypes.GenesisState

	Params MarginParams `json:"params"`
}

type MarginParams struct {
	margintypes.Params

	EpochLength      json.Number `json:"epoch_length"`
	MaxOpenPositions json.Number `json:"max_open_positions"`
}

type AuthParams struct {
	authtypes.Params

	MaxMemoCharacters      json.Number `json:"max_memo_characters"`
	TxSigLimit             json.Number `json:"tx_sig_limit"`
	TxSizeCostPerByte      json.Number `json:"tx_size_cost_per_byte"`
	SigVerifyCostEd25519   json.Number `json:"sig_verify_cost_ed25519"`
	SigVerifyCostSecp256K1 json.Number `json:"sig_verify_cost_secp256k1"`
}

type BaseAccount struct {
	Address       string      `json:"address"`
	PubKey        interface{} `json:"pub_key"`
	AccountNumber json.Number `json:"account_number"`
	Sequence      json.Number `json:"sequence"`
}

type ModuleAccount struct {
	BaseAccount BaseAccount `json:"base_account"`
	Name        string      `json:"name"`
	Permissions []string    `json:"permissions"`
}

type Account struct {
	*BaseAccount
	*ModuleAccount

	Type string `json:"@type"`
}

type Auth struct {
	authtypes.GenesisState

	Params   AuthParams `json:"params"`
	Accounts []Account  `json:"accounts"`
}

// KeyOutput represents the JSON structure of the output from the add key command
type KeyOutput struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Address  string `json:"address"`
	PubKey   string `json:"pubkey"`
	Mnemonic string `json:"mnemonic"`
}

// StatusOutput represents the JSON structure of the output from the status command
type StatusOutput struct {
	SyncInfo struct {
		LatestBlockHeight string `json:"latest_block_height"`
	} `json:"SyncInfo"`
}

// ProposalsOutput represents the JSON structure of the output from the query proposals command
type ProposalsOutput struct {
	Proposals []struct {
		ProposalId string `json:"proposal_id"`
	} `json:"proposals"`
}

// Colors
const (
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
)
