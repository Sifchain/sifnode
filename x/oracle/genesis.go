package oracle

import (
	"github.com/Sifchain/sifnode/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func InitGenesis(ctx sdk.Context, keeper Keeper, data types.GenesisState) (res []abci.ValidatorUpdate) {

	if data.AddressWhitelist != nil {
		for k, v := range data.AddressWhitelist {
			keeper.SetOracleWhiteList(ctx, NewNetworkDescriptor(k), v)
		}

	}

	if data.AdminAddress != nil {
		keeper.SetAdminAccount(ctx, data.AdminAddress)
	}

	return []abci.ValidatorUpdate{}
}

func ExportGenesis(ctx sdk.Context, keeper Keeper) types.GenesisState {
	var whiteList map[int32][]sdk.ValAddress
	var i int32 = 0
	for ; i <= MaxNetworkDescriptor; i++ {
		whiteList[i] = keeper.GetOracleWhiteList(ctx, NewNetworkDescriptor(i))
	}

	return GenesisState{
		AddressWhitelist: whiteList,
	}
}

// ValidateGenesis validates the oracle genesis parameters
func ValidateGenesis(data GenesisState) error {
	return nil
}
