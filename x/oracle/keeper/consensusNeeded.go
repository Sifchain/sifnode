package keeper

import (
	"encoding/binary"
	"fmt"

	"github.com/Sifchain/sifnode/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SetCrossChainFee set the crosschain fee for a network.
func (k Keeper) SetConsensusNeeded(ctx sdk.Context,
	networkIdentity types.NetworkIdentity,
	consensusNeeded uint32) {
	store := ctx.KVStore(k.storeKey)
	key := networkIdentity.GetConsensusNeededPrefix()

	fmt.Println("++++++ SetConsensusNeeded")

	bs := make([]byte, 4)
	binary.BigEndian.PutUint32(bs, consensusNeeded)

	store.Set(key, bs)
}

// GetCrossChainFeeConfig return crosschain fee config
func (k Keeper) GetConsensusNeeded(ctx sdk.Context, networkIdentity types.NetworkIdentity) (uint32, error) {
	store := ctx.KVStore(k.storeKey)
	key := networkIdentity.GetConsensusNeededPrefix()

	if !store.Has(key) {
		return 0.0, fmt.Errorf("%s%s", "ConsensusNeeded not set for ", networkIdentity.NetworkDescriptor.String())
	}

	bz := store.Get(key)
	return binary.BigEndian.Uint32(bz), nil
}
