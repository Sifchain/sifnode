package types

// clp module event types

const (
	EventTypeCreatePool              = "created_new_pool"
	EventTypeDecommissionPool        = "decommission_pool"
	EventTypeCreateLiquidityProvider = "created_new_liquidity_provider"
	EventTypeAddLiquidity            = "added_liquidity"
	EventTypeRemoveLiquidity         = "removed_liquidity"
	EventTypeSwap                    = "swap"
	AttributeKeySwapAmount           = "swap_amount"
	AttributeKeyLiquidityFee         = "liquidity_fee"
	AttributeKeyTradeSlip            = "trade_slip"
	AttributeKeyPool                 = "pool"
	AttributeKeyHeight               = "height"
	AttributeKeyLiquidityProvider    = "liquidity_provider"
	AttributeValueCategory           = ModuleName
)
