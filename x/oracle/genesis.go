package oracle

import (
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/Sifchain/sifnode/x/oracle/keeper"
	"github.com/Sifchain/sifnode/x/oracle/types"
)

func InitGenesis(ctx sdk.Context, keeper keeper.Keeper, data types.GenesisState) (res []abci.ValidatorUpdate) {
	if data.AddressWhitelist != nil {
		for networkDescriptor, list := range data.AddressWhitelist {
			keeper.SetOracleWhiteList(ctx, types.NewNetworkIdentity(types.NetworkDescriptor(networkDescriptor)), *list)
		}
	}

	if len(strings.TrimSpace(data.AdminAddress)) != 0 {
		adminAddress, err := sdk.AccAddressFromBech32(data.AdminAddress)
		if err != nil {
			panic(err)
		}
		keeper.SetAdminAccount(ctx, adminAddress)
	}

	for _, prophecy := range data.Prophecies {

		keeper.SetProphecy(ctx, *prophecy)
	}

	return []abci.ValidatorUpdate{}
}

func ExportGenesis(ctx sdk.Context, keeper keeper.Keeper) *types.GenesisState {
	whiteList := keeper.GetAllWhiteList(ctx)
	wl := make(map[uint32]*types.ValidatorWhiteList, len(whiteList))
	for key, entry := range whiteList {
		wl[uint32(key)] = &entry
	}
	adminAcc := keeper.GetAdminAccount(ctx)
	prophecies := keeper.GetProphecies(ctx)

	dbProphecies := make([]*types.Prophecy, len(prophecies))
	for i, dbProphecy := range prophecies {

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
