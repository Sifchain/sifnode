package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
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

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// GetProphecy gets the entire prophecy data struct for a given id
func (k Keeper) GetProphecy(ctx sdk.Context, id string) (types.Prophecy, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get([]byte(id))
	if bz == nil {
		return types.Prophecy{}, false
	}

	var dbProphecy types.DBProphecy
	k.cdc.MustUnmarshalBinaryBare(bz, &dbProphecy)

	deSerializedProphecy, err := dbProphecy.DeserializeFromDB()
	if err != nil {
		return types.Prophecy{}, false
	}

	return deSerializedProphecy, true
}

// setProphecy saves a prophecy with an initial claim
func (k Keeper) setProphecy(ctx sdk.Context, prophecy types.Prophecy) {
	store := ctx.KVStore(k.storeKey)
	serializedProphecy, err := prophecy.SerializeForDB()
	if err != nil {
		panic(err)
	}

	store.Set([]byte(prophecy.ID), k.cdc.MustMarshalBinaryBare(&serializedProphecy))
}

// ProcessClaim ...
func (k Keeper) ProcessClaim(ctx sdk.Context, claim types.Claim) (types.Status, error) {
	logger := k.Logger(ctx)
	inWhiteList := false
	// Check if claim from whitelist validators
	whiteList := k.GetOracleWhiteList(ctx)
	logger.Info("oracle whitelist is", "whitelist", whiteList)
	for _, address := range whiteList {

		if address.String() == (claim.ValidatorAddress) {
			inWhiteList = true
			break
		}
	}

	if !inWhiteList {
		logger.Error("sifnode oracle keeper ProcessClaim validator no in whitelist.")
		return types.Status{}, types.ErrValidatorNotInWhiteList
	}

	valAddr, err := sdk.ValAddressFromBech32(claim.ValidatorAddress)
	if err != nil {
		return types.Status{}, err
	}

	activeValidator := k.checkActiveValidator(ctx, valAddr)
	if !activeValidator {
		logger.Error("sifnode oracle keeper ProcessClaim validator not active.")
		return types.Status{}, types.ErrInvalidValidator
	}

	if claim.Id == "" {
		logger.Error("sifnode oracle keeper ProcessClaim wrong claim id.", "claimID", claim.Id)
		return types.Status{}, types.ErrInvalidIdentifier
	}

	if claim.Content == "" {
		logger.Error("sifnode oracle keeper ProcessClaim claim content is empty.")
		return types.Status{}, types.ErrInvalidClaim
	}

	prophecy, found := k.GetProphecy(ctx, claim.Id)
	if !found {
		prophecy = types.NewProphecy(claim.Id)
	}
	switch prophecy.Status.Text {
	case types.StatusText_STATUS_TEXT_PENDING:
		// continue processing
	default:
		return types.Status{}, types.ErrProphecyFinalized
	}

	if prophecy.ValidatorClaims[claim.ValidatorAddress] != "" {
		return types.Status{}, types.ErrDuplicateMessage
	}

	prophecy.AddClaim(valAddr, claim.Content)
	prophecy = k.processCompletion(ctx, prophecy)

	k.setProphecy(ctx, prophecy)
	return prophecy.Status, nil
}

func (k Keeper) checkActiveValidator(ctx sdk.Context, validatorAddress sdk.ValAddress) bool {
	validator, found := k.stakeKeeper.GetValidator(ctx, validatorAddress)
	if !found {
		return false
	}

	return validator.IsBonded()
}

// ProcessUpdateWhiteListValidator processes the update whitelist validator from admin
func (k Keeper) ProcessUpdateWhiteListValidator(ctx sdk.Context, cosmosSender sdk.AccAddress, validator sdk.ValAddress, operationtype string) error {
	logger := k.Logger(ctx)
	if !k.IsAdminAccount(ctx, cosmosSender) {
		logger.Error("cosmos sender is not admin account.")
		return types.ErrNotAdminAccount
	}

	switch operationtype {
	case "add":
		k.AddOracleWhiteList(ctx, validator)
	case "remove":
		k.RemoveOracleWhiteList(ctx, validator)
	default:
		return types.ErrInvalidOperationType
	}

	return nil
}

// processCompletion looks at a given prophecy
// an assesses whether the claim with the highest power on that prophecy has enough
// power to be considered successful, or alternatively,
// will never be able to become successful due to not enough validation power being
// left to push it over the threshold required for consensus.
func (k Keeper) processCompletion(ctx sdk.Context, prophecy types.Prophecy) types.Prophecy {
	highestClaim, highestClaimPower, totalClaimsPower, totalPower := prophecy.FindHighestClaim(ctx, k.stakeKeeper, k.GetOracleWhiteList(ctx))

	highestConsensusRatio := float64(highestClaimPower) / float64(totalPower)
	remainingPossibleClaimPower := totalPower - totalClaimsPower
	highestPossibleClaimPower := highestClaimPower + remainingPossibleClaimPower
	highestPossibleConsensusRatio := float64(highestPossibleClaimPower) / float64(totalPower)

	if highestConsensusRatio >= k.consensusNeeded {
		prophecy.Status.Text = types.StatusText_STATUS_TEXT_SUCCESS
		prophecy.Status.FinalClaim = highestClaim
	} else if highestPossibleConsensusRatio < k.consensusNeeded {
		prophecy.Status.Text = types.StatusText_STATUS_TEXT_FAILED
	}

	return prophecy
}

// Exists check if the key exists
func (k Keeper) Exists(ctx sdk.Context, key []byte) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(key)
}
