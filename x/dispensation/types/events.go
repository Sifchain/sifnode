package types

// dispensation module event types

const (
	AttributeValueCategory        = ModuleName
	EventTypeDistributionStarted  = "distribution_started"
	AttributeKeyFromModuleAccount = "module_account"
	EventTypeClaimCreated         = "userClaim_new"
	AttributeKeyClaimUser         = "userClaim_creator"
	AttributeKeyClaimType         = "userClaim_type"
	AttributeKeyClaimTime         = "userClaim_creationTime"
)
