package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/supply"
	"strings"
)

type Asset struct {
	SourceChain string `json:"source_chain"`
	Symbol      string `json:"symbol"`
	Ticker      string `json:"ticker"`
}

type Assets []Asset

// NewAsset returns a new Asset
func NewAsset(sourceChain string, symbol string, ticker string) Asset {
	return Asset{
		SourceChain: sourceChain,
		Symbol:      symbol,
		Ticker:      ticker,
	}
}

// implement fmt.Stringer
func (a Asset) String() string {
	return strings.TrimSpace(fmt.Sprintf(`SourceChain: %s
Symbol: %s
Ticker: %s`, a.SourceChain, a.Symbol, a.Ticker))
}

func (a Asset) Validate() bool {
	if !VerifyRange(len(strings.TrimSpace(a.SourceChain)), 0, MaxSourceChainLength) {
		return false
	}
	if !VerifyRange(len(strings.TrimSpace(a.Symbol)), 0, MaxSymbolLength) {
		return false
	}
	if !VerifyRange(len(strings.TrimSpace(a.Ticker)), 0, MaxTickerLength) {
		return false
	}
	coin := sdk.NewCoin(a.Ticker, sdk.OneInt())
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
	return a.SourceChain == (a2.SourceChain) && a.Symbol == (a2.Symbol) && a.Ticker == (a2.Ticker)
}

func (a Asset) IsEmpty() bool {
	return a.SourceChain == "" || a.Symbol == "" || a.Ticker == ""
}

func GetSettlementAsset() Asset {
	return Asset{
		SourceChain: NativeChain,
		Symbol:      NativeSymbol,
		Ticker:      NativeTicker,
	}

}

func GetCLPModuleAddress() sdk.AccAddress {
	return supply.NewModuleAddress(ModuleName)
}
