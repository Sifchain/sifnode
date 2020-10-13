package types

// clp module event types

const (
	EventTypeCreatePool              = "created_new_pool"
	EventTypeCreateLiquidityProvider = "created_new_liquidity_provider"
	EventTypeAddLiquidity            = "added_liqudity"
	EventTypeRemoveLiquidity         = "removed_liquidity"
	AttributeKeyPool                 = "pool"
	AttributeKeyHeight               = "height"
	AttributeKeyLiquidityProvider    = "liquidity_provider"
	AttributeValueCategory           = ModuleName
)
