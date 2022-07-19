package ante

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	disttypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

var minCommission = sdk.NewDecWithPrec(5, 2)   // 5%
var maxVotingPower = sdk.NewDecWithPrec(10, 2) // 10%

// TODO: remove once Cosmos SDK is upgraded to v0.46, refer to https://github.com/cosmos/cosmos-sdk/pull/10529#issuecomment-1026320612

// ValidateMinCommissionDecorator validates that the validator commission is always
// greater than or equal to the min commission rate
type ValidateMinCommissionDecorator struct {
	sk         stakingkeeper.Keeper
	bankKeeper stakingtypes.BankKeeper
}

// ValidateMinCommissionDecorator creates a new ValidateMinCommissionDecorator
func NewValidateMinCommissionDecorator(sk stakingkeeper.Keeper, bk stakingtypes.BankKeeper) ValidateMinCommissionDecorator {
	return ValidateMinCommissionDecorator{
		sk:         sk,
		bankKeeper: bk,
	}
}

func (vcd ValidateMinCommissionDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	for _, msg := range tx.GetMsgs() {
		if err := vcd.validateMsg(ctx, msg); err != nil {
			return ctx, err
		}
	}

	return next(ctx, tx, simulate)
}

// getValidator returns the validator belonging to a given bech32 validator address
func (vcd ValidateMinCommissionDecorator) getValidator(ctx sdk.Context, bech32ValAddr string) (stakingtypes.ValidatorI, error) {
	valAddr, err := sdk.ValAddressFromBech32(bech32ValAddr)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, bech32ValAddr)
	}

	val := vcd.sk.Validator(ctx, valAddr)
	if val == nil {
		return nil, disttypes.ErrNoValidatorExists
	}

	return val, nil
}

func (vcd ValidateMinCommissionDecorator) validateMsg(ctx sdk.Context, msg sdk.Msg) error {
	switch msg := msg.(type) {
	case *stakingtypes.MsgCreateValidator:
		if msg.Commission.Rate.LT(minCommission) {
			return sdkerrors.Wrapf(
				sdkerrors.ErrInvalidRequest,
				"validator commission %s cannot be lower than minimum of %s", msg.Commission.Rate, minCommission)
		}
	case *stakingtypes.MsgEditValidator:
		if msg.CommissionRate != nil && msg.CommissionRate.LT(minCommission) {
			return sdkerrors.Wrapf(
				sdkerrors.ErrInvalidRequest,
				"validator commission %s cannot be lower than minimum of %s", msg.CommissionRate, minCommission)
		}
	case *stakingtypes.MsgDelegate:
		val, err := vcd.getValidator(ctx, msg.ValidatorAddress)
		if err != nil {
			return err
		}

		votingPower := vcd.calculateProvisionalVotingPower(ctx, val, sdk.NewDecFromInt(msg.Amount.Amount))
		if votingPower.GTE(maxVotingPower) {
			return sdkerrors.Wrapf(
				sdkerrors.ErrInvalidRequest,
				"validator has %s voting power, cannot delegate to a validator with %s or higher voting power, please choose another validator", votingPower, maxVotingPower)
		}
	case *stakingtypes.MsgBeginRedelegate:
		val, err := vcd.getValidator(ctx, msg.ValidatorDstAddress)
		if err != nil {
			return err
		}

		votingPower := vcd.calculateProvisionalVotingPower(ctx, val, sdk.NewDecFromInt(msg.Amount.Amount))
		if votingPower.GTE(maxVotingPower) {
			return sdkerrors.Wrapf(
				sdkerrors.ErrInvalidRequest,
				"validator has %s voting power, cannot redelegate to a validator with %s or higher voting power, please choose another validator", votingPower, maxVotingPower)
		}
	}
	return nil
}

func (vcd ValidateMinCommissionDecorator) calculateProvisionalVotingPower(ctx sdk.Context, validator stakingtypes.ValidatorI, delegateAmount sdk.Dec) sdk.Dec {
	// TODO: watch for divide by zero?
	validatorTokens := sdk.NewDecFromInt(validator.GetTokens())
	provisionalValidatorTokens := validatorTokens.Add(delegateAmount)
	totalBondedTokens := sdk.NewDecFromInt(vcd.sk.TotalBondedTokens(ctx))

	//return validatorTokens.Quo(totalBondedTokens)

	provisionalTotalBondedTokens := sdk.ZeroDec()
	if validator.IsBonded() {
		provisionalTotalBondedTokens = totalBondedTokens.Add(delegateAmount)
	} else {
		bondedValidators := vcd.sk.GetBondedValidatorsByPower(ctx)
		weakestValidatorTokens := sdk.Dec(bondedValidators[len(bondedValidators)-1].Tokens)

		if weakestValidatorTokens.LT(provisionalValidatorTokens) {
			//validator will still not be bonded so will have no voting power
			return sdk.ZeroDec()
		}

		// validator will become bonded
		provisionalTotalBondedTokens = totalBondedTokens.Add(delegateAmount).Sub(weakestValidatorTokens)
	}

	return provisionalValidatorTokens.Quo(provisionalTotalBondedTokens)
}

func (vcd ValidateMinCommissionDecorator) getDelegatedTokens(ctx sdk.Context) sdk.Int {

	bondDenom := vcd.sk.BondDenom(ctx)
	bondedPool := vcd.sk.GetBondedPool(ctx)
	notBondedPool := vcd.sk.GetNotBondedPool(ctx)

	bondedTokens := vcd.bankKeeper.GetBalance(ctx, notBondedPool.GetAddress(), bondDenom).Amount
	notBondedTokens := vcd.bankKeeper.GetBalance(ctx, bondedPool.GetAddress(), bondDenom).Amount

	return bondedTokens.Add(notBondedTokens)

}
