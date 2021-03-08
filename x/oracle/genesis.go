package oracle

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/Sifchain/sifnode/x/oracle/keeper"
	"github.com/Sifchain/sifnode/x/oracle/types"
)

func InitGenesis(ctx sdk.Context, keeper keeper.Keeper, data types.GenesisState) (res []abci.ValidatorUpdate) {

	if data.AddressWhitelist != nil {
		keeper.SetOracleWhiteList(ctx, data.AddressWhitelist)
	}

	if data.AdminAddress != nil {
		keeper.SetAdminAccount(ctx, data.AdminAddress)
	}

	return []abci.ValidatorUpdate{}
}

func ExportGenesis(ctx sdk.Context, keeper keeper.Keeper) types.GenesisState {
	whiteList := keeper.GetOracleWhiteList(ctx)
	return types.GenesisState{
		AddressWhitelist: whiteList,
	}
}

// ValidateGenesis validates the oracle genesis parameters
func ValidateGenesis(data *types.GenesisState) error {
	return nil
}
