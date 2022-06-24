package types

// clp module event types

const (
	EventTypeCreatePool              = "created_new_pool"
	EventTypeDecommissionPool        = "decommission_pool"
	EventTypeAddNewPmtpPolicy        = "pmtp_new_policy"
	EventTypeEndPmtpPolicy           = "pmtp_end_policy"
	EventTypeCreateLiquidityProvider = "created_new_liquidity_provider"
	EventTypeAddLiquidity            = "added_liquidity"
	EventTypeRemoveLiquidity         = "removed_liquidity"
	EventTypeRequestUnlock           = "request_unlock_liquidity"
	EventTypeCancelUnlock            = "cancel_unlock_liquidity"
	EventTypeSwap                    = "swap_successful"
	EventTypeSwapFailed              = "swap_failed"
	EventTypeProcessedRemovalQueue   = "processed_removal_queue"
	EventTypeQueueRemovalRequest     = "queue_removal_request"
	EventTypeDequeueRemovalRequest   = "dequeue_removal_request"
	AttributeKeyThreshold            = "min_threshold"
	AttributeKeySwapAmount           = "swap_amount"
	AttributeKeyLiquidityFee         = "liquidity_fee"
	AttributeKeyPriceImpact          = "price_impact"
	AttributeKeyInPool               = "in_pool"
	AttributeKeyOutPool              = "out_pool"
	AttributePmtpBlockRate           = "pmtp_block_rate"
	AttributePmtpCurrentRunningRate  = "pmtp_current_running_rate"
	AttributeKeyPool                 = "pool"
	AttributeKeyHeight               = "height"
	AttributeKeyLiquidityProvider    = "liquidity_provider"
	AttributeKeyUnits                = "liquidity_units"
	AttributeKeyPmtpPolicyParams     = "pmtp_policy_params"
	AttributeKeyPmtpRateParams       = "pmtp_rate_params"
	AttributeValueCategory           = ModuleName
)
