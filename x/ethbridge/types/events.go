package types

// Ethbridge module event types
var (
	EventTypeCreateClaim              = "create_claim"
	EventTypeProphecyStatus           = "prophecy_status"
	EventTypeBurn                     = "burn"
	EventTypeLock                     = "lock"
	EventTypeUpdateWhiteListValidator = "update_whitelist_validator"

	AttributeKeyEthereumSender      = "ethereum_sender"
	AttributeKeyEthereumSenderNonce = "ethereum_sender_nonce"
	AttributeKeyCosmosReceiver      = "cosmos_receiver"
	AttributeKeyAmount              = "amount"
	AttributeKeyCethAmount          = "ceth_amount"
	AttributeKeySymbol              = "symbol"
	AttributeKeyCoins               = "coins"
	AttributeKeyStatus              = "status"
	AttributeKeyClaimType           = "claim_type"
	AttributeKeyValidator           = "validator"
	AttributeKeyOperationType       = "operation_type"
	AttributeKeyCethReceiverAccount = "ceth_receiver_account"

	AttributeKeyEthereumChainID      = "ethereum_chain_id"
	AttributeKeyTokenContract        = "token_contract_address"
	AttributeKeyCosmosSender         = "cosmos_sender"
	AttributeKeyCosmosSenderSequence = "cosmos_sender_sequence"
	AttributeKeyEthereumReceiver     = "ethereum_receiver"

	AttributeValueCategory = ModuleName
)
