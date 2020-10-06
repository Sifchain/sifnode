package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// TODO: Fill out some custom errors for the module
// You can see how they are constructed below:
var (
	ErrInvalid                    = sdkerrors.Register(ModuleName, 1, "Invalid")
	ErrPoolDoesNotExist           = sdkerrors.Register(ModuleName, 1, "Pool does not exists")
	LiquidityProviderDoesNotExist = sdkerrors.Register(ModuleName, 1, "Liquidity Provider does not exists")
)
