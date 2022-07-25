//go:build FEATURE_TOGGLE_MARGIN_CLI_ALPHA
// +build FEATURE_TOGGLE_MARGIN_CLI_ALPHA

package types

const EventOpen = "mtp_open"
const EventClose = "mtp_close"
const EventForceClose = "mtp_force_close"
const EventInterestRateComputation = "margin_interest_rate_computation"
const EventMarginUpdateParams = "margin_update_params"
const EventRepayInsuranceFund = "repay_insurance_fund"
const AttributeKeyPoolInterestRate = "pool_interest_rate"
const AttributeKeyMarginParams = "margin_params"
