package v39

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const ModuleName = "clp"

type GenesisState struct {
	Params                Params             `json:"params" yaml:"params"`
	AddressWhitelist      []sdk.AccAddress   `json:"address_whitelist"`
	PoolList              Pools              `json:"pool_list"`
	LiquidityProviderList LiquidityProviders `json:"liquidity_provider_list"`
}

// Params - used for initializing default parameter for clp at genesis
type Params struct {
	MinCreatePoolThreshold uint `json:"min_create_pool_threshold"`
}

type Pools []Pool
type LiquidityProviders []LiquidityProvider

type Pool struct {
	ExternalAsset        Asset    `json:"external_asset"`
	NativeAssetBalance   sdk.Uint `json:"native_asset_balance"`
	ExternalAssetBalance sdk.Uint `json:"external_asset_balance"`
	PoolUnits            sdk.Uint `json:"pool_units"`
}

type LiquidityProvider struct {
	Asset                    Asset          `json:"asset"`
	LiquidityProviderUnits   sdk.Uint       `json:"liquidity_provider_units"`
	LiquidityProviderAddress sdk.AccAddress `json:"liquidity_provider_address"`
}

type Asset struct {
	Symbol string `json:"symbol"`
}

func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cryptocodec.RegisterCrypto(cdc)
}
