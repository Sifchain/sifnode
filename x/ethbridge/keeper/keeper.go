package keeper

import (
	"errors"
	"fmt"
	oracletypes "github.com/Sifchain/sifnode/x/oracle/types"
	tokenregistrytypes "github.com/Sifchain/sifnode/x/tokenregistry/types"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/Sifchain/sifnode/x/ethbridge/types"
)

const errorMessageKey = "errorMessageKey"

// Keeper maintains the link to data storage and
// exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
	cdc                 codec.BinaryCodec // The wire codec for binary encoding/decoding.
	accountKeeper       types.AccountKeeper
	bankKeeper          types.BankKeeper
	oracleKeeper        types.OracleKeeper
	tokenRegistryKeeper types.TokenRegistryKeeper
	storeKey            sdk.StoreKey
}

// NewKeeper creates new instances of the oracle Keeper
func NewKeeper(cdc codec.BinaryCodec, bankKeeper types.BankKeeper, oracleKeeper types.OracleKeeper, accountKeeper types.AccountKeeper, tokenRegistryKeeper tokenregistrytypes.Keeper, storeKey sdk.StoreKey) Keeper {
	return Keeper{
		cdc:                 cdc,
		accountKeeper:       accountKeeper,
		bankKeeper:          bankKeeper,
		oracleKeeper:        oracleKeeper,
		tokenRegistryKeeper: tokenRegistryKeeper,
		storeKey:            storeKey,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// ProcessClaim processes a new claim coming in from a validator
func (k Keeper) ProcessClaim(ctx sdk.Context, claim *types.EthBridgeClaim) (oracletypes.Status, error) {
	logger := k.Logger(ctx)
	oracleClaim, err := types.CreateOracleClaimFromEthClaim(claim)
	if err != nil {
		logger.Error("failed to create oracle claim from eth claim.",
			errorMessageKey, err.Error())
		return oracletypes.Status{}, err
	}
	return k.oracleKeeper.ProcessClaim(ctx, oracleClaim)
}

// ProcessSuccessfulClaim processes a claim that has just completed successfully with consensus
func (k Keeper) ProcessSuccessfulClaim(ctx sdk.Context, claim string) error {
	logger := k.Logger(ctx)
	oracleClaim, err := types.CreateOracleClaimFromOracleString(claim)
	if err != nil {
		logger.Error("failed to create oracle claim from oracle string.",
			errorMessageKey, err.Error())
		return err
	}
	receiverAddress := oracleClaim.CosmosReceiver
	var coins sdk.Coins
	switch oracleClaim.ClaimType {
	case types.ClaimType_CLAIM_TYPE_LOCK:
		symbol := fmt.Sprintf("%v%v", types.PeggedCoinPrefix, oracleClaim.Symbol)
		k.AddPeggyToken(ctx, symbol)
		coins = sdk.NewCoins(sdk.NewCoin(symbol, oracleClaim.Amount))
		err = k.bankKeeper.MintCoins(ctx, types.ModuleName, coins)
	case types.ClaimType_CLAIM_TYPE_BURN:
		coins = sdk.NewCoins(sdk.NewCoin(oracleClaim.Symbol, oracleClaim.Amount))
		err = k.bankKeeper.MintCoins(ctx, types.ModuleName, coins)
	default:
		err = types.ErrInvalidClaimType
	}
	if err != nil {
		logger.Error("failed to process successful claim.",
			errorMessageKey, err.Error())
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
func (k Keeper) ProcessBurn(ctx sdk.Context, cosmosSender sdk.AccAddress, msg *types.MsgBurn) error {
	logger := k.Logger(ctx)
	var coins sdk.Coins

	if k.IsBlacklisted(ctx, msg.EthereumReceiver) {
		return types.ErrInvalidEthAddress
	}

	if k.IsCethReceiverAccountSet(ctx) {
		coins = sdk.NewCoins(sdk.NewCoin(types.CethSymbol, msg.CethAmount))
		err := k.bankKeeper.SendCoins(ctx, cosmosSender, k.GetCethReceiverAccount(ctx), coins)
		if err != nil {
			logger.Error("failed to send ceth from account to account.",
				errorMessageKey, err.Error())
			return err
		}
		coins = sdk.NewCoins(sdk.NewCoin(msg.Symbol, msg.Amount))
	} else {
		if msg.Symbol == types.CethSymbol {
			coins = sdk.NewCoins(sdk.NewCoin(types.CethSymbol, msg.CethAmount.Add(msg.Amount)))
		} else {
			coins = sdk.NewCoins(sdk.NewCoin(msg.Symbol, msg.Amount), sdk.NewCoin(types.CethSymbol, msg.CethAmount))
		}
	}
	err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, cosmosSender, types.ModuleName, coins)
	if err != nil {
		logger.Error("failed to send ceth from module to account.",
			errorMessageKey, err.Error())
		return err
	}
	coins = sdk.NewCoins(sdk.NewCoin(msg.Symbol, msg.Amount))
	err = k.bankKeeper.BurnCoins(ctx, types.ModuleName, coins)
	if err != nil {
		logger.Error("failed to burn locked coin.",
			errorMessageKey, err.Error())
		return err
	}
	return nil
}

// ProcessLock processes the lockup of cosmos coins from the given sender
func (k Keeper) ProcessLock(ctx sdk.Context, cosmosSender sdk.AccAddress, msg *types.MsgLock) error {
	logger := k.Logger(ctx)

	if k.IsBlacklisted(ctx, msg.EthereumReceiver) {
		return types.ErrInvalidEthAddress
	}

	var coins sdk.Coins
	if k.IsCethReceiverAccountSet(ctx) {
		coins = sdk.NewCoins(sdk.NewCoin(types.CethSymbol, msg.CethAmount))
		err := k.bankKeeper.SendCoins(ctx, cosmosSender, k.GetCethReceiverAccount(ctx), coins)
		if err != nil {
			logger.Error("failed to send ceth from account to account.",
				errorMessageKey, err.Error())
			return err
		}
		coins = sdk.NewCoins(sdk.NewCoin(msg.Symbol, msg.Amount))
	} else {
		coins = sdk.NewCoins(sdk.NewCoin(msg.Symbol, msg.Amount), sdk.NewCoin(types.CethSymbol, msg.CethAmount))
	}
	err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, cosmosSender, types.ModuleName, coins)
	if err != nil {
		logger.Error("failed to transfer coin from account to module.", errorMessageKey, err.Error())
		return err
	}
	coins = sdk.NewCoins(sdk.NewCoin(msg.Symbol, msg.Amount))
	err = k.bankKeeper.BurnCoins(ctx, types.ModuleName, coins)
	if err != nil {
		logger.Error("failed to burn burned coin.", errorMessageKey, err.Error())
		return err
	}
	return nil
}

// ProcessUpdateWhiteListValidator processes the update whitelist validator from admin
func (k Keeper) ProcessUpdateWhiteListValidator(ctx sdk.Context, cosmosSender sdk.AccAddress, validator sdk.ValAddress, operationtype string) error {
	return k.oracleKeeper.ProcessUpdateWhiteListValidator(ctx, cosmosSender, validator, operationtype)
}

// ProcessUpdateCethReceiverAccount processes the update whitelist validator from admin
func (k Keeper) ProcessUpdateCethReceiverAccount(ctx sdk.Context, cosmosSender sdk.AccAddress, cethReceiverAccount sdk.AccAddress) error {
	logger := k.Logger(ctx)
	if !k.oracleKeeper.IsAdminAccount(ctx, cosmosSender) {
		logger.Error("cosmos sender is not admin account.")
		return errors.New("only admin account can update ceth receiver account")
	}
	k.SetCethReceiverAccount(ctx, cethReceiverAccount)
	return nil
}

// ProcessRescueCeth transfer ceth from ethbridge module to an account
func (k Keeper) ProcessRescueCeth(ctx sdk.Context, msg *types.MsgRescueCeth) error {
	logger := k.Logger(ctx)
	cosmosSender, err := sdk.AccAddressFromBech32(msg.CosmosSender)
	if err != nil {
		return err
	}
	cosmosReceiver, err := sdk.AccAddressFromBech32(msg.CosmosReceiver)
	if err != nil {
		return err
	}
	if !k.oracleKeeper.IsAdminAccount(ctx, cosmosSender) {
		logger.Error("cosmos sender is not admin account.")
		return errors.New("only admin account can call rescue ceth")
	}
	coins := sdk.NewCoins(sdk.NewCoin(types.CethSymbol, msg.CethAmount))
	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, cosmosReceiver, coins)
	if err != nil {
		logger.Error("failed to transfer coin from module to account.",
			errorMessageKey, err.Error())
		return err
	}
	return nil
}

// Exists chec if the key existed in db.
func (k Keeper) Exists(ctx sdk.Context, key []byte) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(key)
}
