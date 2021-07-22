package keeper

import (
	"errors"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/Sifchain/sifnode/x/oracle/types"
)

// ProphecyLiftTime is used to clean outdated prophecy info from keeper
const ProphecyLiftTime = 1000

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

// GetProphecyInfo return a prophecy's signatures
func (k Keeper) GetProphecyInfo(ctx sdk.Context, prophecyID []byte) (types.ProphecyInfo, bool) {
	var prophecySignatures types.ProphecyInfo
	store := ctx.KVStore(k.storeKey)

	// storePrefix := append(types.SignaturePrefix, prophecyID[:]...)
	bz := store.Get(append(types.SignaturePrefix, prophecyID[:]...))
	if bz == nil {
		return types.ProphecyInfo{}, false
	}

	k.cdc.MustUnmarshalBinaryBare(bz, &prophecySignatures)

	return prophecySignatures, true
}

// SetProphecyInfo saves a prophecy with an initial value
func (k Keeper) SetProphecyInfo(ctx sdk.Context, prophecyID []byte, networkDescriptor types.NetworkDescriptor,
	cosmosSender string,
	cosmosSenderSequence uint64,
	ethereumReceiver string,
	tokenSymbol string,
	tokenAmount sdk.Int,
	crosschainFee sdk.Int,
	doublePeg bool,
	globalNonce uint64) error {

	store := ctx.KVStore(k.storeKey)

	storePrefix := append(types.SignaturePrefix, prophecyID[:]...)

	prophecyInfo := types.ProphecyInfo{
		ProphecyId:           prophecyID,
		NetworkDescriptor:    networkDescriptor,
		CosmosSender:         cosmosSender,
		CosmosSenderSequence: cosmosSenderSequence,
		EthereumReceiver:     ethereumReceiver,
		TokenSymbol:          tokenSymbol,
		TokenAmount:          tokenAmount,
		DoublePeg:            doublePeg,
		GlobalNonce:          globalNonce,
		CrosschainFee:        crosschainFee,
		EthereumAddress:      []string{},
		Signatures:           []string{},
		BlockNumber:          uint64(k.currentHeight),
	}

	store.Set(storePrefix, k.cdc.MustMarshalBinaryBare(&prophecyInfo))
	return nil
}

// AppendSignature add a new ethereum address and signature to prophecy
func (k Keeper) AppendSignature(ctx sdk.Context, prophecyID []byte, ethereumAddress, signature string) error {
	store := ctx.KVStore(k.storeKey)

	prophecySignatures, ok := k.GetProphecyInfo(ctx, prophecyID)
	if !ok {
		return errors.New("can not get the prophecy")
	}

	prophecySignatures.EthereumAddress = append(prophecySignatures.EthereumAddress, ethereumAddress)
	prophecySignatures.Signatures = append(prophecySignatures.Signatures, signature)

	storePrefix := append(types.SignaturePrefix, prophecyID[:]...)

	store.Set(storePrefix, k.cdc.MustMarshalBinaryBare(&prophecySignatures))
	return nil
}

// Clean up outdated prophecy
func (k Keeper) CleanUpProphecy(ctx sdk.Context) {
	var prophecyInfo types.ProphecyInfo
	currentHeight := uint64(k.currentHeight)

	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.SignaturePrefix)
	for ; iter.Valid(); iter.Next() {
		k.cdc.MustUnmarshalBinaryBare(iter.Value(), &prophecyInfo)
		if prophecyInfo.BlockNumber-currentHeight > ProphecyLiftTime {
			storePrefix := append(types.SignaturePrefix, prophecyInfo.ProphecyId[:]...)
			store.Delete(storePrefix)
		}
	}
}
