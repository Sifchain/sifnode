package keeper

import (
	"github.com/Sifchain/sifnode/x/tokenregistry/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k keeper) getPeggy1DenomMappingPrefix(ctx sdk.Context, peggy1_denom string) []byte {
	return append(types.Peggy1DenomMappingPrefix, []byte(peggy1_denom)...)
}

func (k keeper) GetPeggy2Denom(ctx sdk.Context, peggy1_denom string) (string, bool) {
	denom_prefix := k.getPeggy1DenomMappingPrefix(ctx, peggy1_denom)

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.Peggy1DenomMappingPrefix)

	b := store.Get(denom_prefix)

	if b == nil {
		return "", false
	}
	return string(b), true
}

func (k keeper) SetPeggy2Denom(ctx sdk.Context, peggy1_denom string, peggy2_denom string) {
	denom_prefix := k.getPeggy1DenomMappingPrefix(ctx, peggy1_denom)

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.Peggy1DenomMappingPrefix)

	store.Set(denom_prefix, []byte(peggy2_denom))
}
