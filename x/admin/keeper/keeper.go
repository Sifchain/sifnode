package keeper

import (
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
	for _, adminAccount := range state.AdminAccounts {
		k.SetAdminAccount(ctx, adminAccount)
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
	if len(accounts) == 0 {
		return false
	}
	for _, account := range accounts {
		if types.StringCompare(account.AdminAddress, adminAccount.String()) {
			return true
		}
	}
	return false
}

func (k Keeper) GetAdminAccountIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, types.AdminAccountStorePrefix)
}

func (k Keeper) GetAdminAccountsForType(ctx sdk.Context, adminType types.AdminType) []*types.AdminAccount {
	var res []*types.AdminAccount
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
			res = append(res, &al)
		}
	}
	return res
}

func (k Keeper) GetAdminAccounts(ctx sdk.Context) []*types.AdminAccount {
	var accounts []*types.AdminAccount
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
		accounts = append(accounts, &al)
	}
	return accounts
}

func (k Keeper) SetParams(ctx sdk.Context, params *types.Params) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.ParamsStorePrefix, k.cdc.MustMarshal(params))
}

func (k Keeper) GetParams(ctx sdk.Context) *types.Params {
	defaultSubmitProposalFee := sdk.NewUintFromString("5000000000000000000000") // 5000

	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.ParamsStorePrefix)
	if bz == nil {
		return &types.Params{SubmitProposalFee: defaultSubmitProposalFee}
	}
	var params types.Params
	k.cdc.MustUnmarshal(bz, &params)
	return &params
}
