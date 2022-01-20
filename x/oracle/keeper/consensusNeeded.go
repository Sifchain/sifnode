package keeper

import (
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/Sifchain/sifnode/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SetConsensusNeeded for a network.
func (k Keeper) SetConsensusNeeded(ctx sdk.Context,
	networkIdentity types.NetworkIdentity,
	consensusNeeded uint32) {
	store := ctx.KVStore(k.storeKey)
	key := networkIdentity.GetConsensusNeededPrefix()

	bs := make([]byte, 4)
	binary.BigEndian.PutUint32(bs, consensusNeeded)

	store.Set(key, bs)
}

// GetConsensusNeeded for a network
func (k Keeper) GetConsensusNeeded(ctx sdk.Context, networkIdentity types.NetworkIdentity) (uint32, error) {
	store := ctx.KVStore(k.storeKey)
	key := networkIdentity.GetConsensusNeededPrefix()

	if !store.Has(key) {
		return 0.0, fmt.Errorf("%s%s", "ConsensusNeeded not set for ", networkIdentity.NetworkDescriptor.String())
	}

	bz := store.Get(key)
	consensusNeeded := binary.BigEndian.Uint32(bz)
	if consensusNeeded > 100 {
		return 0, errors.New("consensusNeeded stored is too large")
	}
	return consensusNeeded, nil
}
