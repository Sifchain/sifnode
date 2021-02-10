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
	cdc *codec.Codec // The wire codec for binary encoding/decoding.

	supplyKeeper types.SupplyKeeper
	oracleKeeper types.OracleKeeper
}

// NewKeeper creates new instances of the oracle Keeper
func NewKeeper(cdc *codec.Codec, supplyKeeper types.SupplyKeeper, oracleKeeper types.OracleKeeper) Keeper {
	return Keeper{
		cdc:          cdc,
		supplyKeeper: supplyKeeper,
		oracleKeeper: oracleKeeper,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// ProcessClaim processes a new claim coming in from a validator
func (k Keeper) ProcessClaim(ctx sdk.Context, claim types.EthBridgeClaim) (oracle.Status, error) {
	oracleClaim, err := types.CreateOracleClaimFromEthClaim(k.cdc, claim)
	if err != nil {
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
		coins = sdk.Coins{sdk.NewInt64Coin(symbol, oracleClaim.Amount)}
		err = k.supplyKeeper.MintCoins(ctx, types.ModuleName, coins)
	case types.BurnText:
		coins = sdk.Coins{sdk.NewInt64Coin(oracleClaim.Symbol, oracleClaim.Amount)}
		err = k.supplyKeeper.MintCoins(ctx, types.ModuleName, coins)
	default:
		err = types.ErrInvalidClaimType
	}

	if err != nil {
		return err
	}

	if err := k.supplyKeeper.SendCoinsFromModuleToAccount(
		ctx, types.ModuleName, receiverAddress, coins,
	); err != nil {
		panic(err)
	}

	return nil
}

// ProcessBurn processes the burn of bridged coins from the given sender
func (k Keeper) ProcessBurn(ctx sdk.Context, cosmosSender sdk.AccAddress, amount sdk.Coins) error {
	if err := k.supplyKeeper.SendCoinsFromAccountToModule(
		ctx, cosmosSender, types.ModuleName, amount,
	); err != nil {
		return err
	}

	if err := k.supplyKeeper.BurnCoins(ctx, types.ModuleName, amount); err != nil {
		panic(err)
	}

	return nil
}

// ProcessLock processes the lockup of cosmos coins from the given sender
func (k Keeper) ProcessLock(ctx sdk.Context, cosmosSender sdk.AccAddress, amount sdk.Coins) error {
	return k.supplyKeeper.SendCoinsFromAccountToModule(ctx, cosmosSender, types.ModuleName, amount)
}
