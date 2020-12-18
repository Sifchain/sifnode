package keeper

import (
	"fmt"

	"github.com/Sifchain/sifnode/x/faucet/types"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Keeper of the faucet store
type Keeper struct {
	storeKey     sdk.StoreKey
	cdc          *codec.Codec
	supplyKeeper types.SupplyKeeper
	bankKeeper   types.BankKeeper
}

// NewKeeper creates a faucet keeper
func NewKeeper(supplyKeeper types.SupplyKeeper, cdc *codec.Codec, key sdk.StoreKey, bankKeeper types.BankKeeper) Keeper {
	keeper := Keeper{
		supplyKeeper: supplyKeeper,
		bankKeeper:   bankKeeper,
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

// Get returns the pubkey from the adddress-pubkey relation
func (k Keeper) Get(ctx sdk.Context, key string) (interface{} /* TODO: Fill out this type */, error) {
	store := ctx.KVStore(k.storeKey)
	var item interface{} /* TODO: Fill out this type */
	byteKey := []byte(key)
	err := k.cdc.UnmarshalBinaryLengthPrefixed(store.Get(byteKey), &item)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (k Keeper) set(ctx sdk.Context, key string, value interface{} /* TODO: fill out this type */) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(value)
	store.Set([]byte(key), bz)
}

func (k Keeper) delete(ctx sdk.Context, key string) {
	store := ctx.KVStore(k.storeKey)
	store.Delete([]byte(key))
}

func (k Keeper) GetBankKeeper() types.BankKeeper {
	return k.bankKeeper
}

func (k Keeper) GetSupplyKeeper() types.SupplyKeeper {
	return k.supplyKeeper
}

func (k Keeper) HasCoins(ctx sdk.Context, user sdk.AccAddress, coins sdk.Coins) bool {
	return k.bankKeeper.HasCoins(ctx, user, coins)
}

func (k Keeper) SendCoins(ctx sdk.Context, from sdk.AccAddress, to sdk.AccAddress, coins sdk.Coins) error {
	return k.bankKeeper.SendCoins(ctx, from, to, coins)
}
