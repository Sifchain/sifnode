package keeper

import (
	"fmt"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/Sifchain/sifnode/x/clp/types"
)

// nat wraps sdk.Uint which are >= 1
type nat struct {
	i *sdk.Uint
}

type Nat interface {
}

func NewNat(n *sdk.Uint) (*nat, error) {
	if n.IsZero() {
		return nil, types.ErrUintIsZero
	}

	return &nat{n}, nil
}

func NewMustNat(n *sdk.Uint) *nat {
	if n.IsZero() {
		panic(fmt.Errorf("NewMustNat: Uint was 0"))
	}

	return &nat{n}
}

func (n nat) BigInt() *big.Int {
	return n.i.BigInt()
}

func (n nat) Uint() *sdk.Uint {
	return n.i
}
