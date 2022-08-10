package types

// Ethbridge module event types
const (
	EventTypeCreateClaim              = "sif/" + ModuleName + "/create_claim"
	EventTypeProphecyStatus           = "sif/" + ModuleName + "/prophecy_status"
	EventTypeBurn                     = "sif/" + ModuleName + "/sif_burn"
	EventTypeLock                     = "sif/" + ModuleName + "/sif_lock"
	EventTypeUpdateWhiteListValidator = "sif/" + ModuleName + "/update_whitelist_validator"
	EventTypeSetCrossChainFee         = "sif/" + ModuleName + "/set_cross_chain_fee"
	EventTypeSignProphecy             = "sif/" + ModuleName + "/sign_prophecy"
	EventTypeUpdateConsensusNeeded    = "sif/" + ModuleName + "/update_consensus_needed"

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
