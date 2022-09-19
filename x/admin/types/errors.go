package types

import sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

var ErrPermissionDenied = sdkerrors.Register(ModuleName, 1, "permission denied")
