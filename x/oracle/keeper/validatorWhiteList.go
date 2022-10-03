package keeper

import (
	"bytes"

	"github.com/Sifchain/sifnode/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetAllWhiteList get the validators for all networks.
func (k Keeper) GetAllWhiteList(ctx sdk.Context) []*types.GenesisValidatorWhiteList {
	genesisValidatorWhiteList := make([]*types.GenesisValidatorWhiteList, 0)
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.WhiteListValidatorPrefix)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		networkIdentity, err := types.GetFromPrefix(k.cdc, iterator.Key(), types.WhiteListValidatorPrefix)
		if err != nil {
			panic("key for validator whitelist is invalid")
		}
		var validatorWhiteList types.ValidatorWhiteList
		k.cdc.MustUnmarshal(iterator.Value(), &validatorWhiteList)

		genesisValidatorWhiteList = append(genesisValidatorWhiteList, &types.GenesisValidatorWhiteList{
			NetworkDescriptor:  networkIdentity.NetworkDescriptor,
			ValidatorWhitelist: &validatorWhiteList,
		})

	}

	return genesisValidatorWhiteList
}

// SetOracleWhiteList set the validator list for a network.
// func (k Keeper) SetOracleWhiteList(ctx sdk.Context, networkDescriptor types.NetworkIdentity, validatorList types.ValidatorWhiteList) {
// 	store := ctx.KVStore(k.storeKey)
// 	key := networkDescriptor.GetPrefix(k.cdc)
// 	store.Set(key, k.cdc.MustMarshal(&validatorList))
// }

// RemoveOracleWhiteList remove the validator list for a network.
func (k Keeper) RemoveOracleWhiteList(ctx sdk.Context, networkDescriptor types.NetworkIdentity) {
	store := ctx.KVStore(k.storeKey)
	key := networkDescriptor.GetPrefix(k.cdc)
	store.Delete(key)
}

// ExistsOracleWhiteList check if the key exist
func (k Keeper) ExistsOracleWhiteList(ctx sdk.Context, networkDescriptor types.NetworkIdentity) bool {
	key := networkDescriptor.GetPrefix(k.cdc)
	return k.Exists(ctx, key)
}

// GetOracleWhiteList return validator list
func (k Keeper) GetOracleWhiteList(ctx sdk.Context, networkIdentity types.NetworkIdentity) types.ValidatorWhiteList {
	store := ctx.KVStore(k.storeKey)

	key := k.GetWhiteListKey(networkIdentity)
	value := store.Get(key)
	var whiteList types.ValidatorWhiteList
	k.cdc.MustUnmarshal(value, &whiteList)
	return whiteList
}

// // GetAllValidators return validator list
// func (k Keeper) GetAllValidators(ctx sdk.Context, networkIdentity types.NetworkIdentity) []sdk.ValAddress {
// 	valAddresses := k.GetOracleWhiteList(ctx, networkIdentity)

// 	vl := []sdk.ValAddress{}
// 	for i, power := range valAddresses.GetWhiteList() {
// 		addr, err := sdk.ValAddressFromBech32(i)
// 		if err != nil {
// 			panic(err)
// 		}
// 		if power > 0 {
// 			vl = append(vl, addr)
// 		}
// 	}

// 	return vl
// }

// ValidateAddress is a validator in whitelist
func (k Keeper) ValidateAddress(ctx sdk.Context, networkIdentity types.NetworkIdentity, address sdk.ValAddress) bool {
	if !k.ExistsOracleWhiteList(ctx, networkIdentity) {
		return false
	}
	whiteList := k.GetOracleWhiteList(ctx, networkIdentity)

	for _, value := range whiteList.ValidatorPower {
		if bytes.Compare(value.ValidatorAddress, address) == 0 {
			if value.VotingPower > 0 {
				return true
			}
		}

	}

	return false
}

// GetValidatorPower return validator's power
// func (k Keeper) GetValidatorPowerMap(ctx sdk.Context, networkDescriptor types.NetworkDescriptor) map[string]types.VotingPower {
// 	powerMap := make(map[string]types.VotingPower)
// 	store := ctx.KVStore(k.storeKey)

// 	iterator := sdk.KVStorePrefixIterator(store, types.WhiteListValidatorPrefix)
// 	defer iterator.Close()
// 	for ; iterator.Valid(); iterator.Next() {

// 		votingKey, err := types.GetValidatorVotingPowerKeyFromRawKey(k.cdc, iterator.Key())
// 		if err != nil {
// 			panic(err.Error())
// 		}

// 		var votingPower types.VotingPower
// 		k.cdc.MustUnmarshal(iterator.Value(), &votingPower)

// 		validatorWhiteList = append(validatorWhiteList, &types.ValidatorWhiteList{
// 			VotingKey:   &votingKey,
// 			VotingPower: &votingPower,
// 		})
// 	}

// 	return 0
// }

// UpdateOracleWhiteList validator's power
func (k Keeper) UpdateOracleWhiteList(ctx sdk.Context, networkDescriptor types.NetworkDescriptor, validator sdk.ValAddress, power uint32) {
	store := ctx.KVStore(k.storeKey)

	key := k.GetWhiteListKey(types.NewNetworkIdentity(networkDescriptor))
	value := store.Get(key)

	var validatorWhiteList types.ValidatorWhiteList
	k.cdc.MustUnmarshal(value, &validatorWhiteList)

	validatorWhiteList.UpdateValidatorPower(validator, power)
	store.Set(key, k.cdc.MustMarshal(&validatorWhiteList))
}

// GetWhiteListKey get the key for storage, key = WhiteListValidatorPrefix + networkDescriptor
func (k Keeper) GetWhiteListKey(networkIdentity types.NetworkIdentity) []byte {
	buf := k.cdc.MustMarshal(&networkIdentity)
	return append(types.WhiteListValidatorPrefix, buf...)
}
