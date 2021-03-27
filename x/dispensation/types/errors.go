package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrInvalid      = sdkerrors.Register(ModuleName, 1, "invalid")
	ErrFileNotFound = sdkerrors.Register(ModuleName, 2, "File not found")
	ErrKeyInvalid   = sdkerrors.Register(ModuleName, 3, "Address in input list is not part of multi sig key")
)
