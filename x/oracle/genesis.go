package oracle

import (
	"github.com/Sifchain/sifnode/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func InitGenesis(ctx sdk.Context, keeper Keeper, data types.GenesisState) (res []abci.ValidatorUpdate) {

	keeper.SetOracleWhiteList(ctx, data.AddressWhitelist)

	keeper.SetAdminAccount(ctx, data.AdminAddress)

	for _, p := range data.Prophecies {
		keeper.SetDBProphecy(ctx, p)
	}

	return []abci.ValidatorUpdate{}
}

func ExportGenesis(ctx sdk.Context, keeper Keeper) types.GenesisState {
	whiteList := keeper.GetOracleWhiteList(ctx)
	adminAddress := keeper.GetAdminAccount(ctx)
	prophecies := keeper.GetProphecies(ctx)

	dbProphecies := make([]types.DBProphecy, len(prophecies))
	for i, p := range prophecies {
		dbP, err := p.SerializeForDB()
		if err != nil {
			panic(err)
		}
		dbProphecies[i] = dbP
	}

	return GenesisState{
		AddressWhitelist: whiteList,
		AdminAddress:     adminAddress,
		Prophecies:       dbProphecies,
	}
}
