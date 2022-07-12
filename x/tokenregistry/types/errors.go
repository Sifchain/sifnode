package types

import sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

var (
	ErrNotFound              = sdkerrors.Register(ModuleName, 1, "denom not found in registry")
	ErrPermissionDenied      = sdkerrors.Register(ModuleName, 2, "permission denied for denom")
	ErrNotAllowedToSellAsset = sdkerrors.Register(ModuleName, 3, "Unable to swap, not allowed to sell selected asset")
	ErrNotAllowedToBuyAsset  = sdkerrors.Register(ModuleName, 4, "Unable to swap, not allowed to buy selected asset")
	ErrInvalidMetadataName   = sdkerrors.Register(ModuleName, 5, "name in token metadata is not valid")
	ErrInvalidMetadataSymbol = sdkerrors.Register(ModuleName, 6, "symbol in token metadata is not valid")
)
