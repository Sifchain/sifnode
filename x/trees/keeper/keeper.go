package keeper

import (
	"fmt"

	"github.com/Sifchain/sifnode/x/trees/types"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Keeper struct {
	storeKey     sdk.StoreKey
	cdc          *codec.Codec
	bankKeeper   types.BankKeeper
}

func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, bankKeeper types.BankKeeper) Keeper {
	keeper := Keeper{
		bankKeeper:   bankKeeper,
		storeKey:     key,
		cdc:          cdc,
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
