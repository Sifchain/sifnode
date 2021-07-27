package keeper

import (
	"fmt"

	whitelisttypes "github.com/Sifchain/sifnode/x/tokenregistry/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/Sifchain/sifnode/x/clp/types"
)

// Keeper of the clp store
type Keeper struct {
	storeKey        sdk.StoreKey
	cdc             codec.BinaryMarshaler
	bankKeeper      types.BankKeeper
	authKeeper      types.AuthKeeper
	whitelistKeeper types.WhitelistKeeper
	paramstore      paramtypes.Subspace
}

// NewKeeper creates a clp keeper
func NewKeeper(cdc codec.BinaryMarshaler, key sdk.StoreKey, bankkeeper types.BankKeeper, accountKeeper types.AuthKeeper, whitelistKeeper whitelisttypes.Keeper, ps paramtypes.Subspace) Keeper {
	// set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	keeper := Keeper{
		storeKey:        key,
		cdc:             cdc,
		bankKeeper:      bankkeeper,
		authKeeper:      accountKeeper,
		whitelistKeeper: whitelistKeeper,
		paramstore:      ps,
	}
	return keeper
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) Codec() codec.BinaryMarshaler {
	return k.cdc
}

func (k Keeper) GetBankKeeper() types.BankKeeper {
	return k.bankKeeper
}

func (k Keeper) GetAuthKeeper() types.AuthKeeper {
	return k.authKeeper
}

func (k Keeper) Exists(ctx sdk.Context, key []byte) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(key)
}

func (k Keeper) SendCoins(ctx sdk.Context, from sdk.AccAddress, to sdk.AccAddress, coins sdk.Coins) error {
	return k.bankKeeper.SendCoins(ctx, from, to, coins)
}

func (k Keeper) HasBalance(ctx sdk.Context, addr sdk.AccAddress, coin sdk.Coin) bool {
	return k.bankKeeper.HasBalance(ctx, addr, coin)
}

func (k Keeper) GetNormalizationFactor(ctx sdk.Context, denom string) (sdk.Dec, bool) {
	normalizationFactor := sdk.NewDec(1)
	adjustExternalToken := false
	entry := k.whitelistKeeper.GetDenom(ctx, denom)
	if !entry.IsWhitelisted {
		return normalizationFactor, adjustExternalToken
	}
	nf := entry.Decimals
	if nf != 18 {
		adjustExternalToken = true
		diffFactor := 18 - nf
		if diffFactor < 0 {
			diffFactor = nf - 18
			adjustExternalToken = false
		}
		normalizationFactor = sdk.NewDec(10).Power(uint64(diffFactor))
	}
	return normalizationFactor, adjustExternalToken
}
