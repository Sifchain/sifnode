package types

import (
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

type Assets []Asset

// NewAsset returns a new Asset
func NewAsset(symbol string) Asset {
	return Asset{
		Symbol: symbol,
	}
}

func (a Asset) Validate() bool {
	if !VerifyRange(len(strings.TrimSpace(a.Symbol)), 0, MaxSymbolLength) {
		return false
	}
	coin := sdk.NewCoin(a.Symbol, sdk.OneInt())
	return coin.IsValid()
}

func VerifyRange(num, low, high int) bool {
	if num >= high {
		return false
	}
	if num <= low {
		return false
	}
	return true
}

func (a Asset) Equals(a2 Asset) bool {
	return a.Symbol == (a2.Symbol)
}

func (a Asset) IsEmpty() bool {
	return a.Symbol == ""
}

func GetSettlementAsset() Asset {
	return Asset{
		Symbol: NativeSymbol,
	}

}

func GetCLPModuleAddress() sdk.AccAddress {
	return authtypes.NewModuleAddress(ModuleName)
}
