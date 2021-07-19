package keeper

import (
	"errors"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	gethCommon "github.com/ethereum/go-ethereum/common"
	gethCrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/Sifchain/sifnode/x/oracle/types"
)

// Keeper maintains the link to data storage and
// exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
	cdc      codec.BinaryMarshaler // The wire codec for binary encoding/decoding.
	storeKey sdk.StoreKey          // Unexposed key to access store from sdk.Context

	stakeKeeper types.StakingKeeper
	// TODO: use this as param instead
	consensusNeeded float64 // The minimum % of stake needed to sign claims in order for consensus to occur
}

// NewKeeper creates new instances of the oracle Keeper
func NewKeeper(
	cdc codec.BinaryMarshaler, storeKey sdk.StoreKey, stakeKeeper types.StakingKeeper, consensusNeeded float64,
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

// GetCdc return keeper's cdc
func (k Keeper) GetCdc() codec.BinaryMarshaler {
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
		logger.Error("sifnode oracle keeper ProcessClaim validator not white list.")
		return types.StatusText_STATUS_TEXT_UNSPECIFIED, errors.New("validator not in white list")
	}

	activeValidator := k.checkActiveValidator(ctx, valAddr)
	if !activeValidator {
		logger.Error("sifnode oracle keeper ProcessClaim validator not active.")
		return types.StatusText_STATUS_TEXT_UNSPECIFIED, types.ErrInvalidValidator
	}

	if len(prophecyID) == 0 {
		logger.Error("sifnode oracle keeper ProcessClaim wrong claim id.", "claimID", prophecyID)
		return types.StatusText_STATUS_TEXT_UNSPECIFIED, types.ErrInvalidIdentifier
	}

	prophecy, ok := k.GetProphecy(ctx, prophecyID)
	if !ok {
		prophecy.Id = prophecyID
		prophecy.Status = types.StatusText_STATUS_TEXT_PENDING
	}

	switch prophecy.Status {
	case types.StatusText_STATUS_TEXT_PENDING:

		err = prophecy.AddClaim(valAddr)
		if err != nil {
			return types.StatusText_STATUS_TEXT_UNSPECIFIED, err

		}

		prophecy = k.processCompletion(ctx, networkDescriptor, prophecy)
		k.SetProphecy(ctx, prophecy)
		return prophecy.Status, nil

	case types.StatusText_STATUS_TEXT_SUCCESS:

		err = prophecy.AddClaim(valAddr)
		if err != nil {
			return types.StatusText_STATUS_TEXT_UNSPECIFIED, err
		}
		k.SetProphecy(ctx, prophecy)
		return prophecy.Status, types.ErrProphecyFinalized

	default:
		return types.StatusText_STATUS_TEXT_UNSPECIFIED, types.ErrInvalidProphecyStatus
	}

}

func (k Keeper) checkActiveValidator(ctx sdk.Context, validatorAddress sdk.ValAddress) bool {
	validator, found := k.stakeKeeper.GetValidator(ctx, validatorAddress)
	if !found {
		return false
	}

	return validator.IsBonded()
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

	if voteRate >= k.consensusNeeded {
		prophecy.Status = types.StatusText_STATUS_TEXT_SUCCESS
	}
	return prophecy
}

// SetFeeInfo set crosschain fee for a network
func (k Keeper) SetFeeInfo(ctx sdk.Context, networkDescriptor types.NetworkDescriptor, crossChainFee string, gas, lockCost, burnCost sdk.Int) error {
	k.SetCrossChainFee(ctx, types.NewNetworkIdentity(networkDescriptor), crossChainFee, gas, lockCost, burnCost)
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
func (k Keeper) ProcessSignProphecy(ctx sdk.Context, cosmosSender, prophecyID, ethereumAddress, signature string) error {
	prophecy, ok := k.GetProphecy(ctx, []byte(prophecyID))
	if !ok {
		return fmt.Errorf("failed to get prophecy with ID as %s", prophecyID)
	}

	// verify the signature
	publicKey, err := gethCrypto.Ecrecover([]byte(prophecyID), gethCommon.FromHex(signature))
	if err != nil {
		return err
	}

	ok = gethCrypto.VerifySignature(publicKey, []byte(prophecyID), []byte(signature))
	if !ok {
		return errors.New("incorrect signature")
	}

	prophecy.ClaimValidators = append(prophecy.ClaimValidators, cosmosSender)
	return nil
}

// Exists check if the key exists
func (k Keeper) Exists(ctx sdk.Context, key []byte) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(key)
}
