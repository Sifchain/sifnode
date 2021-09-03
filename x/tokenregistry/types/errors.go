package types

import sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

var (
	ErrNotFound         = sdkerrors.Register(ModuleName, 1, "denom not found in registry")
	ErrPermissionDenied = sdkerrors.Register(ModuleName, 2, "permission denied for denom")
)
