package oracle

import (
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/Sifchain/sifnode/x/oracle/keeper"
	"github.com/Sifchain/sifnode/x/oracle/types"
)

func InitGenesis(ctx sdk.Context, keeper keeper.Keeper, data types.GenesisState) (res []abci.ValidatorUpdate) {

	var wl []sdk.ValAddress
	for _, addr := range data.AddressWhitelist {
		if len(strings.TrimSpace(addr)) == 0 {
			continue
		}
		wlAddress, err := sdk.ValAddressFromBech32(addr)
		if err != nil {
			panic(err)
		}
		wl = append(wl, wlAddress)
	}

	keeper.SetOracleWhiteList(ctx, wl)

	if len(strings.TrimSpace(data.AdminAddress)) != 0 {
		adminAddress, err := sdk.AccAddressFromBech32(data.AdminAddress)
		if err != nil {
			panic(err)
		}
		keeper.SetAdminAccount(ctx, adminAddress)
	}

	for _, dbProphecy := range data.Prophecies {
		keeper.SetDBProphecy(ctx, *dbProphecy)
	}

	return []abci.ValidatorUpdate{}
}

func ExportGenesis(ctx sdk.Context, keeper keeper.Keeper) *types.GenesisState {
	whiteList := keeper.GetOracleWhiteList(ctx)
	wl := make([]string, len(whiteList))
	for i, entry := range whiteList {
		wl[i] = entry.String()
	}

	adminAcc := keeper.GetAdminAccount(ctx)

	prophecies := keeper.GetProphecies(ctx)
	dbProphecies := make([]*types.DBProphecy, len(prophecies))
	for i, p := range prophecies {
		dbProphecy, err := p.SerializeForDB()
		if err != nil {
			panic(err)
		}
		dbProphecies[i] = &dbProphecy
	}

	return &types.GenesisState{
		AddressWhitelist: wl,
		AdminAddress:     adminAcc.String(),
		Prophecies:       dbProphecies,
	}
}

// ValidateGenesis validates the oracle genesis parameters
func ValidateGenesis(_ *types.GenesisState) error {
	return nil
}
