package types

// Oracle module event types
var (
	EventTypeProphecyCompleted = "prophecy_completed"

	AttributeKeyAmount               = "amount"
	AttributeKeycrossChainFee        = "cross_chain_fee_amount"
	AttributeKeyTokenContractAddress = "token_contract_address"
	AttributeKeyCosmosSender         = "cosmos_sender"
	AttributeKeyCosmosSenderSequence = "cosmos_sender_sequence"
	AttributeKeyEthereumReceiver     = "ethereum_receiver"
	AttributeKeyNetworkDescriptor    = "network_id"
	AttributeKeyProphecyID           = "prophecy_id"
	AttributeKeyTokenAddress         = "token_address"
	AttributeKeyDoublePeggy          = "double_peggy"
	AttributeKeyGlobalNonce          = "global_nonce"
	AttributeKeyEthereumAddresses    = "ethereum_addresses"

	AttributeKeySignatures = "ethereum_signatures"

	AttributeValueCategory = ModuleName
)
