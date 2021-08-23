package types

import sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

var (
	ErrConvertingToUnitDenom         = sdkerrors.Register(ModuleName, 1, "error convertin to unit denom")
	ErrConvertingToCounterpartyDenom = sdkerrors.Register(ModuleName, 2, "error convertin to counterparty denom")
)
