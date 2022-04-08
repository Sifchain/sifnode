package keeper

import (
	"bytes"
	"fmt"

	tokenregistrytypes "github.com/Sifchain/sifnode/x/tokenregistry/types"
	gogotypes "github.com/gogo/protobuf/types"

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
	paramstore          paramtypes.Subspace
}

// NewKeeper creates a clp keeper
func NewKeeper(cdc codec.BinaryCodec, key sdk.StoreKey, bankkeeper types.BankKeeper, accountKeeper types.AuthKeeper, tokenRegistryKeeper tokenregistrytypes.Keeper, ps paramtypes.Subspace) Keeper {
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

func (k Keeper) SetAdminAccount(ctx sdk.Context, adminAccount sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	key := types.AdminAccountStorePrefix
	store.Set(key, k.cdc.MustMarshal(&gogotypes.BytesValue{Value: adminAccount}))
}

func (k Keeper) IsAdminAccount(ctx sdk.Context, adminAccount sdk.AccAddress) bool {
	account := k.GetAdminAccount(ctx)
	if account == nil {
		return false
	}
	return bytes.Equal(account, adminAccount)
}

func (k Keeper) GetAdminAccount(ctx sdk.Context) (adminAccount sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	key := types.AdminAccountStorePrefix
	bz := store.Get(key)
	acc := gogotypes.BytesValue{}
	k.cdc.MustUnmarshal(bz, &acc)
	adminAccount = sdk.AccAddress(acc.Value)
	return adminAccount
}
