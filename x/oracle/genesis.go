package oracle

import (
	"github.com/Sifchain/sifnode/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func InitGenesis(ctx sdk.Context, keeper Keeper, data types.GenesisState) (res []abci.ValidatorUpdate) {

	if data.AddressWhitelist != nil {
		keeper.SetOracleWhiteList(ctx, data.AddressWhitelist)
	}

	if data.AdminAddress != nil {
		keeper.SetAdminAccount(ctx, data.AdminAddress)
	}

	if data.Prophecies != nil {
		for _, p := range data.Prophecies {
			keeper.SetProphecy(ctx, p)
		}
	}

	return []abci.ValidatorUpdate{}
}

func ExportGenesis(ctx sdk.Context, keeper Keeper) types.GenesisState {
	whiteList := keeper.GetOracleWhiteList(ctx)
	adminAddress := keeper.GetAdminAccount(ctx)
	prophecies := keeper.GetProphecies(ctx)

	return GenesisState{
		AddressWhitelist: whiteList,
		AdminAddress:     adminAddress,
		Prophecies:       prophecies,
	}
}
