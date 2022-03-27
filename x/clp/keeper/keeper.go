package keeper

import (
	"fmt"

	tokenregistrytypes "github.com/Sifchain/sifnode/x/tokenregistry/types"

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
func (keeper Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (keeper Keeper) Codec() codec.BinaryCodec {
	return keeper.cdc
}

func (keeper Keeper) GetBankKeeper() types.BankKeeper {
	return keeper.bankKeeper
}

func (keeper Keeper) GetAuthKeeper() types.AuthKeeper {
	return keeper.authKeeper
}

func (keeper Keeper) Exists(ctx sdk.Context, key []byte) bool {
	store := ctx.KVStore(keeper.storeKey)
	return store.Has(key)
}

func (keeper Keeper) SendCoins(ctx sdk.Context, from sdk.AccAddress, to sdk.AccAddress, coins sdk.Coins) error {
	return keeper.bankKeeper.SendCoins(ctx, from, to, coins)
}

func (keeper Keeper) HasBalance(ctx sdk.Context, addr sdk.AccAddress, coin sdk.Coin) bool {
	return keeper.bankKeeper.HasBalance(ctx, addr, coin)
}

func (keeper Keeper) GetNormalizationFactor(decimals int64) (sdk.Dec, bool) {
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
