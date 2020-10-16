package types

import (
	"fmt"
	"strings"
)

type Asset struct {
	SourceChain string `json:"source_chain"`
	Symbol      string `json:"symbol"`
	Ticker      string `json:"ticker"`
}

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
	return strings.TrimSpace(fmt.Sprintf(`SourceChain: %s Symbol: %s Ticker: %s`, a.SourceChain, a.Symbol, a.Ticker))
}

func (a Asset) Validate() bool {
	if len(strings.TrimSpace(a.SourceChain)) == 0 {
		return false
	}
	if a.SourceChain == a.Ticker {
		return false
	}
	if len(strings.TrimSpace(a.Symbol)) == 0 {
		return false
	}
	if len(strings.TrimSpace(a.Ticker)) == 0 {
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
