package keeper

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/Sifchain/sifnode/x/ethbridge/types"
	"github.com/Sifchain/sifnode/x/oracle"
)

// Keeper maintains the link to data storage and
// exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
	cdc codec.BinaryMarshaler // The wire codec for binary encoding/decoding.

	bankKeeper   types.BankKeeper
	oracleKeeper types.OracleKeeper
	storeKey     sdk.StoreKey
}

// NewKeeper creates new instances of the oracle Keeper
func NewKeeper(cdc codec.BinaryMarshaler, bankKeeper types.BankKeeper, oracleKeeper types.OracleKeeper, storeKey sdk.StoreKey) Keeper {
	return Keeper{
		cdc:          cdc,
		bankKeeper:   bankKeeper,
		oracleKeeper: oracleKeeper,
		storeKey:     storeKey,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// ProcessClaim processes a new claim coming in from a validator
func (k Keeper) ProcessClaim(ctx sdk.Context, claim types.EthBridgeClaim) (oracle.Status, error) {
	fmt.Println("sifnode ethbridge keeper ProcessClaim")
	oracleClaim, err := types.CreateOracleClaimFromEthClaim(k.cdc, claim)
	if err != nil {
		fmt.Printf("sifnode ethbridge keeper ProcessClaim oracle %s\n", err.Error())
		return oracle.Status{}, err
	}

	return k.oracleKeeper.ProcessClaim(ctx, oracleClaim)
}

// ProcessSuccessfulClaim processes a claim that has just completed successfully with consensus
func (k Keeper) ProcessSuccessfulClaim(ctx sdk.Context, claim string) error {
	oracleClaim, err := types.CreateOracleClaimFromOracleString(claim)
	if err != nil {
		return err
	}

	receiverAddress := oracleClaim.CosmosReceiver

	var coins sdk.Coins
	switch oracleClaim.ClaimType {
	case types.LockText:
		symbol := fmt.Sprintf("%v%v", types.PeggedCoinPrefix, oracleClaim.Symbol)
		k.AddPeggyToken(ctx, symbol)

		coins = sdk.Coins{sdk.NewCoin(symbol, oracleClaim.Amount)}
		err = k.bankKeeper.MintCoins(ctx, types.ModuleName, coins)
	case types.BurnText:
		coins = sdk.Coins{sdk.NewCoin(oracleClaim.Symbol, oracleClaim.Amount)}
		err = k.bankKeeper.MintCoins(ctx, types.ModuleName, coins)
	default:
		err = types.ErrInvalidClaimType
	}

	if err != nil {
		return err
	}

	if err := k.bankKeeper.SendCoinsFromModuleToAccount(
		ctx, types.ModuleName, receiverAddress, coins,
	); err != nil {
		panic(err)
	}

	return nil
}

// ProcessBurn processes the burn of bridged coins from the given sender
func (k Keeper) ProcessBurn(ctx sdk.Context, cosmosSender sdk.AccAddress, amount sdk.Coins) error {
	if err := k.bankKeeper.SendCoinsFromAccountToModule(
		ctx, cosmosSender, types.ModuleName, amount,
	); err != nil {
		return err
	}

	if err := k.bankKeeper.BurnCoins(ctx, types.ModuleName, amount); err != nil {
		panic(err)
	}

	return nil
}

// ProcessLock processes the lockup of cosmos coins from the given sender
func (k Keeper) ProcessLock(ctx sdk.Context, cosmosSender sdk.AccAddress, amount sdk.Coins) error {
	return k.bankKeeper.SendCoinsFromAccountToModule(ctx, cosmosSender, types.ModuleName, amount)
}

// ProcessUpdateWhiteListValidator processes the update whitelist validator from admin
func (k Keeper) ProcessUpdateWhiteListValidator(ctx sdk.Context, cosmosSender sdk.AccAddress, validator sdk.ValAddress, operationtype string) error {
	return k.oracleKeeper.ProcessUpdateWhiteListValidator(ctx, cosmosSender, validator, operationtype)
}

// Exists chec if the key existed in db.
func (k Keeper) Exists(ctx sdk.Context, key []byte) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(key)
}
