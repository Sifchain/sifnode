package types

import sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

var (
	ErrMTPDoesNotExist = sdkerrors.Register(ModuleName, 1, "mtp not found")
	ErrMTPInvalid      = sdkerrors.Register(ModuleName, 2, "mtp invalid")
)
