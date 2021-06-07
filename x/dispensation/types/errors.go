package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrInvalid       = sdkerrors.Register(ModuleName, 1, "invalid")
	ErrFailedInputs  = sdkerrors.Register(ModuleName, 4, "Failed in collecting funds for airdrop")
	ErrFailedOutputs = sdkerrors.Register(ModuleName, 5, "Failed in distributing funds for airdrop")
	ErrDistribution  = sdkerrors.Register(ModuleName, 6, "DistributionFailed")
)
