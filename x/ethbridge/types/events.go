package types

// Ethbridge module event types
const (
	EventTypeCreateClaim                        = "create_claim"
	EventTypeProphecyStatus                     = "prophecy_status"
	EventTypeBurn                               = "sif_burn"
	EventTypeLock                               = "sif_lock"
	EventTypeUpdateWhiteListValidator           = "update_whitelist_validator"
	EventTypeSetCrossChainFee                   = "set_cross_chain_fee"
	EventTypeSignProphecy                       = "sign_prophecy"
	EventTypeUpdateConsensusNeeded              = "update_consensus_needed"
	EventTypeUpdateCrossChainFeeReceiverAccount = "update_cross_chain_fee_receiver_account"

	AttributeKeyEthereumSender               = "ethereum_sender"
	AttributeKeyEthereumSenderSequence       = "ethereum_sender_sequence"
	AttributeKeyCosmosReceiver               = "cosmos_receiver"
	AttributeKeyAmount                       = "amount"
	AttributeKeycrossChainFee                = "cross_chain_fee_amount"
	AttributeKeySymbol                       = "symbol"
	AttributeKeyCoins                        = "coins"
	AttributeKeyStatus                       = "status"
	AttributeKeyClaimType                    = "claim_type"
	AttributeKeyValidator                    = "validator"
	AttributeKeyPowerType                    = "power"
	AttributeKeyCrossChainFeeReceiverAccount = "cross_chain_fee_receiver_account"

	AttributeKeyTokenContract         = "token_contract_address"
	AttributeKeyCosmosSender          = "cosmos_sender"
	AttributeKeyCosmosSenderSequence  = "cosmos_sender_sequence"
	AttributeKeyEthereumReceiver      = "ethereum_receiver"
	AttributeKeyNetworkDescriptor     = "network_id"
	AttributeKeyCrossChainFee         = "cross_chain_fee"
	AttributeKeyCrossChainFeeGas      = "cross_chain_fee_gas"
	AttributeKeyMinimumLockCost       = "minimum_lock_cost"
	AttributeKeyMinimumBurnCost       = "minimum_burn_cost"
	AttributeKeyProphecyID            = "prophecy_id"
	AttributeKeyGlobalSequence        = "global_sequence"
	AttributeKeyUpdateConsensusNeeded = "consensus_needed"

	AttributeValueCategory = ModuleName
)
