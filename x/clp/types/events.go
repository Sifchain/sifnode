package types

// clp module event types

const (
	EventTypeCreatePool              = "created_new_pool"
	EventTypeDecommissionPool        = "decommission_pool"
	EventTypeCreateLiquidityProvider = "created_new_liquidity_provider"
	EventTypeAddLiquidity            = "added_liquidity"
	EventTypeRemoveLiquidity         = "removed_liquidity"
	EventTypeSwap                    = "swap_successful"
	EventTypeSwapFailed              = "swap_failed"
	AttributeKeyThreshold            = "min_threshold"
	AttributeKeySwapAmount           = "swap_amount"
	AttributeKeyLiquidityFee         = "liquidity_fee"
	AttributeKeyPriceImpact          = "price_impact"
	AttributeKeyInPool               = "in_pool"
	AttributeKeyOutPool              = "out_pool"
	AttributeKeyPool                 = "pool"
	AttributeKeyHeight               = "height"
	AttributeKeyLiquidityProvider    = "liquidity_provider"
	AttributeValueCategory           = ModuleName
)
