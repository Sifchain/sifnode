package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/Sifchain/sifnode/x/clp/types"
)

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

func (k Keeper) GetAssetDecimals(ctx sdk.Context, asset types.Asset) (uint8, error) {
	registry := k.tokenRegistryKeeper.GetRegistry(ctx)
	registryEntry, err := k.tokenRegistryKeeper.GetEntry(registry, asset.Symbol)
	if err != nil {
		return 0, err
	}
	return Int64ToUint8Safe(registryEntry.Decimals)
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
