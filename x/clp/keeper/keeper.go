package keeper

import (
	"fmt"
	"github.com/Sifchain/sifnode/x/clp/types"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Keeper of the clp store
type Keeper struct {
	storeKey   sdk.StoreKey
	cdc        *codec.Codec
	bankKeeper types.BankKeeper
}

// NewKeeper creates a clp keeper
func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, bankkeeper types.BankKeeper, paramspace types.ParamSubspace) Keeper {
	keeper := Keeper{
		storeKey:   key,
		cdc:        cdc,
		bankKeeper: bankkeeper,
		//paramspace: paramspace.WithKeyTable(types.ParamKeyTable()),
	}
	return keeper
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) SetAsset(ctx sdk.Context, assetName string, asset types.Asset) {

}

func (k Keeper) GetAsset(tx sdk.Context, assetName string) types.Asset {
	return types.Asset{}
}
