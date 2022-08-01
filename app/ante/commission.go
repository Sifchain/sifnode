package ante

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
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
	bankkeeper bankkeeper.Keeper
}

// ValidateMinCommissionDecorator creates a new ValidateMinCommissionDecorator
func NewValidateMinCommissionDecorator(sk stakingkeeper.Keeper, bk bankkeeper.Keeper) ValidateMinCommissionDecorator {
	return ValidateMinCommissionDecorator{
		sk:         sk,
		bankkeeper: bk,
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
func (vcd ValidateMinCommissionDecorator) getValidator(ctx sdk.Context, bech32ValAddr string) (stakingtypes.Validator, error) {
	valAddr, err := sdk.ValAddressFromBech32(bech32ValAddr)
	if err != nil {
		return stakingtypes.Validator{}, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, bech32ValAddr)
	}

	val, found := vcd.sk.GetValidator(ctx, valAddr)
	if !found {
		return stakingtypes.Validator{}, disttypes.ErrNoValidatorExists
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

		projectedVotingPower := vcd.CalculateDelegateProjectedVotingPower(ctx, val, sdk.NewDecFromInt(msg.Amount.Amount))
		if projectedVotingPower.GTE(maxVotingPower) {
			return sdkerrors.Wrapf(
				sdkerrors.ErrInvalidRequest,
				"validator would have %s voting power, cannot delegate to a validator with projected %s or higher voting power", projectedVotingPower, maxVotingPower)
		}
	case *stakingtypes.MsgBeginRedelegate:
		dstVal, err := vcd.getValidator(ctx, msg.ValidatorDstAddress)
		if err != nil {
			return err
		}

		var delegateAmount sdk.Dec
		if msg.ValidatorSrcAddress == msg.ValidatorDstAddress {
			// This is blocked later on by the SDK. However we may as well calculate the correct projected voting power.
			// Since this is a self redelegation, no additional tokens are delegated to this validator hence delegateAmount = 0
			delegateAmount = sdk.ZeroDec()
		} else {
			delegateAmount = sdk.NewDecFromInt(msg.Amount.Amount)
		}

		projectedVotingPower := vcd.CalculateRedelegateProjectedVotingPower(ctx, dstVal, delegateAmount)
		if projectedVotingPower.GTE(maxVotingPower) {
			return sdkerrors.Wrapf(
				sdkerrors.ErrInvalidRequest,
				"validator would have %s voting power, cannot redelegate to a validator with projected %s or higher voting power", projectedVotingPower, maxVotingPower)
		}
	}
	return nil
}

func (vcd ValidateMinCommissionDecorator) getTotalDelegatedTokens(ctx sdk.Context) sdk.Int {
	bondDenom := vcd.sk.BondDenom(ctx)
	bondedPool := vcd.sk.GetBondedPool(ctx)
	notBondedPool := vcd.sk.GetNotBondedPool(ctx)

	notBondedAmount := vcd.bankkeeper.GetBalance(ctx, notBondedPool.GetAddress(), bondDenom).Amount
	bondedAmount := vcd.bankkeeper.GetBalance(ctx, bondedPool.GetAddress(), bondDenom).Amount

	return notBondedAmount.Add(bondedAmount)
}

func (vcd ValidateMinCommissionDecorator) CalculateDelegateProjectedVotingPower(ctx sdk.Context, validator stakingtypes.ValidatorI, delegateAmount sdk.Dec) sdk.Dec {
	validatorTokens := sdk.NewDecFromInt(validator.GetTokens())
	totalDelegatedTokens := sdk.NewDecFromInt(vcd.getTotalDelegatedTokens(ctx))

	projectedTotalDelegatedTokens := totalDelegatedTokens.Add(delegateAmount)
	projectedValidatorTokens := validatorTokens.Add(delegateAmount)

	return projectedValidatorTokens.Quo(projectedTotalDelegatedTokens)
}

func (vcd ValidateMinCommissionDecorator) CalculateRedelegateProjectedVotingPower(ctx sdk.Context, validator stakingtypes.ValidatorI, delegateAmount sdk.Dec) sdk.Dec {
	validatorTokens := sdk.NewDecFromInt(validator.GetTokens())
	projectedTotalDelegatedTokens := sdk.NewDecFromInt(vcd.getTotalDelegatedTokens(ctx)) // no additional delegated tokens

	projectedValidatorTokens := validatorTokens.Add(delegateAmount)

	return projectedValidatorTokens.Quo(projectedTotalDelegatedTokens)
}
