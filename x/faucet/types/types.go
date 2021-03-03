package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

func GetFaucetModuleAddress() sdk.AccAddress {
	return authtypes.NewModuleAddress(ModuleName)
}
