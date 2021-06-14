package types

// dispensation module event types

const (
	AttributeValueCategory                = ModuleName
	EventTypeDistributionStarted          = "distribution_started"
	EventTypeDistributionRun              = "distribution_run"
	EventTypeDistributionRecordsList      = "distribution_record_"
	AttributeKeyFromModuleAccount         = "module_account"
	AttributeKeyDistributionName          = "distribution_name"
	AttributeKeyDistributionRunner        = "distribution_runner"
	AttributeKeyDistributionType          = "distribution_type"
	AttributeKeyDistributionRecordType    = "type"
	AttributeKeyDistributionRecordAddress = "recipient_address"
	AttributeKeyDistributionRecordAmount  = "amount"

	EventTypeClaimCreated = "userClaim_new"
	AttributeKeyClaimUser = "userClaim_creator"
	AttributeKeyClaimType = "userClaim_type"
	AttributeKeyClaimTime = "userClaim_creationTime"
)
