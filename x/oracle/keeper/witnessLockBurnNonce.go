package keeper

import (
	"bytes"
	"encoding/binary"

	"github.com/Sifchain/sifnode/x/oracle/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SetWitnessLockBurnNonce set the Witness lock burn nonce for each relayer
func (k Keeper) SetWitnessLockBurnNonce(ctx sdk.Context, networkDescriptor types.NetworkDescriptor, valAccount sdk.ValAddress, newNonce uint64) {
	store := ctx.KVStore(k.storeKey)
	key := k.getWitnessLockBurnNoncePrefix(networkDescriptor, valAccount)

	bs := make([]byte, 8)
	binary.BigEndian.PutUint64(bs, newNonce)

	store.Set(key, bs)
}

// GetWitnessLockBurnNonce return Witness lock burn nonce
func (k Keeper) GetWitnessLockBurnNonce(ctx sdk.Context, networkDescriptor types.NetworkDescriptor, valAccount sdk.ValAddress) uint64 {
	store := ctx.KVStore(k.storeKey)
	key := k.getWitnessLockBurnNoncePrefix(networkDescriptor, valAccount)

	// nonce start from 1, 0 represent the relayer is a new one
	if !store.Has(key) {
		return 0
	}

	bz := store.Get(key)
	return binary.BigEndian.Uint64(bz)
}

// getWitnessLockBurnNoncePrefix return storage prefix
func (k Keeper) getWitnessLockBurnNoncePrefix(networkDescriptor types.NetworkDescriptor, valAccount sdk.ValAddress) []byte {
	bytebuf := bytes.NewBuffer([]byte{})
	_ = binary.Write(bytebuf, binary.BigEndian, networkDescriptor)
	tmpKey := append(types.WitnessLockBurnNoncePrefix, bytebuf.Bytes()...)
	return append(tmpKey, valAccount...)
}
