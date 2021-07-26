package keeper

import (
	"encoding/binary"

	"github.com/Sifchain/sifnode/x/ethbridge/types"
	oracletypes "github.com/Sifchain/sifnode/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetAndUpdateGlobalNonce get current global nonce and update it
func (k Keeper) GetAndUpdateGlobalNonce(ctx sdk.Context, networkDescriptor oracletypes.NetworkDescriptor) uint64 {
	prefix := k.GetGlobalNoncePrefix(ctx, networkDescriptor)
	store := ctx.KVStore(k.storeKey)

	if !k.ExistsGlobalNonce(ctx, prefix) {
		nextGlobalNonce := uint64(1)
		bs := make([]byte, 8)
		binary.LittleEndian.PutUint64(bs, nextGlobalNonce)

		store.Set(prefix, bs)
		return uint64(0)
	}

	value := store.Get(prefix)
	globalNonce := binary.LittleEndian.Uint64(value)

	bs := make([]byte, 8)
	binary.LittleEndian.PutUint64(bs, globalNonce+1)

	return globalNonce
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
