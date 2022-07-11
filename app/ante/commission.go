package ante

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	disttypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

var minCommission = sdk.NewDecWithPrec(5, 2)    // 5%
var maxVotingPower = sdk.NewDecWithPrec(100, 2) // 100%

// TODO: remove once Cosmos SDK is upgraded to v0.46, refer to https://github.com/cosmos/cosmos-sdk/pull/10529#issuecomment-1026320612

// ValidateMinCommissionDecorator validates that the validator commission is always
// greater than or equal to the min commission rate
type ValidateMinCommissionDecorator struct {
	sk stakingkeeper.Keeper
}

// ValidateMinCommissionDecorator creates a new ValidateMinCommissionDecorator
func NewValidateMinCommissionDecorator(sk stakingkeeper.Keeper) ValidateMinCommissionDecorator {
	return ValidateMinCommissionDecorator{
		sk: sk,
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

// validateMsg checks if the tx contains one of the following msg types & errors if the validator's commission rate is below the min threshold
// For validators: create validator or edit validator
// For delegators: delegate to validator or withdraw delegator rewards
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
		if val.GetCommission().LT(minCommission) {
			return sdkerrors.Wrapf(
				sdkerrors.ErrInvalidRequest,
				"cannot delegate to validator with commission lower than minimum of %s", minCommission)
		}

		validatorTokens := sdk.NewDecFromInt(val.GetTokens())
		stakingSupply := sdk.NewDecFromInt(vcd.sk.StakingTokenSupply(ctx))
		var votingPower = validatorTokens.Quo(stakingSupply).Mul(sdk.NewDec(100))
		if err != nil {
			return err
		}
		if votingPower.GTE(maxVotingPower) {
			return sdkerrors.Wrapf(
				sdkerrors.ErrInvalidRequest,
				"Cannot delegate to a validator with %s or higher voting power, please choose another validator", maxVotingPower)
		}
	case *stakingtypes.MsgBeginRedelegate:
		val, err := vcd.getValidator(ctx, msg.ValidatorDstAddress)
		if err != nil {
			return err
		}
		if val.GetCommission().LT(minCommission) {
			return sdkerrors.Wrapf(
				sdkerrors.ErrInvalidRequest,
				"cannot redelegate to validator with commission lower than minimum of %s", minCommission)
		}

		validatorTokens := sdk.NewDecFromInt(val.GetTokens())
		stakingSupply := sdk.NewDecFromInt(vcd.sk.StakingTokenSupply(ctx))
		var votingPower = validatorTokens.Quo(stakingSupply).Mul(sdk.NewDec(100))
		if err != nil {
			return err
		}
		if votingPower.GTE(maxVotingPower) {
			return sdkerrors.Wrapf(
				sdkerrors.ErrInvalidRequest,
				"Cannot delegate to a validator with %s or higher voting power, please choose another validator", maxVotingPower)
		}
	}
	return nil
}
