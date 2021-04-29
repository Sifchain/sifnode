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
		wl := make([]sdk.ValAddress, len(data.AddressWhitelist))
		for i, entry := range data.AddressWhitelist {
			wlAddress, err := sdk.ValAddressFromBech32(entry)
			if err != nil {
				panic(err)
			}
			wl[i] = wlAddress
		}

		keeper.SetOracleWhiteList(ctx, wl)
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
	whiteList := keeper.GetOracleWhiteList(ctx)
	wl := make([]string, len(whiteList), 0)
	for i, entry := range whiteList {
		wl[i] = entry.String()
	}
	return &types.GenesisState{
		AddressWhitelist: wl,
	}
}

// ValidateGenesis validates the oracle genesis parameters
func ValidateGenesis(data *types.GenesisState) error {
	return nil
}
