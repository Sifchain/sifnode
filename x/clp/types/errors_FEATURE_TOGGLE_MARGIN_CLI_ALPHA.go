//go:build FEATURE_TOGGLE_MARGIN_CLI_ALPHA
// +build FEATURE_TOGGLE_MARGIN_CLI_ALPHA

package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrQueued = sdkerrors.Register(ModuleName, 37, "Cannot process immediately, request has been queued")
)
