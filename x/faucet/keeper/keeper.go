package keeper

import (
	"fmt"

	"github.com/Sifchain/sifnode/x/faucet/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Keeper of the faucet store
type Keeper struct {
	storeKey     sdk.StoreKey
	cdc          *codec.Codec
	supplyKeeper types.SupplyKeeper
	bankKeeper   types.BankKeeper
}

// NewKeeper creates a faucet keeper
func NewKeeper(supplyKeeper types.SupplyKeeper, cdc *codec.Codec, key sdk.StoreKey, bankKeeper types.BankKeeper) Keeper {
	keeper := Keeper{
		supplyKeeper: supplyKeeper,
		bankKeeper:   bankKeeper,
		storeKey:     key,
		cdc:          cdc,
		// paramspace: paramspace.WithKeyTable(types.ParamKeyTable()),
	}
	return keeper
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) GetBankKeeper() types.BankKeeper {
	return k.bankKeeper
}

func (k Keeper) GetSupplyKeeper() types.SupplyKeeper {
	return k.supplyKeeper
}

// GetWithdrawnAmountInEpoch validates if a user has utilized faucet functionality within the last 4 hours
// If a withdrawal action has occurred the module will block withdraws until the timer is reset
func (k Keeper) GetWithdrawnAmountInEpoch(ctx sdk.Context, user string, token string) (sdk.Int, error) {
	faucetTracker := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.FaucetPrefix))
	ok := faucetTracker.Has(types.GetBalanceKey(user, token))
	if !ok {
		return sdk.ZeroInt(), nil
	}
	amount := faucetTracker.Get(types.GetBalanceKey(user, token))
	var am sdk.Int
	err := k.cdc.UnmarshalBinaryBare(amount, &am)
	if err != nil {
		return sdk.ZeroInt(), err
	}
	return am, nil
}

func (k Keeper) SetWithdrawnAmountInEpoch(ctx sdk.Context, user string, amount sdk.Int, token string) error {
	faucetTracker := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.FaucetPrefix))
	withdrawnAmount, err := k.GetWithdrawnAmountInEpoch(ctx, user, token)
	if err != nil {
		return err
	}
	totalAmountInEpoch := withdrawnAmount.Add(amount)
	bz, err := k.cdc.MarshalBinaryBare(&totalAmountInEpoch)
	if err != nil {
		return err
	}
	faucetTracker.Set(types.GetBalanceKey(user, token), bz)
	return nil
}

func (k Keeper) StartNextEpoch(ctx sdk.Context) {
	faucetTracker := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(faucetTracker, types.KeyPrefix(types.FaucetPrefix))
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		faucetTracker.Delete(iterator.Key())
	}
}

func (k Keeper) CanRequest(ctx sdk.Context, user string, coins sdk.Coins) (bool, error) {
	alreadyWithraw, err := k.GetWithdrawnAmountInEpoch(ctx, user, types.FaucetToken)
	if err != nil {
		return false, err
	}
	maxAllowedWithdraw, ok := sdk.NewIntFromString(types.MaxWithdrawAmountPerEpoch)
	if !ok {
		return false, nil
	}
	amount := coins.AmountOf(types.FaucetToken)
	if alreadyWithraw.Add(amount).GT(maxAllowedWithdraw) {
		return false, types.ErrInvalid
	}
	return true, nil
}

func (k Keeper) ExecuteRequest(ctx sdk.Context, user string, coins sdk.Coins) (bool, error) {
	amount := coins.AmountOf(types.FaucetToken)
	err := k.SetWithdrawnAmountInEpoch(ctx, user, amount, types.FaucetToken)
	if err != nil {
		return false, err
	}
	return true, nil
}
