package keeper

import (
	"encoding/binary"
	"fmt"

	"github.com/Sifchain/sifnode/x/ethbridge/types"
	oracletypes "github.com/Sifchain/sifnode/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetAndUpdateBlobalNonce get current global nonce and update it
func (k Keeper) GetAndUpdateBlobalNonce(ctx sdk.Context, networkDescriptor oracletypes.NetworkDescriptor) uint64 {
	prefix := k.GetBlobalNoncePrefix(ctx, networkDescriptor)
	store := ctx.KVStore(k.storeKey)

	if !k.ExistsBlobalNonce(ctx, prefix) {
		nextGlobalNonce := uint64(1)
		bs := make([]byte, 8)
		binary.LittleEndian.PutUint64(bs, nextGlobalNonce)

		store.Set(prefix, bs)
		return uint64(0)
	}

	value := store.Get(prefix)
	fmt.Printf("value is %v\n", value)
	globalNonce := binary.LittleEndian.Uint64(value)

	bs := make([]byte, 8)
	binary.LittleEndian.PutUint64(bs, globalNonce+1)

	return globalNonce
}

// ExistsBlobalNonce get peggy token list
func (k Keeper) ExistsBlobalNonce(ctx sdk.Context, prefix []byte) bool {
	if !k.Exists(ctx, prefix) {
		return false
	}
	return true
}

// GetBlobalNoncePrefix compute the prefix
func (k Keeper) GetBlobalNoncePrefix(ctx sdk.Context, networkDescriptor oracletypes.NetworkDescriptor) []byte {
	bs := make([]byte, 4)
	binary.LittleEndian.PutUint32(bs, uint32(networkDescriptor))

	return append(types.GlobalNoncePrefix, bs[:]...)
}
