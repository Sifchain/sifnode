package types

import (
	"github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrInvalid            = errors.Register(ModuleName, 1, "Invalid Input")
	NotEnoughBalance      = errors.Register(ModuleName, 2, "Faucet does not have enough balance ")
	ErrorRequestingTokens = errors.Register(ModuleName, 3, "Faucet cannot fund the specified address")
	ErrorAddingTokens     = errors.Register(ModuleName, 4, "Unable to fund faucet")
)
