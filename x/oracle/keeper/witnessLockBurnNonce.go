package keeper

import (
	"bytes"
	"encoding/binary"

	"github.com/Sifchain/sifnode/x/instrumentation"
	"github.com/Sifchain/sifnode/x/oracle/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SetWitnessLockBurnNonce set the Witness lock burn nonce for each relayer
func (k Keeper) SetWitnessLockBurnNonce(ctx sdk.Context, networkDescriptor types.NetworkDescriptor, valAccount sdk.ValAddress, newNonce uint64) {
	store := ctx.KVStore(k.storeKey)
	key := k.GetWitnessLockBurnSequencePrefix(networkDescriptor, valAccount)

	bs := make([]byte, 8)
	binary.BigEndian.PutUint64(bs, newNonce)

	instrumentation.PeggyCheckpoint(ctx.Logger(), instrumentation.SetWitnessLockBurnNonce, "networkDescriptor", networkDescriptor, "valAccount", valAccount, "newNonce", newNonce, "key", key)

	store.Set(key, bs)
}

// GetWitnessLockBurnSequence return Witness lock burn nonce
func (k Keeper) GetWitnessLockBurnSequence(ctx sdk.Context, networkDescriptor types.NetworkDescriptor, valAccount sdk.ValAddress) uint64 {
	store := ctx.KVStore(k.storeKey)
	key := k.GetWitnessLockBurnSequencePrefix(networkDescriptor, valAccount)

	// nonce start from 1, 0 represent the relayer is a new one
	if !store.Has(key) {
		return 0
	}

	bz := store.Get(key)
	return binary.BigEndian.Uint64(bz)
}

// GetWitnessLockBurnSequencePrefix return storage prefix
func (k Keeper) GetWitnessLockBurnSequencePrefix(networkDescriptor types.NetworkDescriptor, valAccount sdk.ValAddress) []byte {
	bytebuf := bytes.NewBuffer([]byte{})
	err := binary.Write(bytebuf, binary.BigEndian, networkDescriptor)
	if err != nil {
		panic(err.Error())
	}
	tmpKey := append(types.WitnessLockBurnNoncePrefix, bytebuf.Bytes()...)
	return append(tmpKey, valAccount...)
}

// GetAllWitnessLockBurnSequence get all witnessLockBurnSequence needed for all validators
func (k Keeper) GetAllWitnessLockBurnSequence(ctx sdk.Context) map[string]uint64 {
	sequences := make(map[string]uint64)
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, types.WitnessLockBurnNoncePrefix)
	defer func(iterator sdk.Iterator) {
		err := iterator.Close()
		if err != nil {
			panic(err)
		}
	}(iterator)
	for ; iterator.Valid(); iterator.Next() {
		key := iterator.Key()
		value := iterator.Value()
		sequences[string(key)] = binary.BigEndian.Uint64(value)
	}
	return sequences
}

// DecodeWitnessLockBurnNonceKey extract the NetworkDescriptor and ValAddress from raw key
func DecodeWitnessLockBurnNonceKey(key []byte) (types.NetworkDescriptor, sdk.ValAddress) {
	prefixLen := len(types.WitnessLockBurnNoncePrefix)

	// the length must larger than 5
	if len(key) < prefixLen+4 {
		panic("key for WitnessLockBurnSequence with wrong length")
	}

	networkDescriptorKey := key[prefixLen : prefixLen+4]
	addressKey := key[prefixLen+4:]

	networkDescriptor := binary.BigEndian.Uint32(networkDescriptorKey)
	address := sdk.ValAddress(addressKey)

	return types.NetworkDescriptor(networkDescriptor), address
}

func (k Keeper) SetWitnessLockBurnNonceViaRawKey(ctx sdk.Context, key []byte, nonce uint64) {
	networkDescriptor, valAddress := DecodeWitnessLockBurnNonceKey(key)

	k.SetWitnessLockBurnNonce(ctx, networkDescriptor, valAddress, nonce)
}
