package types

import (
	_ "github.com/cosmos/cosmos-sdk/types/errors"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrInvalid            = sdkerrors.Register(ModuleName, 1, "Invalid Input")
	NotEnoughBalance      = sdkerrors.Register(ModuleName, 2, "Faucet does not have enough balance ")
	ErrorRequestingTokens = sdkerrors.Register(ModuleName, 3, "Faucet cannot fund the specified address")
	ErrorAddingTokens     = sdkerrors.Register(ModuleName, 4, "Unable to fund faucet")
)
