package types

import sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

var (
	ErrMTPDoesNotExist             = sdkerrors.Register(ModuleName, 1, "mtp not found")
	ErrMTPInvalid                  = sdkerrors.Register(ModuleName, 2, "mtp invalid")
	ErrMTPDisabled                 = sdkerrors.Register(ModuleName, 3, "margin not enabled for pool")
	ErrUnknownRequest              = sdkerrors.Register(ModuleName, 4, "unknown request")
	ErrMTPHealthy                  = sdkerrors.Register(ModuleName, 5, "mtp health above force close threshold")
	ErrInvalidPosition             = sdkerrors.Register(ModuleName, 6, "mtp position invalid")
	ErrMaxOpenPositions            = sdkerrors.Register(ModuleName, 7, "max open positions reached")
	ErrUnauthorised                = sdkerrors.Register(ModuleName, 8, "address not on whitelist")
	ErrBorrowTooLow                = sdkerrors.Register(ModuleName, 9, "borrowed amount is too low")
	ErrBorrowTooHigh               = sdkerrors.Register(ModuleName, 10, "borrowed amount is higher than pool depth")
	ErrCustodyTooHigh              = sdkerrors.Register(ModuleName, 11, "custody amount is higher than pool depth")
	ErrMTPUnhealthy                = sdkerrors.Register(ModuleName, 12, "mtp health would be too low for safety factor")
	ErrRowanAsCollateralNotAllowed = sdkerrors.Register(ModuleName, 13, "using rowan as collateral asset is not allowed")
)
