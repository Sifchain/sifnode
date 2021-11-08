package keeper

import (
	"bytes"
	"encoding/binary"

	"github.com/Sifchain/sifnode/x/ethbridge/types"
	oracletypes "github.com/Sifchain/sifnode/x/oracle/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SetEthereumLockBurnSequence set the ethereum lock burn nonce for each relayer
func (k Keeper) SetEthereumLockBurnSequence(ctx sdk.Context, networkDescriptor oracletypes.NetworkDescriptor, valAccount sdk.ValAddress, newNonce uint64) {
	store := ctx.KVStore(k.storeKey)
	key := k.getEthereumLockBurnSequencePrefix(networkDescriptor, valAccount)

	bs := make([]byte, 8)
	binary.BigEndian.PutUint64(bs, newNonce)

	store.Set(key, bs)
}

// GetEthereumLockBurnSequence return ethereum lock burn nonce
func (k Keeper) GetEthereumLockBurnSequence(ctx sdk.Context, networkDescriptor oracletypes.NetworkDescriptor, valAccount sdk.ValAddress) uint64 {
	store := ctx.KVStore(k.storeKey)
	key := k.getEthereumLockBurnSequencePrefix(networkDescriptor, valAccount)

	// nonces start from 0, and the first ethereum transaction
	// should have a nonce of 1
	if !store.Has(key) {
		return 0
	}

	bz := store.Get(key)
	return binary.BigEndian.Uint64(bz)
}

// getEthereumLockBurnSequencePrefix return storage prefix
func (k Keeper) getEthereumLockBurnSequencePrefix(networkDescriptor oracletypes.NetworkDescriptor, valAccount sdk.ValAddress) []byte {
	bytebuf := bytes.NewBuffer([]byte{})
	_ = binary.Write(bytebuf, binary.BigEndian, networkDescriptor)
	tmpKey := append(types.EthereumLockBurnSequencePrefix, bytebuf.Bytes()...)
	return append(tmpKey, valAccount...)
}
