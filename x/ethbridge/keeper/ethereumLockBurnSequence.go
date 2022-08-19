package keeper

import (
	"bytes"
	"encoding/binary"

	"github.com/Sifchain/sifnode/x/ethbridge/types"
	oracletypes "github.com/Sifchain/sifnode/x/oracle/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SetEthereumLockBurnSequence set the ethereum lock burn nonce for each relayer
func (k Keeper) SetEthereumLockBurnSequence(ctx sdk.Context, networkDescriptor oracletypes.NetworkDescriptor, valAccount sdk.ValAddress, newSequence uint64) {
	store := ctx.KVStore(k.storeKey)
	key := k.GetEthereumLockBurnSequencePrefix(networkDescriptor, valAccount)

	bs := make([]byte, 8)
	binary.BigEndian.PutUint64(bs, newSequence)

	store.Set(key, bs)
}

// GetEthereumLockBurnSequence return ethereum lock burn nonce
func (k Keeper) GetEthereumLockBurnSequence(ctx sdk.Context, networkDescriptor oracletypes.NetworkDescriptor, valAccount sdk.ValAddress) uint64 {
	store := ctx.KVStore(k.storeKey)
	key := k.GetEthereumLockBurnSequencePrefix(networkDescriptor, valAccount)

	// nonces start from 0, and the first ethereum transaction
	// should have a nonce of 1
	if !store.Has(key) {
		return 0
	}

	bz := store.Get(key)
	return binary.BigEndian.Uint64(bz)
}

// GetEthereumLockBurnSequencePrefix return storage prefix
func (k Keeper) GetEthereumLockBurnSequencePrefix(networkDescriptor oracletypes.NetworkDescriptor, valAccount sdk.ValAddress) []byte {
	bytebuf := bytes.NewBuffer([]byte{})
	_ = binary.Write(bytebuf, binary.BigEndian, networkDescriptor)
	tmpKey := append(types.EthereumLockBurnSequencePrefix, bytebuf.Bytes()...)
	return append(tmpKey, valAccount...)
}

func (k Keeper) getEthereumLockBurnSequenceIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, types.EthereumLockBurnSequencePrefix)
}

// GetEthereumLockBurnSequences get all sequences from keeper
func (k Keeper) GetEthereumLockBurnSequences(ctx sdk.Context) map[string]uint64 {
	sequences := make(map[string]uint64)
	iterator := k.getEthereumLockBurnSequenceIterator(ctx)
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

// SetSequenceViaRawKey used in import sequence from genesis
func (k Keeper) SetSequenceViaRawKey(ctx sdk.Context, key []byte, newSequence uint64) {
	network, address := DecodeKey(key)
	k.SetEthereumLockBurnSequence(ctx, network, address, newSequence)
}

// DecodeKey extract the NetworkDescriptor and ValAddress from raw key
func DecodeKey(key []byte) (oracletypes.NetworkDescriptor, sdk.ValAddress) {
	prefixLen := len(types.EthereumLockBurnSequencePrefix)

	// the length must larger than 5
	if len(key) < prefixLen+4 {
		panic("key for EthereumLockBurnSequence with wrong length")
	}

	networkDescriptorKey := key[prefixLen : prefixLen+4]
	addressKey := key[prefixLen+4:]

	networkDescriptor := binary.BigEndian.Uint32(networkDescriptorKey)
	address := sdk.ValAddress(addressKey)

	return oracletypes.NetworkDescriptor(networkDescriptor), address
}
