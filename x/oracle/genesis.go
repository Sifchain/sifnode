package oracle

import (
	"log"

	"github.com/Sifchain/sifnode/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func InitGenesis(ctx sdk.Context, keeper Keeper, data types.GenesisState) (res []abci.ValidatorUpdate) {
	log.Println("Oracle module genesis InitGenesis")
	log.Println("Oracle module genesis InitGenesis %v", data.AddressWhitelist[0].String())

	if data.AddressWhitelist != nil {
		log.Println("Oracle nil module genesis InitGenesis %v", data)
		keeper.SetOracleWhiteList(ctx, data.AddressWhitelist)
	}

	return []abci.ValidatorUpdate{}
}

func ExportGenesis(ctx sdk.Context, keeper Keeper) types.GenesisState {
	whiteList := keeper.GetOracleWhiteList(ctx)
	return GenesisState{
		AddressWhitelist: whiteList,
	}
}

// ValidateGenesis validates the oracle genesis parameters
func ValidateGenesis(data GenesisState) error {
	return nil
}
