package keeper

import (
	"encoding/binary"
	"errors"
	"github.com/Sifchain/sifnode/x/instrumentation"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/Sifchain/sifnode/x/oracle/types"
)

// ProphecyLiftTime is used to clean outdated prophecy info from keeper
const ProphecyLiftTime = 520000

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

	instrumentation.PeggyCheckpoint(ctx.Logger(), "SetProphecy", "prophecy", prophecy, "storePrefix", storePrefix)

	store.Set(storePrefix, k.cdc.MustMarshalBinaryBare(&prophecy))
}

// GetProphecyInfo return a prophecy's signatures
func (k Keeper) GetProphecyInfo(ctx sdk.Context, prophecyID []byte) (types.ProphecyInfo, bool) {
	var prophecySignatures types.ProphecyInfo
	store := ctx.KVStore(k.storeKey)

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
	tokenDenomHash string,
	tokenContractAddress string,
	tokenAmount sdk.Int,
	crosschainFee sdk.Int,
	doublePeg bool,
	globalSequence uint64) error {

	store := ctx.KVStore(k.storeKey)

	storePrefix := append(types.SignaturePrefix, prophecyID[:]...)

	prophecyInfo := types.ProphecyInfo{
		ProphecyId:           prophecyID,
		NetworkDescriptor:    networkDescriptor,
		CosmosSender:         cosmosSender,
		CosmosSenderSequence: cosmosSenderSequence,
		EthereumReceiver:     ethereumReceiver,
		TokenDenomHash:       tokenDenomHash,
		TokenContractAddress: tokenContractAddress,
		TokenAmount:          tokenAmount,
		DoublePeg:            doublePeg,
		GlobalSequence:       globalSequence,
		CrosschainFee:        crosschainFee,
		EthereumAddress:      []string{},
		Signatures:           []string{},
		BlockNumber:          uint64(k.currentHeight),
	}

	instrumentation.PeggyCheckpoint(ctx.Logger(), "SetProphecyInfo", prophecyInfo)

	k.SetGlobalNonceProphecyID(ctx, networkDescriptor, globalSequence, prophecyID)
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

	instrumentation.PeggyCheckpoint(ctx.Logger(), "AppendSignature", "storePrefix", storePrefix, "prophecySignatures", prophecySignatures)

	store.Set(storePrefix, k.cdc.MustMarshalBinaryBare(&prophecySignatures))
	return nil
}

// CleanUpProphecy clean up outdated prophecy
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
			storePrefix = k.getKeyViaNetworkDescriptorGlobalNonce(prophecyInfo.NetworkDescriptor, prophecyInfo.GlobalSequence)
			store.Delete(storePrefix)
		}
	}
}

// GetProphecyIDByNetworkDescriptorGlobalNonce get the prophecy id via network descriptor + global sequence
func (k Keeper) GetProphecyIDByNetworkDescriptorGlobalNonce(ctx sdk.Context,
	networkDescriptor types.NetworkDescriptor,
	globalSequence uint64) ([]byte, bool) {
	store := ctx.KVStore(k.storeKey)
	storeKey := k.getKeyViaNetworkDescriptorGlobalNonce(networkDescriptor, globalSequence)

	bz := store.Get(storeKey)
	if bz == nil {
		return bz, false
	}
	return bz, true
}

// SetGlobalNonceProphecyID store the map from network descriptor + global sequence to prophecy id
func (k Keeper) SetGlobalNonceProphecyID(ctx sdk.Context,
	networkDescriptor types.NetworkDescriptor,
	globalSequence uint64,
	prophecyID []byte) {
	store := ctx.KVStore(k.storeKey)
	storeKey := k.getKeyViaNetworkDescriptorGlobalNonce(networkDescriptor, globalSequence)

	instrumentation.PeggyCheckpoint(ctx.Logger(), "SetGlobalNonceProphecyID",
		"storeKey", storeKey,
		"prophecyID", prophecyID,
		"networkDescriptor", networkDescriptor,
		"globalSequence", globalSequence,
	)

	store.Set(storeKey, prophecyID)
}

func (k Keeper) getKeyViaNetworkDescriptorGlobalNonce(networkDescriptor types.NetworkDescriptor,
	globalSequence uint64) []byte {
	bs1 := make([]byte, 4)
	binary.LittleEndian.PutUint32(bs1, uint32(networkDescriptor))

	bs2 := make([]byte, 8)
	binary.LittleEndian.PutUint64(bs2, globalSequence)

	storeKey := append(append(types.GlobalNonceProphecyIDPrefix, bs1[:]...), bs2[:]...)
	return storeKey
}

// GetProphecyInfoWithScopeGlobalSequence get the prophecy id via network descriptor + global sequence
func (k Keeper) GetProphecyInfoWithScopeGlobalSequence(ctx sdk.Context,
	networkDescriptor types.NetworkDescriptor,
	startGlobalSequence uint64) []*types.ProphecyInfo {
	result := []*types.ProphecyInfo{}

	globalSequence := startGlobalSequence
	for {
		prophecyID, ok := k.GetProphecyIDByNetworkDescriptorGlobalNonce(ctx, networkDescriptor, globalSequence)
		if !ok {
			return result
		}

		prophecy, ok := k.GetProphecy(ctx, prophecyID)
		if !ok {
			return result
		}

		if prophecy.Status != types.StatusText_STATUS_TEXT_SUCCESS {
			return result
		}

		prophecyInfo, ok := k.GetProphecyInfo(ctx, prophecyID)
		if !ok {
			return result
		}
		globalSequence++
		result = append(result, &prophecyInfo)
	}

}
