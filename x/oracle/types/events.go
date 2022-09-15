package types

// Oracle module event types
var (
	EventTypeProphecyCompleted = "prophecy_completed"

	AttributeKeyAmount               = "amount"
	AttributeKeyCrossChainFee        = "cross_chain_fee_amount"
	AttributeKeyTokenContractAddress = "token_contract_address"
	AttributeKeyCosmosSender         = "cosmos_sender"
	AttributeKeyCosmosSenderSequence = "cosmos_sender_sequence"
	AttributeKeyEthereumReceiver     = "ethereum_receiver"
	AttributeKeyNetworkDescriptor    = "network_id"
	AttributeKeyProphecyID           = "prophecy_id"
	AttributeKeyTokenAddress         = "token_address"
	AttributeKeyBridgeToken          = "bridge_token"
	AttributeKeyGlobalNonce          = "global_sequence"
	AttributeKeyEthereumAddresses    = "ethereum_addresses"

	AttributeKeySignatures = "ethereum_signatures"

	AttributeValueCategory = ModuleName
)
