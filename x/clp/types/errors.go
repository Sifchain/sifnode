package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// TODO: Fill out some custom errors for the module
// You can see how they are constructed below:
var (
	ErrInvalid                    = sdkerrors.Register(ModuleName, 1, "Invalid")
	ErrPoolDoesNotExist           = sdkerrors.Register(ModuleName, 2, "Pool does not exist")
	LiquidityProviderDoesNotExist = sdkerrors.Register(ModuleName, 3, "Liquidity Provider does not exist")
	InValidAsset                  = sdkerrors.Register(ModuleName, 4, "Asset is invalid")
	InValidAmount                 = sdkerrors.Register(ModuleName, 5, "Amount is invalid")
	PoolListIsEmpty               = sdkerrors.Register(ModuleName, 6, "PoolList Is Empty")
	TotalAmountTooLow             = sdkerrors.Register(ModuleName, 7, "Total amount is less than minimum threshold")
)
