package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// TODO: Fill out some custom errors for the module
// You can see how they are constructed below:
var (
	ErrBalanceTooHigh                = sdkerrors.Register(ModuleName, 11, "Pool Balance too high to be decommissioned")
	ErrInvalid                       = sdkerrors.Register(ModuleName, 1, "invalid")
	ErrPoolDoesNotExist              = sdkerrors.Register(ModuleName, 2, "pool does not exist")
	ErrLiquidityProviderDoesNotExist = sdkerrors.Register(ModuleName, 3, "liquidity Provider does not exist")
	ErrInValidAsset                  = sdkerrors.Register(ModuleName, 4, "asset is invalid")
	ErrInValidAmount                 = sdkerrors.Register(ModuleName, 5, "amount is invalid")
	ErrPoolListIsEmpty               = sdkerrors.Register(ModuleName, 6, "poolList Is Empty")
	ErrTotalAmountTooLow             = sdkerrors.Register(ModuleName, 7, "total amount is less than minimum threshold")
	ErrNotEnoughAssetTokens          = sdkerrors.Register(ModuleName, 8, "not enough received asset tokens to swap")
	ErrInvalidAsymmetry              = sdkerrors.Register(ModuleName, 9, "Asymmetry has to be 1,-1 or 0")
	ErrInvalidWBasis                 = sdkerrors.Register(ModuleName, 10, "WBasisPoints has to be positive")
)
