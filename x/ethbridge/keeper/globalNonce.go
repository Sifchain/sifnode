package keeper

import (
	"encoding/binary"

	"github.com/Sifchain/sifnode/x/ethbridge/types"
	oracletypes "github.com/Sifchain/sifnode/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetGlobalNonce get current global nonce and update it
func (k Keeper) GetGlobalNonce(ctx sdk.Context, networkDescriptor oracletypes.NetworkDescriptor) uint64 {
	prefix := k.GetGlobalNoncePrefix(ctx, networkDescriptor)
	store := ctx.KVStore(k.storeKey)

	if !k.ExistsGlobalNonce(ctx, prefix) {
		// global nonce start from 1
		return uint64(1)
	}

	value := store.Get(prefix)
	globalNonce := binary.LittleEndian.Uint64(value)

	return globalNonce
}

// UpdateGlobalNonce get current global nonce and update it
func (k Keeper) UpdateGlobalNonce(ctx sdk.Context, networkDescriptor oracletypes.NetworkDescriptor) {
	prefix := k.GetGlobalNoncePrefix(ctx, networkDescriptor)
	store := ctx.KVStore(k.storeKey)

	globalNonce := k.GetGlobalNonce(ctx, networkDescriptor)

	bs := make([]byte, 8)
	binary.LittleEndian.PutUint64(bs, globalNonce+1)
	store.Set(prefix, bs)
}

// ExistsGlobalNonce check if the global nonce exists
func (k Keeper) ExistsGlobalNonce(ctx sdk.Context, prefix []byte) bool {
	if !k.Exists(ctx, prefix) {
		return false
	}
	return true
}

// GetGlobalNoncePrefix compute the prefix
func (k Keeper) GetGlobalNoncePrefix(ctx sdk.Context, networkDescriptor oracletypes.NetworkDescriptor) []byte {
	bs := make([]byte, 4)
	binary.LittleEndian.PutUint32(bs, uint32(networkDescriptor))

	return append(types.GlobalNoncePrefix, bs[:]...)
}
