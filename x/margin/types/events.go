package types

const EventOpen = "margin/mtp_open"
const EventClose = "margin/mtp_close"
const EventForceClose = "margin/mtp_force_close"
const EventAdminClose = "margin/mtp_admin_close"
const EventAdminCloseAll = "margin/mtp_admin_close_all"
const EventInterestRateComputation = "margin/interest_rate_computation"
const EventMarginUpdateParams = "margin/update_params"
const EventRepayFund = "margin/repay_fund"
const EventBelowRemovalThreshold = "margin/below_removal_threshold"
const EventAboveRemovalThreshold = "margin/above_removal_threshold"
const EventIncrementalPayFund = "margin/incremental_pay_fund"

const AttributeKeyPoolInterestRate = "margin_pool_interest_rate"
const AttributeKeyMarginParams = "margin_params"
