package keeper

import (
	"errors"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/Sifchain/sifnode/x/oracle/types"
)

// GetProphecies returns all prophecies
func (k Keeper) GetProphecies(ctx sdk.Context) []types.Prophecy {
	var prophecies []types.Prophecy
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.ProphecyPrefix)
	for ; iter.Valid(); iter.Next() {
		var prophecy types.Prophecy
		k.cdc.MustUnmarshalBinaryBare(iter.Value(), &prophecy)
		prophecies = append(prophecies, prophecy)
	}
	return prophecies
}

// GetProphecy gets the entire prophecy data struct for a given id
func (k Keeper) GetProphecy(ctx sdk.Context, prophecyID []byte) (types.Prophecy, bool) {

	store := ctx.KVStore(k.storeKey)
	bz := store.Get(append(types.ProphecyPrefix, prophecyID[:]...))

	if bz == nil {
		return types.Prophecy{}, false
	}

	var prophecy types.Prophecy
	k.cdc.MustUnmarshalBinaryBare(bz, &prophecy)

	return prophecy, true
}

// SetProphecy saves a prophecy with an initial claim
func (k Keeper) SetProphecy(ctx sdk.Context, prophecy types.Prophecy) {
	store := ctx.KVStore(k.storeKey)

	storePrefix := append(types.ProphecyPrefix, prophecy.Id[:]...)

	store.Set(storePrefix, k.cdc.MustMarshalBinaryBare(&prophecy))
}

// GetSignature return a prophecy's signatures
func (k Keeper) GetSignature(ctx sdk.Context, prophecyID []byte) (types.ProphecySignatures, bool) {
	var prophecySignatures types.ProphecySignatures
	store := ctx.KVStore(k.storeKey)

	// storePrefix := append(types.SignaturePrefix, prophecyID[:]...)
	bz := store.Get(append(types.SignaturePrefix, prophecyID[:]...))
	if bz == nil {
		return types.ProphecySignatures{}, false
	}

	k.cdc.MustUnmarshalBinaryBare(bz, &prophecySignatures)

	return prophecySignatures, true
}

// SetSignature saves a prophecy with an initial value
func (k Keeper) SetSignature(ctx sdk.Context, prophecyID []byte, networkDescriptor types.NetworkDescriptor) {

	store := ctx.KVStore(k.storeKey)

	storePrefix := append(types.SignaturePrefix, prophecyID[:]...)

	prophecySignatures := types.ProphecySignatures{
		NetworkDescriptor: networkDescriptor,
		EthereumAddress:   []string{},
		Signatures:        []string{},
	}

	store.Set(storePrefix, k.cdc.MustMarshalBinaryBare(&prophecySignatures))
}

// AppendSignature add a new ethereum address and signature to prophecy
func (k Keeper) AppendSignature(ctx sdk.Context, prophecyID []byte, ethereumAddress, signature string) error {
	store := ctx.KVStore(k.storeKey)

	prophecySignatures, ok := k.GetSignature(ctx, prophecyID)
	if !ok {
		return errors.New("can not get the prophecy")
	}

	prophecySignatures.EthereumAddress = append(prophecySignatures.EthereumAddress, ethereumAddress)
	prophecySignatures.Signatures = append(prophecySignatures.Signatures, signature)

	storePrefix := append(types.SignaturePrefix, prophecyID[:]...)

	store.Set(storePrefix, k.cdc.MustMarshalBinaryBare(&prophecySignatures))
	return nil
}
