package keeper

import (
	"bytes"
	"encoding/binary"

	"github.com/Sifchain/sifnode/x/ethbridge/types"
	oracletypes "github.com/Sifchain/sifnode/x/oracle/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SetEthereumLockBurnNonce set the ethereum lock burn nonce after prophecy completed in Sifchain
func (k Keeper) SetEthereumLockBurnNonce(ctx sdk.Context, networkDescriptor oracletypes.NetworkDescriptor, valAccount sdk.ValAddress, newNonce uint64) {
	store := ctx.KVStore(k.storeKey)
	key := k.getEthereumLockBurnNoncePrefix(networkDescriptor, valAccount)

	bs := make([]byte, 8)
	binary.BigEndian.PutUint64(bs, newNonce)

	store.Set(key, bs)
}

// GetEthereumLockBurnNonce return ethereum lock burn nonce
func (k Keeper) GetEthereumLockBurnNonce(ctx sdk.Context, networkDescriptor oracletypes.NetworkDescriptor, valAccount sdk.ValAddress) uint64 {
	store := ctx.KVStore(k.storeKey)
	key := k.getEthereumLockBurnNoncePrefix(networkDescriptor, valAccount)

	// nonce start from 1, 0 represent the relayer is a new one
	if !store.Has(key) {
		return 0
	}

	bz := store.Get(key)
	return binary.BigEndian.Uint64(bz)
}

// GetCrossChainFee return crosschain fee
// GetCrossChainFeePrefix return storage prefix
func (k Keeper) getEthereumLockBurnNoncePrefix(networkDescriptor oracletypes.NetworkDescriptor, valAccount sdk.ValAddress) []byte {
	bytebuf := bytes.NewBuffer([]byte{})
	_ = binary.Write(bytebuf, binary.BigEndian, networkDescriptor)
	tmpKey := append(types.EthereumLockBurnNoncePrefix, bytebuf.Bytes()...)
	return append(tmpKey, valAccount...)
}
