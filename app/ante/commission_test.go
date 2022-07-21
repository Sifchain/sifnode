package ante_test

import (
	"testing"
	"time"

	sifapp "github.com/Sifchain/sifnode/app"
	sifAnte "github.com/Sifchain/sifnode/app/ante"
	sdk "github.com/cosmos/cosmos-sdk/types"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"

	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/require"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
)

func TestAnte_CalculateProjectedVotingPower(t *testing.T) {
	testcases := []struct {
		name                   string
		delegateAmount         sdk.Int
		validatorTokens        sdk.Int
		validatorBondStatus    stakingtypes.BondStatus
		weakestValidatorTokens sdk.Int
		expectedVotingPower    sdk.Dec
	}{
		{
			name:                   "not bonded, projected to be more powerful than weakest",
			delegateAmount:         sdk.NewIntFromUint64(100),
			validatorTokens:        sdk.NewIntFromUint64(1000),
			validatorBondStatus:    stakingtypes.Unbonded,
			weakestValidatorTokens: sdk.NewIntFromUint64(1050),
			expectedVotingPower:    sdk.MustNewDecFromStr("1.0000000000000"),
		},
		{
			name:                   "not bonded, NOT projected to be more powerful than weakest",
			delegateAmount:         sdk.NewIntFromUint64(100),
			validatorTokens:        sdk.NewIntFromUint64(1000),
			validatorBondStatus:    stakingtypes.Unbonded,
			weakestValidatorTokens: sdk.NewIntFromUint64(2000),
			expectedVotingPower:    sdk.MustNewDecFromStr("0.0000000"),
		},
		{
			name:                   "bonded",
			delegateAmount:         sdk.NewIntFromUint64(100),
			validatorTokens:        sdk.NewIntFromUint64(1000),
			validatorBondStatus:    stakingtypes.Bonded,
			weakestValidatorTokens: sdk.NewIntFromUint64(2000),
			expectedVotingPower:    sdk.MustNewDecFromStr("0.354838709677419355"),
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			app := sifapp.Setup(false)
			ctx := app.BaseApp.NewContext(false, tmproto.Header{})
			vcd := sifAnte.NewValidateMinCommissionDecorator(app.StakingKeeper)

			delegatorAddresses := sifapp.AddTestAddrs(app, ctx, 1, tc.weakestValidatorTokens.Add(tc.validatorTokens))
			delegatorAddress := delegatorAddresses[0]
			pubKeys := sifapp.CreateTestPubKeys(2)

			// create weakest validator
			createValidator(app, ctx, tc.weakestValidatorTokens, stakingtypes.Bonded, delegatorAddress, pubKeys[0])

			// create validator for which we'll calculate projected voting power
			valAddress := createValidator(app, ctx, tc.validatorTokens, tc.validatorBondStatus, delegatorAddress, pubKeys[1])
			validator, _ := app.StakingKeeper.GetValidator(ctx, valAddress)

			calcVotingPower := vcd.CalculateProjectedVotingPower(ctx, validator, tc.delegateAmount.ToDec())
			require.Equal(t, tc.expectedVotingPower.String(), calcVotingPower.String())
		})
	}
}

//TODO: edge case if num validators < max validators

func createValidator(app *sifapp.SifchainApp, ctx sdk.Context, delegateAmount sdk.Int, status stakingtypes.BondStatus, delegatorAddress sdk.AccAddress, valPubKey cryptotypes.PubKey) sdk.ValAddress {
	pkAny, err := codectypes.NewAnyWithValue(valPubKey)
	if err != nil {
		panic(err)
	}
	operatorAddress := sdk.ValAddress(valPubKey.Address())

	validator := stakingtypes.Validator{
		OperatorAddress: operatorAddress.String(),
		ConsensusPubkey: pkAny,
		Jailed:          false,
		Status:          status,
		Tokens:          sdk.ZeroInt(),
		DelegatorShares: sdk.ZeroDec(),
		Description: stakingtypes.Description{
			Moniker:         "moniker",
			Identity:        "id",
			Website:         "www",
			SecurityContact: "alice",
			Details:         "details",
		},
		UnbondingHeight: 0,
		UnbondingTime:   time.Time{},
		Commission: stakingtypes.Commission{
			CommissionRates: stakingtypes.CommissionRates{
				Rate:          sdk.NewDecWithPrec(5, 2),
				MaxRate:       sdk.NewDecWithPrec(10, 2),
				MaxChangeRate: sdk.NewDecWithPrec(1, 2),
			},
			UpdateTime: time.Time{},
		},
		MinSelfDelegation: sdk.NewInt(1),
	}

	app.StakingKeeper.SetValidator(ctx, validator)
	err = app.StakingKeeper.SetValidatorByConsAddr(ctx, validator)
	if err != nil {
		panic(err)
	}
	app.StakingKeeper.SetNewValidatorByPowerIndex(ctx, validator)
	app.StakingKeeper.AfterValidatorCreated(ctx, validator.GetOperator())

	_, err = app.StakingKeeper.Delegate(ctx, delegatorAddress, delegateAmount, stakingtypes.Unbonded, validator, true)
	if err != nil {
		panic(err)
	}

	return operatorAddress
}
