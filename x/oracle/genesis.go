package oracle

import (
	"strconv"

	"github.com/Sifchain/sifnode/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func InitGenesis(ctx sdk.Context, keeper Keeper, data types.GenesisState) (res []abci.ValidatorUpdate) {

	if data.AddressWhitelist != nil {
		for k, v := range data.AddressWhitelist {
			networkID, err := strconv.ParseUint(k, 10, 32)
			if err != nil {
				panic("white list can't parse from genesis data")
			}
			keeper.SetOracleWhiteList(ctx, NewNetworkDescriptor(uint32(networkID)), types.NewValidatorWhitelistFromData(v))
		}

	}

	if data.AdminAddress != nil {
		keeper.SetAdminAccount(ctx, data.AdminAddress)
	}

	return []abci.ValidatorUpdate{}
}

func ExportGenesis(ctx sdk.Context, keeper Keeper) types.GenesisState {
	whiteList := make(map[string]map[string]uint32)
	var i uint32 = 0
	for ; i < MaxNetworkDescriptor; i++ {
		whiteList[strconv.Itoa(int(i))] = keeper.GetOracleWhiteList(ctx, NewNetworkDescriptor(i)).Whitelist
	}

	return GenesisState{
		AddressWhitelist: whiteList,
		AdminAddress:     keeper.GetAdminAccount(ctx),
	}
}

// ValidateGenesis validates the oracle genesis parameters
func ValidateGenesis(data GenesisState) error {
	return nil
}
