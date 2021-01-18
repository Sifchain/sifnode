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
	storeKey   sdk.StoreKey
	cdc        *codec.BinaryMarshaler
	bankKeeper types.BankKeeper
	authKeeper types.AccountKeeper
	paramSpace types.ParamSubspace
}

// NewKeeper creates a clp keeper
func NewKeeper(cdc *codec.BinaryMarshaler, key sdk.StoreKey, bankkeeper types.BankKeeper, supplyKeeper types.AccountKeeper, paramsSubspace types.ParamSubspace) Keeper {
	keeper := Keeper{
		storeKey:   key,
		cdc:        cdc,
		bankKeeper: bankkeeper,
		authKeeper: supplyKeeper,
		paramSpace: paramsSubspace,
	}
	return keeper
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) Codec() *codec.BinaryMarshaler {
	return k.cdc
}

func (k Keeper) GetBankKeeper() types.BankKeeper {
	return k.bankKeeper
}

func (k Keeper) GetAccountKeeper() types.AccountKeeper {
	return k.authKeeper
}

func (k Keeper) Exists(ctx sdk.Context, key []byte) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(key)
}

func (k Keeper) SendCoins(ctx sdk.Context, from sdk.AccAddress, to sdk.AccAddress, coins sdk.Coins) error {
	return k.bankKeeper.SendCoins(ctx, from, to, coins)
}

func (k Keeper) HasCoins(ctx sdk.Context, user sdk.AccAddress, coins sdk.Coins) bool {
	return k.bankKeeper.HasCoins(ctx, user, coins)
}
