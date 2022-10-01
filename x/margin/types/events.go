package types

const (
	EventOpen                    = "margin/mtp_open"
	EventClose                   = "margin/mtp_close"
	EventForceClose              = "margin/mtp_force_close"
	EventAdminClose              = "margin/mtp_admin_close"
	EventAdminCloseAll           = "margin/mtp_admin_close_all"
	EventInterestRateComputation = "margin/interest_rate_computation"
	EventMarginUpdateParams      = "margin/update_params"
	EventRepayFund               = "margin/repay_fund"
	EventBelowRemovalThreshold   = "margin/below_removal_threshold"
	EventAboveRemovalThreshold   = "margin/above_removal_threshold"
	EventIncrementalPayFund      = "margin/incremental_pay_fund"
)

const (
	AttributeKeyPoolInterestRate = "margin_pool_interest_rate"
	AttributeKeyMarginParams     = "margin_params"
)
