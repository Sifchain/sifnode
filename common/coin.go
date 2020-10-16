package common

import (
	"errors"
	"fmt"
	"github.com/Sifchain/sifnode/x/clp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"sort"
	"strings"
)

type Coin struct {
	Asset  clp.Asset `json:"asset"`
	Amount sdk.Uint  `json:"amount"`
}

var NoCoin = Coin{
	Amount: sdk.ZeroUint(),
}

type Coins []Coin

// NewCoin return a new instance of Coin
func NewCoin(asset clp.Asset, amount sdk.Uint) Coin {
	return Coin{
		Asset:  asset,
		Amount: amount,
	}
}

func (c Coin) Equals(cc Coin) bool {
	if !c.Asset.Equals(cc.Asset) {
		return false
	}
	if !c.Amount.Equal(cc.Amount) {
		return false
	}
	return true
}

func (c Coin) IsEmpty() bool {
	if c.Asset.IsEmpty() {
		return true
	}
	if c.Amount.IsZero() {
		return true
	}
	return false
}

func (c Coin) Valid() error {
	if c.Asset.IsEmpty() {
		return errors.New("Denom cannot be empty")
	}
	if c.Amount.IsZero() {
		return errors.New("Amount cannot be zero")
	}

	return nil
}

func (c Coin) IsNative() bool {
	return c.Asset.SourceChain == clp.NativeChain
}

func (c Coin) String() string {
	return fmt.Sprintf("%d %s", c.Amount.Uint64(), c.Asset.String())
}

func (cs Coins) Valid() error {
	for _, coin := range cs {
		if err := coin.Valid(); err != nil {
			return err
		}
	}
	return nil
}

// Check if two lists of coins are equal to each other. Order does not matter
func (cs1 Coins) Equals(cs2 Coins) bool {
	if len(cs1) != len(cs2) {
		return false
	}
	// sort both lists
	sort.Slice(cs1[:], func(i, j int) bool {
		return cs1[i].Asset.String() < cs1[j].Asset.String()
	})
	sort.Slice(cs2[:], func(i, j int) bool {
		return cs2[i].Asset.String() < cs2[j].Asset.String()
	})

	for i := range cs1 {
		if !cs1[i].Equals(cs2[i]) {
			return false
		}
	}

	return true
}

func (cs Coins) IsEmpty() bool {
	if len(cs) == 0 {
		return true
	}
	for _, coin := range cs {
		if !coin.IsEmpty() {
			return false
		}
	}
	return true
}

func (cs Coins) String() string {
	coins := make([]string, len(cs))
	for i, c := range cs {
		coins[i] = c.String()
	}
	return strings.Join(coins, ", ")
}

// Contains check whether the given coin is in the list
func (cs Coins) Contains(c Coin) bool {
	for _, item := range cs {
		if c.Equals(item) {
			return true
		}
	}
	return false
}

// Gets a specific coin by asset. Assumes there is only one of this coin in the
// list.
func (cs Coins) GetCoin(asset clp.Asset) Coin {
	for _, item := range cs {
		if asset.Equals(item.Asset) {
			return item
		}
	}
	return NoCoin
}
