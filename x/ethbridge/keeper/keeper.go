package keeper

import (
	"errors"
	"fmt"

	oracletypes "github.com/Sifchain/sifnode/x/oracle/types"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/Sifchain/sifnode/x/ethbridge/types"
)

const errorMessageKey = "errorMessageKey"

// Keeper maintains the link to data storage and
// exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
	cdc codec.BinaryMarshaler // The wire codec for binary encoding/decoding.

	accountKeeper types.AccountKeeper
	bankKeeper    types.BankKeeper
	oracleKeeper  types.OracleKeeper
	storeKey      sdk.StoreKey
}

// GetAccountKeeper
func (k Keeper) GetAccountKeeper() types.AccountKeeper {
	return k.accountKeeper
}

// GetBankKeeper
func (k Keeper) GetBankKeeper() types.BankKeeper {
	return k.bankKeeper
}

// NewKeeper creates new instances of the oracle Keeper
func NewKeeper(cdc codec.BinaryMarshaler, bankKeeper types.BankKeeper, oracleKeeper types.OracleKeeper, accountKeeper types.AccountKeeper, storeKey sdk.StoreKey) Keeper {
	return Keeper{
		cdc:           cdc,
		accountKeeper: accountKeeper,
		bankKeeper:    bankKeeper,
		oracleKeeper:  oracleKeeper,
		storeKey:      storeKey,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// ProcessClaim processes a new claim coming in from a validator
func (k Keeper) ProcessClaim(ctx sdk.Context, claim *types.EthBridgeClaim) (oracletypes.StatusText, error) {
	return k.oracleKeeper.ProcessClaim(ctx, claim.NetworkDescriptor, claim.GetProphecyID(), claim.ValidatorAddress)
}

// ProcessSuccessfulClaim processes a claim that has just completed successfully with consensus
func (k Keeper) ProcessSuccessfulClaim(ctx sdk.Context, claim *types.EthBridgeClaim) error {
	logger := k.Logger(ctx)

	var coins sdk.Coins
	var err error
	switch claim.ClaimType {
	case types.ClaimType_CLAIM_TYPE_LOCK:
		symbol := fmt.Sprintf("%v%v", types.PeggedCoinPrefix, claim.Symbol)
		k.AddPeggyToken(ctx, symbol)

		coins = sdk.Coins{sdk.NewCoin(symbol, claim.Amount)}
		err = k.bankKeeper.MintCoins(ctx, types.ModuleName, coins)
	case types.ClaimType_CLAIM_TYPE_BURN:
		coins = sdk.Coins{sdk.NewCoin(claim.Symbol, claim.Amount)}
		err = k.bankKeeper.MintCoins(ctx, types.ModuleName, coins)
	default:
		err = types.ErrInvalidClaimType
	}

	if err != nil {
		logger.Error("failed to process successful claim.",
			errorMessageKey, err.Error())
		return err
	}

	receiverAddress, err := sdk.AccAddressFromBech32(claim.CosmosReceiver)

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
func (k Keeper) ProcessBurn(ctx sdk.Context, cosmosSender sdk.AccAddress, senderSequence uint64, msg *types.MsgBurn) ([]byte, error) {
	logger := k.Logger(ctx)
	var coins sdk.Coins
	networkIdentity := oracletypes.NewNetworkIdentity(msg.NetworkDescriptor)
	crossChainFeeConfig, err := k.oracleKeeper.GetCrossChainFeeConfig(ctx, networkIdentity)

	if err != nil {
		return []byte{}, err
	}

	minimumBurn := crossChainFeeConfig.MinimumBurnCost.Mul(crossChainFeeConfig.FeeCurrencyGas)
	if msg.CrosschainFee.LT(minimumBurn) {
		return []byte{}, errors.New("crosschain fee amount in message less than minimum burn")
	}

	if k.IsCrossChainFeeReceiverAccountSet(ctx) {
		coins = sdk.NewCoins(sdk.NewCoin(crossChainFeeConfig.FeeCurrency, msg.CrosschainFee))

		err := k.bankKeeper.SendCoins(ctx, cosmosSender, k.GetCrossChainFeeReceiverAccount(ctx), coins)
		if err != nil {
			logger.Error("failed to send crosschain fee from account to account.",
				errorMessageKey, err.Error())
			return []byte{}, err
		}

		coins = sdk.NewCoins(sdk.NewCoin(msg.Symbol, msg.Amount))

	} else {
		if msg.Symbol == crossChainFeeConfig.FeeCurrency {
			coins = sdk.NewCoins(sdk.NewCoin(crossChainFeeConfig.FeeCurrency, msg.CrosschainFee.Add(msg.Amount)))
		} else {
			coins = sdk.NewCoins(sdk.NewCoin(msg.Symbol, msg.Amount), sdk.NewCoin(crossChainFeeConfig.FeeCurrency, msg.CrosschainFee))
		}
	}

	err = k.bankKeeper.SendCoinsFromAccountToModule(ctx, cosmosSender, types.ModuleName, coins)
	if err != nil {
		logger.Error("failed to send crosschain fee from module to account.",
			errorMessageKey, err.Error())
		return []byte{}, err
	}

	coins = sdk.NewCoins(sdk.NewCoin(msg.Symbol, msg.Amount))
	err = k.bankKeeper.BurnCoins(ctx, types.ModuleName, coins)
	if err != nil {
		logger.Error("failed to burn locked coin.",
			errorMessageKey, err.Error())
		return []byte{}, err
	}

	// TODO global sequence will be implemented in other feature
	glocalSequence := uint64(0)

	prophecyID := msg.GetProphecyID(false, senderSequence, glocalSequence)
	k.oracleKeeper.SetProphecyWithInitValue(ctx, prophecyID)

	return prophecyID, nil
}

// ProcessLock processes the lockup of cosmos coins from the given sender
func (k Keeper) ProcessLock(ctx sdk.Context, cosmosSender sdk.AccAddress, senderSequence uint64, msg *types.MsgLock) ([]byte, error) {
	logger := k.Logger(ctx)
	var coins sdk.Coins
	networkIdentity := oracletypes.NewNetworkIdentity(msg.NetworkDescriptor)
	crossChainFeeConfig, err := k.oracleKeeper.GetCrossChainFeeConfig(ctx, networkIdentity)

	if err != nil {
		return []byte{}, err
	}

	minimumLock := crossChainFeeConfig.MinimumLockCost.Mul(crossChainFeeConfig.FeeCurrencyGas)
	if msg.CrosschainFee.LT(minimumLock) {
		return []byte{}, errors.New("crosschain fee amount in message less than minimum lock")
	}

	if k.IsCrossChainFeeReceiverAccountSet(ctx) {
		coins = sdk.NewCoins(sdk.NewCoin(crossChainFeeConfig.FeeCurrency, msg.CrosschainFee))

		err := k.bankKeeper.SendCoins(ctx, cosmosSender, k.GetCrossChainFeeReceiverAccount(ctx), coins)
		if err != nil {
			logger.Error("failed to send crosschain fee from account to account.",
				errorMessageKey, err.Error())
			return []byte{}, err
		}

		coins = sdk.NewCoins(sdk.NewCoin(msg.Symbol, msg.Amount))

	} else {
		coins = sdk.NewCoins(sdk.NewCoin(msg.Symbol, msg.Amount), sdk.NewCoin(crossChainFeeConfig.FeeCurrency, msg.CrosschainFee))
	}

	err = k.bankKeeper.SendCoinsFromAccountToModule(ctx, cosmosSender, types.ModuleName, coins)

	if err != nil {
		logger.Error("failed to transfer coin from account to module.",
			errorMessageKey, err.Error())
		return []byte{}, err
	}

	coins = sdk.NewCoins(sdk.NewCoin(msg.Symbol, msg.Amount))
	err = k.bankKeeper.BurnCoins(ctx, types.ModuleName, coins)
	if err != nil {
		logger.Error("failed to burn burned coin.",
			errorMessageKey, err.Error())
		return []byte{}, err
	}

	// global sequence will be implemented in other feature
	glocalSequence := uint64(0)

	prophecyID := msg.GetProphecyID(false, senderSequence, glocalSequence)
	k.oracleKeeper.SetProphecyWithInitValue(ctx, prophecyID)

	return prophecyID, nil
}

// ProcessUpdateWhiteListValidator processes the update whitelist validator from admin
func (k Keeper) ProcessUpdateWhiteListValidator(ctx sdk.Context, networkDescriptor oracletypes.NetworkDescriptor, cosmosSender sdk.AccAddress, validator sdk.ValAddress, power uint32) error {
	return k.oracleKeeper.ProcessUpdateWhiteListValidator(ctx, networkDescriptor, cosmosSender, validator, power)
}

// ProcessUpdateCrossChainFeeReceiverAccount processes the update crosschain fee receiver account from admin
func (k Keeper) ProcessUpdateCrossChainFeeReceiverAccount(ctx sdk.Context, cosmosSender sdk.AccAddress, crosschainFeeReceiverAccount sdk.AccAddress) error {
	logger := k.Logger(ctx)
	if !k.oracleKeeper.IsAdminAccount(ctx, cosmosSender) {
		logger.Error("cosmos sender is not admin account.")
		return errors.New("only admin account can update CrossChainFee receiver account")
	}

	k.SetCrossChainFeeReceiverAccount(ctx, crosschainFeeReceiverAccount)
	return nil
}

// RescueCrossChainFees transfer CrossChainFee from ethbridge module to an account
func (k Keeper) RescueCrossChainFees(ctx sdk.Context, msg *types.MsgRescueCrossChainFee) error {
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
		return errors.New("only admin account can call rescue CrossChainFee")
	}

	coins := sdk.NewCoins(sdk.NewCoin(msg.CrosschainFeeSymbol, msg.CrosschainFee))
	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, cosmosReceiver, coins)

	if err != nil {
		logger.Error("failed to transfer coin from module to account.",
			errorMessageKey, err.Error())
		return err
	}
	return nil
}

// SetFeeInfo processes the set crosschain fee from admin
func (k Keeper) SetFeeInfo(ctx sdk.Context, msg *types.MsgSetFeeInfo) error {
	logger := k.Logger(ctx)

	cosmosSender, err := sdk.AccAddressFromBech32(msg.CosmosSender)
	if err != nil {
		return err
	}

	if !k.oracleKeeper.IsAdminAccount(ctx, cosmosSender) {
		logger.Error("cosmos sender is not admin account.")
		return errors.New("only admin account can set crosschain fee")
	}
	return k.oracleKeeper.SetFeeInfo(ctx, msg.NetworkDescriptor, msg.FeeCurrency, msg.FeeCurrencyGas, msg.MinimumBurnCost, msg.MinimumLockCost)
}

// ProcessSignProphecy processes the set sign prophecy from validator
func (k Keeper) ProcessSignProphecy(ctx sdk.Context, msg *types.MsgSignProphecy) (oracletypes.StatusText, error) {
	return k.oracleKeeper.ProcessSignProphecy(ctx, msg.NetworkDescriptor, msg.ProphecyId, msg.CosmosSender, msg.EthereumAddress, msg.Signature)
}

// Exists chec if the key existed in db.
func (k Keeper) Exists(ctx sdk.Context, key []byte) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(key)
}
