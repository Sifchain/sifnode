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
		for networkID, list := range data.AddressWhitelist {
			keeper.SetOracleWhiteList(ctx, types.NewNetworkDescriptor(networkID), *list)
		}

	}

	if len(strings.TrimSpace(data.AdminAddress)) != 0 {
		adminAddress, err := sdk.AccAddressFromBech32(data.AdminAddress)
		if err != nil {
			panic(err)
		}
		keeper.SetAdminAccount(ctx, adminAddress)
	}

	return []abci.ValidatorUpdate{}
}

func ExportGenesis(ctx sdk.Context, keeper keeper.Keeper) *types.GenesisState {
	whiteList := keeper.GetAllWhiteList(ctx)
	wl := make(map[uint32]*types.ValidatorWhiteList)

	for i, value := range whiteList {
		wl[i] = &types.ValidatorWhiteList{WhiteList: value.WhiteList}
	}
	return &types.GenesisState{
		AddressWhitelist: wl,
		AdminAddress:     keeper.GetAdminAccount(ctx).String(),
	}
}

// ValidateGenesis validates the oracle genesis parameters
func ValidateGenesis(_ *types.GenesisState) error {
	return nil
}
