package keeper

import (
	"fmt"

	"github.com/Sifchain/sifnode/x/faucet/types"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/supply"
)

// Keeper of the faucet store
type Keeper struct {
	supplyKeeper supply.Keeper
	storeKey     sdk.StoreKey
	cdc          *codec.Codec
}

// NewKeeper creates a faucet keeper
func NewKeeper(supplyKeeper supply.Keeper, cdc *codec.Codec, key sdk.StoreKey) Keeper {
	keeper := Keeper{
		supplyKeeper: supplyKeeper,
		storeKey:     key,
		cdc:          cdc,
		// paramspace: paramspace.WithKeyTable(types.ParamKeyTable()),
	}
	return keeper
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}
