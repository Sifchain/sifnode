package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrInvalid       = sdkerrors.Register(ModuleName, 1, "invalid")
	ErrKeyInvalid    = sdkerrors.Register(ModuleName, 3, "Address in input list is not part of multi sig key")
	ErrFailedInputs  = sdkerrors.Register(ModuleName, 4, "Failed in collecting funds for airdrop")
	ErrFailedOutputs = sdkerrors.Register(ModuleName, 5, "Failed in distributing funds for airdrop")
	ErrAirdrop       = sdkerrors.Register(ModuleName, 6, "AirdropFailed")
)
