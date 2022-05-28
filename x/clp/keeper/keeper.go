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
	mintKeeper          mintkeeper.Keeper
	paramstore          paramtypes.Subspace
}

// NewKeeper creates a clp keeper
func NewKeeper(cdc codec.BinaryCodec, key sdk.StoreKey, bankkeeper types.BankKeeper, accountKeeper types.AuthKeeper, tokenRegistryKeeper tokenregistrytypes.Keeper, mintKeeper mintkeeper.Keeper, ps paramtypes.Subspace) Keeper {
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

func (k Keeper) GetNormalizationFactorFromAsset(ctx sdk.Context, asset types.Asset) (sdk.Dec, bool) {
	registry := k.tokenRegistryKeeper.GetRegistry(ctx)
	registryEntry, err := k.tokenRegistryKeeper.GetEntry(registry, asset.Symbol)
	if err != nil {
		return sdk.Dec{}, false
	}
	return k.GetNormalizationFactor(registryEntry.Decimals)
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

func (k Keeper) GetSymmetryRatio(ctx sdk.Context) sdk.Dec {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.SymmetryThresholdPrefix)
	if bz == nil {
		return sdk.NewDecWithPrec(5, 4)
	}
	var setThreshold types.MsgSetSymmetryThreshold
	k.cdc.MustUnmarshal(bz, &setThreshold)
	return setThreshold.Ratio
}

func (k Keeper) SetSymmetryThreshold(ctx sdk.Context, setThreshold *types.MsgSetSymmetryThreshold) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(setThreshold)
	store.Set(types.SymmetryThresholdPrefix, bz)
}

func (k Keeper) GetSwapPermissionIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, types.SwapPermissionStorePrefix)
}

func (k Keeper) GetSwapPermissions(ctx sdk.Context) []*types.SwapPermission {
	var swapPermissions []*types.SwapPermission
	iterator := k.GetSwapPermissionIterator(ctx)
	defer func(iterator sdk.Iterator) {
		err := iterator.Close()
		if err != nil {
			panic(err)
		}
	}(iterator)
	for ; iterator.Valid(); iterator.Next() {
		var st types.SwapPermission
		bytesValue := iterator.Value()
		k.cdc.MustUnmarshal(bytesValue, &st)
		swapPermissions = append(swapPermissions, &st)
	}
	return swapPermissions
}

func (k Keeper) GetSwapTypes(ctx sdk.Context) []types.SwapType {
	var swapTypes []types.SwapType
	iterator := k.GetSwapPermissionIterator(ctx)
	defer func(iterator sdk.Iterator) {
		err := iterator.Close()
		if err != nil {
			panic(err)
		}
	}(iterator)
	for ; iterator.Valid(); iterator.Next() {
		var st types.SwapPermission
		bytesValue := iterator.Value()
		k.cdc.MustUnmarshal(bytesValue, &st)
		swapTypes = append(swapTypes, st.SwapType)
	}
	return swapTypes
}

func (k Keeper) AddSwapPermission(ctx sdk.Context, swapPermission *types.SwapPermission) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetSwapPermissionKey(swapPermission.SwapType)
	store.Set(key, k.cdc.MustMarshal(swapPermission))
}

func (k Keeper) RemoveSwapPermission(ctx sdk.Context, swapPermission *types.SwapPermission) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetSwapPermissionKey(swapPermission.SwapType)
	store.Delete(key)
}

func (k Keeper) checkSwapPermission(ctx sdk.Context, swapType types.SwapType) bool {
	store := ctx.KVStore(k.storeKey)
	key := types.GetSwapPermissionKey(swapType)
	return store.Has(key)
}
