package keeper

import (
	"errors"
	"fmt"
	"go.uber.org/zap"

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
	cdc codec.BinaryMarshaler // The wire codec for binary encoding/decoding.

	accountKeeper       types.AccountKeeper
	bankKeeper          types.BankKeeper
	oracleKeeper        types.OracleKeeper
	tokenRegistryKeeper tokenregistrytypes.Keeper
	storeKey            sdk.StoreKey
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
func NewKeeper(cdc codec.BinaryMarshaler,
	bankKeeper types.BankKeeper,
	oracleKeeper types.OracleKeeper,
	accountKeeper types.AccountKeeper,
	tokenRegistryKeeper tokenregistrytypes.Keeper,
	storeKey sdk.StoreKey) Keeper {
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
func (k Keeper) ProcessClaim(ctx sdk.Context, claim *types.EthBridgeClaim) (oracletypes.StatusText, error) {
	return k.oracleKeeper.ProcessClaim(ctx, claim.NetworkDescriptor, claim.GetProphecyID(), claim.ValidatorAddress)
}

// ProcessSuccessfulClaim processes a claim that has just completed successfully with consensus
func (k Keeper) ProcessSuccessfulClaim(ctx sdk.Context, claim *types.EthBridgeClaim) error {
	logger := k.Logger(ctx)

	logger.Debug(types.PeggyTestMarker, "kind", "ProcessSuccessfulClaim", "claim", zap.Reflect("claim", claim))

	tokenMetadata, ok := k.tokenRegistryKeeper.GetTokenMetadata(ctx, claim.DenomHash)
	if !ok {
		return fmt.Errorf("token metadata not available for %s", claim.DenomHash)
	}

	var coins sdk.Coins
	var err error
	switch claim.ClaimType {
	case types.ClaimType_CLAIM_TYPE_LOCK:
		coins = sdk.NewCoins(sdk.NewCoin(tokenMetadata.Symbol, claim.Amount))
		err = k.bankKeeper.MintCoins(ctx, types.ModuleName, coins)
	case types.ClaimType_CLAIM_TYPE_BURN:
		coins = sdk.NewCoins(sdk.NewCoin(tokenMetadata.Symbol, claim.Amount))
		err = nil
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

	ctx.Logger().Debug(types.PeggyTestMarker, "kind", "coinsSent", "claim", zap.Reflect("claim", claim), "receiverAddress", receiverAddress, "coins", coins)

	return nil
}

// ProcessBurn processes the burn of bridged coins from the given sender
func (k Keeper) ProcessBurn(ctx sdk.Context,
	cosmosSender sdk.AccAddress,
	senderSequence uint64,
	msg *types.MsgBurn,
	tokenMetadata tokenregistrytypes.TokenMetadata) ([]byte, error) {

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

	if tokenMetadata.NetworkDescriptor.IsSifchain() {
		logger.Error("sifchain natvie token can't be burn.", "tokenSymbol", tokenMetadata.Symbol)
		return []byte{}, fmt.Errorf("sifchain native token %s can't be burn", tokenMetadata.Symbol)
	}

	if k.IsCrossChainFeeReceiverAccountSet(ctx) {
		coins = sdk.NewCoins(sdk.NewCoin(crossChainFeeConfig.FeeCurrency, msg.CrosschainFee))

		err := k.bankKeeper.SendCoins(ctx, cosmosSender, k.GetCrossChainFeeReceiverAccount(ctx), coins)
		if err != nil {
			logger.Error("failed to send crosschain fee from account to account.",
				errorMessageKey, err.Error())
			return []byte{}, err
		}

		coins = sdk.NewCoins(sdk.NewCoin(tokenMetadata.Symbol, msg.Amount))

	} else {
		if msg.DenomHash == crossChainFeeConfig.FeeCurrency {
			coins = sdk.NewCoins(sdk.NewCoin(crossChainFeeConfig.FeeCurrency, msg.CrosschainFee.Add(msg.Amount)))
		} else {
			coins = sdk.NewCoins(sdk.NewCoin(tokenMetadata.Symbol, msg.Amount), sdk.NewCoin(crossChainFeeConfig.FeeCurrency, msg.CrosschainFee))
		}
	}

	err = k.bankKeeper.SendCoinsFromAccountToModule(ctx, cosmosSender, types.ModuleName, coins)
	if err != nil {
		logger.Error("failed to send crosschain fee from module to account.",
			errorMessageKey, err.Error())
		return []byte{}, err
	}

	// not burn the token if it is sifchain native token
	if !tokenMetadata.NetworkDescriptor.IsSifchain() {
		coins = sdk.NewCoins(sdk.NewCoin(tokenMetadata.Symbol, msg.Amount))
		err = k.bankKeeper.BurnCoins(ctx, types.ModuleName, coins)
		if err != nil {
			logger.Error("failed to burn locked coin.",
				errorMessageKey, err.Error())
			return []byte{}, err
		}
	}

	prophecyID := msg.GetProphecyID(false, senderSequence, k.GetGlobalNonce(ctx, msg.NetworkDescriptor), tokenMetadata.TokenAddress)
	k.oracleKeeper.SetProphecyWithInitValue(ctx, prophecyID)

	return prophecyID, nil
}

// ProcessLock processes the lockup of cosmos coins from the given sender
func (k Keeper) ProcessLock(ctx sdk.Context,
	cosmosSender sdk.AccAddress,
	senderSequence uint64,
	msg *types.MsgLock,
	tokenMetadata tokenregistrytypes.TokenMetadata,
	firstDoublePeg bool) ([]byte, error) {

	logger := k.Logger(ctx)
	var coins sdk.Coins
	networkIdentity := oracletypes.NewNetworkIdentity(msg.NetworkDescriptor)
	crossChainFeeConfig, err := k.oracleKeeper.GetCrossChainFeeConfig(ctx, networkIdentity)

	if err != nil {
		return []byte{}, err
	}

	if !tokenMetadata.NetworkDescriptor.IsSifchain() {
		logger.Error("pegged token can't be lock.", "tokenSymbol", tokenMetadata.Symbol)
		return []byte{}, fmt.Errorf("pegged token %s can't be lock", tokenMetadata.Symbol)
	}

	// check if it is the first time to do double peg
	cost := crossChainFeeConfig.MinimumLockCost
	if firstDoublePeg {
		cost = cost.Add(crossChainFeeConfig.FirstLockDoublePeggyCost)
	}

	minimumLock := cost.Mul(crossChainFeeConfig.FeeCurrencyGas)
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

		coins = sdk.NewCoins(sdk.NewCoin(tokenMetadata.Symbol, msg.Amount))

	} else {
		coins = sdk.NewCoins(sdk.NewCoin(tokenMetadata.Symbol, msg.Amount), sdk.NewCoin(crossChainFeeConfig.FeeCurrency, msg.CrosschainFee))
	}

	err = k.bankKeeper.SendCoinsFromAccountToModule(ctx, cosmosSender, types.ModuleName, coins)

	if err != nil {
		logger.Error("failed to transfer coin from account to module.",
			errorMessageKey, err.Error())
		return []byte{}, err
	}

	coins = sdk.NewCoins(sdk.NewCoin(tokenMetadata.Symbol, msg.Amount))
	err = k.bankKeeper.BurnCoins(ctx, types.ModuleName, coins)
	if err != nil {
		logger.Error("failed to burn burned coin.",
			errorMessageKey, err.Error())
		return []byte{}, err
	}

	prophecyID := msg.GetProphecyID(false, senderSequence, k.GetGlobalNonce(ctx, msg.NetworkDescriptor), tokenMetadata.TokenAddress)
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
	return k.oracleKeeper.SetFeeInfo(ctx,
		msg.NetworkDescriptor,
		msg.FeeCurrency,
		msg.FeeCurrencyGas,
		msg.MinimumBurnCost,
		msg.MinimumLockCost,
		msg.FirstLockDoublePeggyCost)
}

// ProcessSignProphecy processes the set sign prophecy from validator
func (k Keeper) ProcessSignProphecy(ctx sdk.Context, msg *types.MsgSignProphecy) error {
	prophecyInfo, ok := k.oracleKeeper.GetProphecyInfo(ctx, msg.ProphecyId)
	if !ok {
		return errors.New("prophecy not found in oracle keeper")
	}

	metadata, ok := k.GetTokenMetadata(ctx, prophecyInfo.TokenDenomHash)
	if !ok {
		return fmt.Errorf("metadata not available for %s", prophecyInfo.TokenDenomHash)
	}

	return k.oracleKeeper.ProcessSignProphecy(ctx, msg.NetworkDescriptor, msg.ProphecyId, msg.CosmosSender, metadata.TokenAddress, msg.EthereumAddress, msg.Signature)
}

// GetTokenRegistryKeeper return token registry keeper
func (k Keeper) GetTokenRegistryKeeper() tokenregistrytypes.Keeper {
	return k.tokenRegistryKeeper
}

// GetTokenMetadata call metadataKeeper's GetTokenMetadata
func (k Keeper) GetTokenMetadata(ctx sdk.Context, denomHash string) (tokenregistrytypes.TokenMetadata, bool) {
	return k.tokenRegistryKeeper.GetTokenMetadata(ctx, denomHash)
}

// AddTokenMetadata call metadataKeeper's AddTokenMetadata
func (k Keeper) AddTokenMetadata(ctx sdk.Context, metadata tokenregistrytypes.TokenMetadata) string {
	return k.tokenRegistryKeeper.AddTokenMetadata(ctx, metadata)
}

// GetWitnessLockBurnNonce return witness lock burn nonce from oracle keeper
func (k Keeper) GetWitnessLockBurnNonce(ctx sdk.Context, networkDescriptor oracletypes.NetworkDescriptor, valAccount sdk.ValAddress) uint64 {
	return k.oracleKeeper.GetWitnessLockBurnNonce(ctx, networkDescriptor, valAccount)
}

// Exists chec if the key existed in db.
func (k Keeper) Exists(ctx sdk.Context, key []byte) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(key)
}
