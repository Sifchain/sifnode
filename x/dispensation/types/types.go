package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"strings"
)

type Airdrop struct {
	NativeAssetBalance   sdk.Uint `json:"native_asset_balance"`
	ExternalAssetBalance sdk.Uint `json:"external_asset_balance"`
	PoolUnits            sdk.Uint `json:"pool_units"`
}

func (p Airdrop) String() string {
	return strings.TrimSpace(fmt.Sprintf(`ExternalAsset: %s
	ExternalAssetBalance: %s
	NativeAssetBalance: %s
	PoolUnits : %s`, p.ExternalAssetBalance, p.NativeAssetBalance, p.PoolUnits))
}

func (p Airdrop) Validate() bool {
	return true
}

// NewPool returns a new Pool
func NewAirdrop( nativeAssetBalance, externalAssetBalance, poolUnits sdk.Uint) (Airdrop, error) {
	ad := Airdrop{
		NativeAssetBalance:   nativeAssetBalance,
		ExternalAssetBalance: externalAssetBalance,
		PoolUnits:            poolUnits}

	return ad, nil
}


// ----------------------------------------------------------------------------
// Client Types

type AirdropResponse struct {
	Airdrop
	Height           int64  `json:"height"`
}

func NewAirdropResponse(ad Airdrop, height int64) AirdropResponse {
	return AirdropResponse{Airdrop: ad, Height: height}
}
