package keeper

import (
	"fmt"

	tokenregistrytypes "github.com/Sifchain/sifnode/x/tokenregistry/types"
	mintkeeper "github.com/cosmos/cosmos-sdk/x/mint/keeper"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/Sifchain/sifnode/x/clp/types"
)

// Keeper of the clp store
type Keeper struct {
	storeKey            sdk.StoreKey
	cdc                 codec.BinaryCodec
	bankKeeper          types.BankKeeper
	authKeeper          types.AuthKeeper
	tokenRegistryKeeper types.TokenRegistryKeeper
	adminKeeper         types.AdminKeeper
	mintKeeper          mintkeeper.Keeper
	paramstore          paramtypes.Subspace
}

// NewKeeper creates a clp keeper
func NewKeeper(cdc codec.BinaryCodec, key sdk.StoreKey, bankkeeper types.BankKeeper, accountKeeper types.AuthKeeper,
	tokenRegistryKeeper tokenregistrytypes.Keeper, adminKeeper types.AdminKeeper, mintKeeper mintkeeper.Keeper, ps paramtypes.Subspace) Keeper {
	// set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}
	keeper := Keeper{
		storeKey:            key,
		cdc:                 cdc,
		bankKeeper:          bankkeeper,
		authKeeper:          accountKeeper,
		tokenRegistryKeeper: tokenRegistryKeeper,
		adminKeeper:         adminKeeper,
		mintKeeper:          mintKeeper,
		paramstore:          ps,
	}
	return keeper
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) Codec() codec.BinaryCodec {
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

func (k Keeper) GetNormalizationFactor(decimals int64) (sdk.Dec, bool) {
	normalizationFactor := sdk.NewDec(1)
	adjustExternalToken := false
	nf := decimals
	if nf != 18 {
		var diffFactor int64
		if nf < 18 {
			diffFactor = 18 - nf
			adjustExternalToken = true
		} else {
			diffFactor = nf - 18
		}
		normalizationFactor = sdk.NewDec(10).Power(uint64(diffFactor))
	}
	return normalizationFactor, adjustExternalToken
}

func (k Keeper) GetNormalizationFactorFromAsset(ctx sdk.Context, asset types.Asset) (sdk.Dec, bool, error) {
	registry := k.tokenRegistryKeeper.GetRegistry(ctx)
	registryEntry, err := k.tokenRegistryKeeper.GetEntry(registry, asset.Symbol)
	if err != nil {
		return sdk.Dec{}, false, tokenregistrytypes.ErrNotFound
	}

	nf, adjust := k.GetNormalizationFactor(registryEntry.Decimals)

	return nf, adjust, nil
}

func (k Keeper) GetSymmetryThreshold(ctx sdk.Context) sdk.Dec {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.SymmetryThresholdPrefix)
	if bz == nil {
		return sdk.NewDecWithPrec(5, 5)
	}
	var setThreshold types.MsgSetSymmetryThreshold
	k.cdc.MustUnmarshal(bz, &setThreshold)
	return setThreshold.Threshold
}

func (k Keeper) SetSymmetryThreshold(ctx sdk.Context, setThreshold *types.MsgSetSymmetryThreshold) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(setThreshold)
	store.Set(types.SymmetryThresholdPrefix, bz)
}
