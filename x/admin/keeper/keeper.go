package keeper

import (
	"strings"

	"github.com/Sifchain/sifnode/x/admin/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

type Keeper struct {
	cdc      codec.BinaryCodec
	storeKey sdk.StoreKey
}

func NewKeeper(cdc codec.Codec, storeKey sdk.StoreKey) Keeper {
	return Keeper{
		cdc:      cdc,
		storeKey: storeKey,
	}
}

func (k Keeper) InitGenesis(ctx sdk.Context, state types.GenesisState) []abci.ValidatorUpdate {
	if state.AdminAccounts != nil {
		for _, adminAccount := range state.AdminAccounts.AdminAccounts {
			k.SetAdminAccount(ctx, adminAccount)
		}
	}
	return []abci.ValidatorUpdate{}
}

func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	return &types.GenesisState{
		AdminAccounts: k.GetAdminAccounts(ctx),
	}
}

func (k Keeper) SetAdminAccount(ctx sdk.Context, account *types.AdminAccount) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetAdminAccountKey(*account)
	store.Set(key, k.cdc.MustMarshal(account))
}

func (k Keeper) RemoveAdminAccount(ctx sdk.Context, account *types.AdminAccount) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetAdminAccountKey(*account)
	store.Delete(key)
}

func (k Keeper) IsAdminAccount(ctx sdk.Context, adminType types.AdminType, adminAccount sdk.AccAddress) bool {
	accounts := k.GetAdminAccountsForType(ctx, adminType)
	if len(accounts.AdminAccounts) == 0 {
		return false
	}
	for _, account := range accounts.AdminAccounts {
		if strings.EqualFold(account.AdminAddress, adminAccount.String()) {
			return true
		}
	}
	return false
}

func (k Keeper) GetAdminAccountIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, types.AdminAccountStorePrefix)
}

func (k Keeper) GetAdminAccountsForType(ctx sdk.Context, adminType types.AdminType) *types.AdminAccounts {
	var res types.AdminAccounts
	iterator := k.GetAdminAccountIterator(ctx)
	defer func(iterator sdk.Iterator) {
		err := iterator.Close()
		if err != nil {
			panic(err)
		}
	}(iterator)
	for ; iterator.Valid(); iterator.Next() {
		var al types.AdminAccount
		bytesValue := iterator.Value()
		k.cdc.MustUnmarshal(bytesValue, &al)
		if al.AdminType == adminType {
			res.AdminAccounts = append(res.AdminAccounts, &al)
		}
	}
	return &res
}

func (k Keeper) GetAdminAccounts(ctx sdk.Context) *types.AdminAccounts {
	var res types.AdminAccounts
	iterator := k.GetAdminAccountIterator(ctx)
	defer func(iterator sdk.Iterator) {
		err := iterator.Close()
		if err != nil {
			panic(err)
		}
	}(iterator)
	for ; iterator.Valid(); iterator.Next() {
		var al types.AdminAccount
		bytesValue := iterator.Value()
		k.cdc.MustUnmarshal(bytesValue, &al)
		res.AdminAccounts = append(res.AdminAccounts, &al)
	}
	return &res
}
