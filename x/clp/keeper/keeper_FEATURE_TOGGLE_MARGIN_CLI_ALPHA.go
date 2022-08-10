//go:build FEATURE_TOGGLE_MARGIN_CLI_ALPHA
// +build FEATURE_TOGGLE_MARGIN_CLI_ALPHA

package keeper

import (
	margintypes "github.com/Sifchain/sifnode/x/margin/types"
	tokenregistrytypes "github.com/Sifchain/sifnode/x/tokenregistry/types"
	mintkeeper "github.com/cosmos/cosmos-sdk/x/mint/keeper"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

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
	getMarginKeeper     func() margintypes.Keeper
	paramstore          paramtypes.Subspace
}

// NewKeeper creates a clp keeper
func NewKeeper(cdc codec.BinaryCodec, key sdk.StoreKey, bankkeeper types.BankKeeper, accountKeeper types.AuthKeeper,
	tokenRegistryKeeper tokenregistrytypes.Keeper, adminKeeper types.AdminKeeper, mintKeeper mintkeeper.Keeper, getMarginKeeper func() margintypes.Keeper, ps paramtypes.Subspace) Keeper {
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
		getMarginKeeper:     getMarginKeeper,
		paramstore:          ps,
	}
	return keeper
}

func (k Keeper) GetMarginKeeper() margintypes.Keeper {
	return k.getMarginKeeper()
}
