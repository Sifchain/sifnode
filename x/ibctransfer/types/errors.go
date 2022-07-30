package types

import sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

var (
	ErrConvertingToUnitDenom         = sdkerrors.Register(ModuleName, 1, "error converting to unit denom")
	ErrConvertingToCounterpartyDenom = sdkerrors.Register(ModuleName, 2, "error converting to counterparty denom")
	ErrAmountTooLowToConvert         = sdkerrors.Register(ModuleName, 3, "amount too low to convert to counterparty denom")
	ErrAmountTooLargeToSend          = sdkerrors.Register(ModuleName, 4, "amount too large to transfer")
)
