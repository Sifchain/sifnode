package keeper

import (
	"errors"

	"github.com/Sifchain/sifnode/x/instrumentation"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/Sifchain/sifnode/x/oracle/types"
)

const (
	// ProphecyLifeTime is used to clean outdated prophecy info from keeper
	ProphecyLifeTime = 520000
	// Max prophecy returned in one query
	MaxProphecyQueryResult = 10
	// Clean up outdated prophecies every 1000 blocks
	CleanUpFrequency = 1000
)

// GetProphecies returns all prophecies
func (k Keeper) GetProphecies(ctx sdk.Context) []types.Prophecy {
	var prophecies []types.Prophecy
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.ProphecyPrefix)
	for ; iter.Valid(); iter.Next() {
		var prophecy types.Prophecy
		k.cdc.MustUnmarshal(iter.Value(), &prophecy)
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
	k.cdc.MustUnmarshal(bz, &prophecy)

	return prophecy, true
}

// SetProphecy saves a prophecy with an initial claim
func (k Keeper) SetProphecy(ctx sdk.Context, prophecy types.Prophecy) {
	store := ctx.KVStore(k.storeKey)

	storePrefix := append(types.ProphecyPrefix, prophecy.Id[:]...)

	instrumentation.PeggyCheckpoint(ctx.Logger(), instrumentation.SetProphecy, "prophecy", prophecy, "validator length", prophecy.ClaimValidators, "storePrefix", string(storePrefix))

	store.Set(storePrefix, k.cdc.MustMarshal(&prophecy))
}

// GetProphecyInfo return a prophecy's signatures
func (k Keeper) GetProphecyInfo(ctx sdk.Context, prophecyID []byte) (types.ProphecyInfo, bool) {
	var prophecySignatures types.ProphecyInfo
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(append(types.SignaturePrefix, prophecyID[:]...))
	if bz == nil {
		return types.ProphecyInfo{}, false
	}

	k.cdc.MustUnmarshal(bz, &prophecySignatures)

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
	crossChainFee sdk.Int,
	bridgeToken bool,
	globalSequence uint64,
	tokenDecimal uint8,
	tokenName string,
	tokenSymbol string) error {

	prophecyInfo := types.ProphecyInfo{
		ProphecyId:           prophecyID,
		NetworkDescriptor:    networkDescriptor,
		CosmosSender:         cosmosSender,
		CosmosSenderSequence: cosmosSenderSequence,
		EthereumReceiver:     ethereumReceiver,
		TokenDenomHash:       tokenDenomHash,
		TokenContractAddress: tokenContractAddress,
		TokenAmount:          tokenAmount,
		BridgeToken:          bridgeToken,
		GlobalSequence:       globalSequence,
		CrosschainFee:        crossChainFee,
		EthereumAddress:      []string{},
		Signatures:           []string{},
		BlockNumber:          uint64(k.currentHeight),
		TokenName:            tokenName,
		TokenSymbol:          tokenSymbol,
		Decimals:             uint32(tokenDecimal),
	}

	instrumentation.PeggyCheckpoint(ctx.Logger(), instrumentation.SetProphecyInfo, prophecyInfo)

	k.SetProphecyInfoObj(ctx, &prophecyInfo)
	return nil
}

func (k Keeper) SetProphecyInfoObj(ctx sdk.Context, prophecyInfo *types.ProphecyInfo) {
	store := ctx.KVStore(k.storeKey)
	storePrefix := append(types.SignaturePrefix, prophecyInfo.ProphecyId[:]...)

	k.SetGlobalNonceProphecyID(ctx, prophecyInfo.NetworkDescriptor, prophecyInfo.GlobalSequence, prophecyInfo.ProphecyId)
	store.Set(storePrefix, k.cdc.MustMarshal(prophecyInfo))
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

	instrumentation.PeggyCheckpoint(ctx.Logger(), instrumentation.AppendSignature, "storePrefix", storePrefix, "prophecySignatures", prophecySignatures)

	store.Set(storePrefix, k.cdc.MustMarshal(&prophecySignatures))
	return nil
}

// CleanUpProphecy clean up outdated prophecy, not check the status of prophecy
// since ProphecyLifeTime is big enough for relayers to handle prophecy
func (k Keeper) CleanUpProphecy(ctx sdk.Context) {
	// it is low efficient to check outdated prophecy each block
	if k.currentHeight%CleanUpFrequency != 0 {
		return
	}
	var prophecyInfo types.ProphecyInfo
	currentHeight := uint64(k.currentHeight)

	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.SignaturePrefix)
	for ; iter.Valid(); iter.Next() {
		k.cdc.MustUnmarshal(iter.Value(), &prophecyInfo)
		if currentHeight > prophecyInfo.BlockNumber+ProphecyLifeTime {
			k.DeleteProphecyInfo(ctx, prophecyInfo)
		}
	}
}

// DeleteProphecyInfo remove both signatures and global sequence in keeper
func (k Keeper) DeleteProphecyInfo(ctx sdk.Context, prophecyInfo types.ProphecyInfo) {
	storePrefix := prophecyInfo.GetSignaturePrefix()
	store := ctx.KVStore(k.storeKey)
	store.Delete(storePrefix)
	storePrefix = k.getKeyViaNetworkDescriptorGlobalNonce(prophecyInfo.NetworkDescriptor, prophecyInfo.GlobalSequence)
	store.Delete(storePrefix)
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

	instrumentation.PeggyCheckpoint(ctx.Logger(), instrumentation.SetGlobalNonceProphecyID,
		"storeKey", storeKey,
		"prophecyID", prophecyID,
		"networkDescriptor", networkDescriptor,
		"globalSequence", globalSequence,
	)

	store.Set(storeKey, prophecyID)
}

func (k Keeper) getKeyViaNetworkDescriptorGlobalNonce(networkDescriptor types.NetworkDescriptor,
	globalSequence uint64) []byte {

	bs := k.cdc.MustMarshal(&types.GlobalSequenceKey{
		NetworkDescriptor: networkDescriptor,
		GlobalSequence:    globalSequence,
	})

	storeKey := append(types.GlobalSequenceProphecyIDPrefix, bs[:]...)
	return storeKey
}

// GetProphecyInfoWithScopeGlobalSequence get the prophecy id via network descriptor + global sequence
func (k Keeper) GetProphecyInfoWithScopeGlobalSequence(ctx sdk.Context,
	networkDescriptor types.NetworkDescriptor,
	startGlobalSequence uint64) []*types.ProphecyInfo {
	result := []*types.ProphecyInfo{}

	globalSequence := startGlobalSequence
	for i := 0; i < MaxProphecyQueryResult; i++ {
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
	return result
}

// GetProphecyInfo return a prophecy's signatures
func (k Keeper) GetAllProphecyInfo(ctx sdk.Context) []*types.GenesisProphecyInfo {
	store := ctx.KVStore(k.storeKey)
	prophecyInfos := make([]*types.GenesisProphecyInfo, 0)

	iterator := sdk.KVStorePrefixIterator(store, types.GlobalSequenceProphecyIDPrefix)

	defer func(iterator sdk.Iterator) {
		err := iterator.Close()
		if err != nil {
			panic(err)
		}
	}(iterator)
	for ; iterator.Valid(); iterator.Next() {
		var globalSequenceKey types.GlobalSequenceKey
		var prophecyInfo types.ProphecyInfo

		globalSequenceKey, err := types.GetGlobalSequenceKeyKeyFromRawKey(k.cdc, iterator.Key())
		if err != nil {
			panic(err)
		}
		prophecyId := iterator.Value()

		storePrefix := append(types.SignaturePrefix, prophecyId[:]...)
		value := store.Get(storePrefix)

		k.cdc.MustUnmarshal(value, &prophecyInfo)
		prophecyInfos = append(prophecyInfos, &types.GenesisProphecyInfo{
			GlobalSequenceKey: &globalSequenceKey,
			ProphecyInfo:      &prophecyInfo,
		})
	}

	return prophecyInfos
}
