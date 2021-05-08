package keeper

import (
	"fmt"
	"github.com/Sifchain/sifnode/x/clp/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/libs/log"
)

// Keeper of the clp store
type Keeper struct {
	storeKey     sdk.StoreKey
	cdc          *codec.Codec
	bankKeeper   types.BankKeeper
	supplyKeeper types.SupplyKeeper
}

// NewKeeper creates a dispensation keeper
func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, bankkeeper types.BankKeeper, supplyKeeper types.SupplyKeeper) Keeper {
	keeper := Keeper{
		storeKey:     key,
		cdc:          cdc,
		bankKeeper:   bankkeeper,
		supplyKeeper: supplyKeeper,
	}
	return keeper
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) Codec() *codec.Codec {
	return k.cdc
}

func (k Keeper) GetBankKeeper() types.BankKeeper {
	return k.bankKeeper
}

func (k Keeper) GetSupplyKeeper() types.SupplyKeeper {
	return k.supplyKeeper
}

func (k Keeper) Exists(ctx sdk.Context, key []byte) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(key)
}
