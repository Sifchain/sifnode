package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/supply"
)

func GetFaucetModuleAddress() sdk.AccAddress {
	return supply.NewModuleAddress(ModuleName)
}
