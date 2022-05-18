package keeper

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/Sifchain/sifnode/x/instrumentation"

	"go.uber.org/zap"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/crypto"
	gethCrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/Sifchain/sifnode/x/oracle/types"
	solsha3 "github.com/miguelmota/go-solidity-sha3"
)

// Keeper maintains the link to data storage and
// exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
	cdc             codec.BinaryCodec // The wire codec for binary encoding/decoding.
	storeKey        sdk.StoreKey      // Unexposed key to access store from sdk.Context
	stakeKeeper     types.StakingKeeper
	consensusNeeded float64 // The minimum % of stake needed to sign claims in order for consensus to occur
	currentHeight   int64
}

// NewKeeper creates new instances of the oracle Keeper
func NewKeeper(
	cdc codec.BinaryCodec, storeKey sdk.StoreKey, stakeKeeper types.StakingKeeper, consensusNeeded float64,
) Keeper {
	if consensusNeeded <= 0 || consensusNeeded > 1 {
		panic(types.ErrMinimumConsensusNeededInvalid.Error())
	}
	return Keeper{
		cdc:             cdc,
		storeKey:        storeKey,
		stakeKeeper:     stakeKeeper,
		consensusNeeded: consensusNeeded,
	}
}

func (k *Keeper) UpdateCurrentHeight(height int64) {
	k.currentHeight = height
}

// GetCdc return keeper's cdc
func (k Keeper) GetCdc() codec.BinaryCodec {
	return k.cdc
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// ProcessClaim handle claim
func (k Keeper) ProcessClaim(ctx sdk.Context, networkDescriptor types.NetworkDescriptor, prophecyID []byte, validator string) (types.StatusText, error) {
	logger := k.Logger(ctx)
	networkIdentity := types.NewNetworkIdentity(networkDescriptor)

	valAddr, err := sdk.ValAddressFromBech32(validator)
	if err != nil {
		return types.StatusText_STATUS_TEXT_UNSPECIFIED, err
	}

	if !k.ValidateAddress(ctx, networkIdentity, valAddr) {
		logger.Error("ProcessClaim: ValidateAddress returned false", "networkDescriptor", networkIdentity, "valAddr", valAddr)
		return types.StatusText_STATUS_TEXT_UNSPECIFIED, errors.New("validator not in white list")
	}

	if len(prophecyID) == 0 {
		logger.Error("sifnode oracle keeper ProcessClaim wrong claim id.", "claimID", prophecyID)
		return types.StatusText_STATUS_TEXT_UNSPECIFIED, types.ErrInvalidIdentifier
	}

	return k.AppendValidatorToProphecy(ctx, networkDescriptor, prophecyID, valAddr)
}

// AppendValidatorToProphecy append the validator's signature to prophecy
func (k Keeper) AppendValidatorToProphecy(ctx sdk.Context, networkDescriptor types.NetworkDescriptor, prophecyID []byte, validator sdk.ValAddress) (types.StatusText, error) {
	prophecy, ok := k.GetProphecy(ctx, prophecyID)
	if !ok {
		prophecy.Id = prophecyID
		prophecy.Status = types.StatusText_STATUS_TEXT_PENDING
	}

	instrumentation.PeggyCheckpoint(ctx.Logger(), instrumentation.AppendValidatorToProphecy, "prophecy", zap.Reflect("prophecy", prophecy))

	switch prophecy.Status {
	case types.StatusText_STATUS_TEXT_PENDING:

		err := prophecy.AddClaim(validator)
		if err != nil {
			return types.StatusText_STATUS_TEXT_UNSPECIFIED, err

		}

		prophecy = k.processCompletion(ctx, networkDescriptor, prophecy)
		k.SetProphecy(ctx, prophecy)

		return prophecy.Status, nil

	case types.StatusText_STATUS_TEXT_SUCCESS:

		err := prophecy.AddClaim(validator)
		if err != nil {
			return types.StatusText_STATUS_TEXT_UNSPECIFIED, err
		}
		k.SetProphecy(ctx, prophecy)
		return prophecy.Status, types.ErrProphecyFinalized

	default:
		return types.StatusText_STATUS_TEXT_UNSPECIFIED, types.ErrInvalidProphecyStatus
	}
}

// ProcessUpdateWhiteListValidator processes the update whitelist validator from admin
func (k Keeper) ProcessUpdateWhiteListValidator(ctx sdk.Context, networkDescriptor types.NetworkDescriptor, cosmosSender sdk.AccAddress, validator sdk.ValAddress, power uint32) error {
	logger := k.Logger(ctx)
	if !k.IsAdminAccount(ctx, cosmosSender) {
		logger.Error("cosmos sender is not admin account.")
		return types.ErrNotAdminAccount
	}

	k.UpdateOracleWhiteList(ctx, types.NewNetworkIdentity(networkDescriptor), validator, power)
	return nil
}

// processCompletion looks at a given prophecy
func (k Keeper) processCompletion(ctx sdk.Context, networkDescriptor types.NetworkDescriptor, prophecy types.Prophecy) types.Prophecy {
	whiteList := k.GetOracleWhiteList(ctx, types.NewNetworkIdentity(networkDescriptor))
	voteRate := whiteList.GetPowerRatio(prophecy.ClaimValidators)

	var consensusNeeded float64
	consensusNeededUint, err := k.GetConsensusNeeded(ctx, types.NewNetworkIdentity(networkDescriptor))
	// consensusNeeded unavailable from keeper, use the default one.
	if err != nil {
		consensusNeeded = k.consensusNeeded
	} else {
		consensusNeeded = float64(consensusNeededUint) / 100.0
	}

	if voteRate >= consensusNeeded {
		prophecy.Status = types.StatusText_STATUS_TEXT_SUCCESS
	}

	instrumentation.PeggyCheckpoint(ctx.Logger(), instrumentation.ProcessCompletion,
		"prophecy", zap.Reflect("prophecy", prophecy),
		"whitelist", zap.Reflect("whiteList", whiteList),
		"consensusNeededUint", consensusNeededUint,
		"voteRate", voteRate,
	)

	return prophecy
}

// SetFeeInfo set crosschain fee for a network
func (k Keeper) SetFeeInfo(ctx sdk.Context, networkDescriptor types.NetworkDescriptor, crossChainFee string, gas, lockCost, burnCost, firstLockDoublePeggyCost sdk.Int) error {
	k.SetCrossChainFee(ctx, types.NewNetworkIdentity(networkDescriptor), crossChainFee, gas, lockCost, burnCost, firstLockDoublePeggyCost)
	return nil
}

// SetProphecyWithInitValue set the prophecy in keeper
func (k Keeper) SetProphecyWithInitValue(ctx sdk.Context, prophecyID []byte) {
	prophecy := types.Prophecy{
		Id:              prophecyID,
		Status:          types.StatusText_STATUS_TEXT_PENDING,
		ClaimValidators: []string{},
	}
	k.SetProphecy(ctx, prophecy)
}

// ProcessSignProphecy deal with the signature from validator
func (k Keeper) ProcessSignProphecy(ctx sdk.Context, networkDescriptor types.NetworkDescriptor, prophecyID []byte, cosmosSender, tokenAddress, ethereumAddress, signature string) error {
	prophecy, ok := k.GetProphecy(ctx, prophecyID)
	if !ok {
		return types.ErrProphecyNotFound
	}

	whiteList := k.GetOracleWhiteList(ctx, types.NewNetworkIdentity(networkDescriptor))
	power, ok := whiteList.WhiteList[cosmosSender]
	if !ok {
		return errors.New("message sender to sign prophecy not in the whitelist")
	}

	if power == 0 {
		return errors.New("message sender to sign prophecy without vote power")
	}

	// verify the signature
	sigData := PrefixMsg(prophecyID)

	publicKey, err := gethCrypto.Ecrecover(sigData, []byte(signature))
	if err != nil {
		return err
	}

	pubKey, err := crypto.UnmarshalPubkey(publicKey)
	if err != nil {
		return err
	}

	recoveredAddr := crypto.PubkeyToAddress(*pubKey)

	if recoveredAddr.String() != ethereumAddress {
		return errors.New("incorrect ethereum signature")
	}

	valAddr, err := sdk.ValAddressFromBech32(cosmosSender)
	if err != nil {
		return err
	}

	prophecyInfo, ok := k.GetProphecyInfo(ctx, prophecyID)
	if !ok {
		return errors.New("prophecy info not available in keeper")
	}

	// check the order of witness lock burn sequence
	lastLockBurnNonce := k.GetWitnessLockBurnSequence(ctx, networkDescriptor, valAddr)
	if lastLockBurnNonce != 0 && lastLockBurnNonce+1 != prophecyInfo.GlobalSequence {
		return errors.New("witness node not send the signature in order")
	}

	oldStatus := prophecy.Status

	newStatus, err := k.AppendValidatorToProphecy(ctx, networkDescriptor, prophecyID, valAddr)
	if err != nil {
		return err
	}

	err = k.AppendSignature(ctx, prophecyID, ethereumAddress, signature)
	if err != nil {
		return err
	}

	// update witness's lock burn sequence
	k.SetWitnessLockBurnNonce(ctx, networkDescriptor, valAddr, prophecyInfo.GlobalSequence)

	// emit the event when status from pending to success
	// old = unspecified, new = pending  the prophecy just created, not emit the event
	// old = unspecified, new = success no such path
	// old = pending, new = success the only case we will emit the event
	// old = success, new = success not emit the same event twice
	if oldStatus == types.StatusText_STATUS_TEXT_PENDING && newStatus == types.StatusText_STATUS_TEXT_SUCCESS {
		event := sdk.NewEvent(
			types.EventTypeProphecyCompleted,
			sdk.NewAttribute(types.AttributeKeyProphecyID, string(prophecyID)),
			sdk.NewAttribute(types.AttributeKeyNetworkDescriptor, prophecyInfo.NetworkDescriptor.String()),
			sdk.NewAttribute(types.AttributeKeyCosmosSender, prophecyInfo.CosmosSender),
			sdk.NewAttribute(types.AttributeKeyCosmosSenderSequence, strconv.FormatInt(int64(prophecyInfo.CosmosSenderSequence), 10)),
			sdk.NewAttribute(types.AttributeKeyEthereumReceiver, prophecyInfo.EthereumReceiver),
			sdk.NewAttribute(types.AttributeKeyTokenContractAddress, tokenAddress),
			sdk.NewAttribute(types.AttributeKeyAmount, strconv.FormatInt(prophecyInfo.TokenAmount.Int64(), 10)),
			sdk.NewAttribute(types.AttributeKeyBridgeToken, strconv.FormatBool(prophecyInfo.BridgeToken)),
			sdk.NewAttribute(types.AttributeKeyGlobalNonce, strconv.FormatInt(int64(prophecyInfo.GlobalSequence), 10)),
			sdk.NewAttribute(types.AttributeKeycrossChainFee, strconv.FormatInt(prophecyInfo.CrosschainFee.Int64(), 10)),
			sdk.NewAttribute(types.AttributeKeySignatures, strings.Join(prophecyInfo.Signatures, ",")),
			sdk.NewAttribute(types.AttributeKeyEthereumAddresses, strings.Join(prophecyInfo.EthereumAddress, ",")),
		)
		ctx.EventManager().EmitEvents(sdk.Events{event})
		instrumentation.PeggyCheckpoint(
			ctx.Logger(),
			instrumentation.ProphecyStatus,
			"event", zap.Reflect("event", event),
			"prophecyInfo", zap.Reflect("prophecyInfo", prophecyInfo),
		)
	}
	return nil
}

// ProcessUpdateConsensusNeeded
func (k Keeper) ProcessUpdateConsensusNeeded(ctx sdk.Context, cosmosSender sdk.AccAddress, networkDescriptor types.NetworkDescriptor, consensusNeeded uint32) error {
	logger := k.Logger(ctx)
	if !k.IsAdminAccount(ctx, cosmosSender) {
		logger.Error("cosmos sender is not admin account.")
		return types.ErrNotAdminAccount
	}

	k.SetConsensusNeeded(ctx, types.NewNetworkIdentity(networkDescriptor), consensusNeeded)
	return nil
}

// Exists check if the key exists
func (k Keeper) Exists(ctx sdk.Context, key []byte) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(key)
}

// PrefixMsg prefixes a message for verification, mimics behavior of web3.eth.sign
func PrefixMsg(msg []byte) []byte {
	return solsha3.SoliditySHA3(solsha3.String("\x19Ethereum Signed Message:\n32"), msg)
}
