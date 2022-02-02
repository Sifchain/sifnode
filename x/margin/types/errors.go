package types

import sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

var (
	ErrMTPDoesNotExist = sdkerrors.Register(ModuleName, 1, "mtp not found")
	ErrMTPInvalid      = sdkerrors.Register(ModuleName, 2, "mtp invalid")
	ErrMTPDisabled     = sdkerrors.Register(ModuleName, 3, "margin not enabled for pool")
	ErrUnknownRequest  = sdkerrors.Register(ModuleName, 4, "unknown request")
	ErrMTPHealthy      = sdkerrors.Register(ModuleName, 5, "mtp health above force close threshold")
	ErrInvalidPosition = sdkerrors.Register(ModuleName, 6, "mtp position invalid")
)
