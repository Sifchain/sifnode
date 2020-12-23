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

func (k Keeper) GetBankKeeper() types.BankKeeper {
	return k.bankKeeper
}

func (k Keeper) GetSupplyKeeper() types.SupplyKeeper {
	return k.supplyKeeper
}

// TODO add functionality to keep track of how much a user withdrew , to prevent spam .
//func (k Keeper) HasCoins(ctx sdk.Context, user sdk.AccAddress, coins sdk.Coins) bool {
//	return k.bankKeeper.HasCoins(ctx, user, coins)
//}
//
//func (k Keeper) SendCoins(ctx sdk.Context, from sdk.AccAddress, to sdk.AccAddress, coins sdk.Coins) error {
//	return k.bankKeeper.SendCoins(ctx, from, to, coins)
//}
