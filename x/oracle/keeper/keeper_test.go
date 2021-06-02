package keeper_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Sifchain/sifnode/app"
	"github.com/Sifchain/sifnode/x/ethbridge/test"
	"github.com/Sifchain/sifnode/x/oracle/types"
)

const (
	TestID                     = "oracleID"
	AlternateTestID            = "altOracleID"
	TestString                 = "{value: 5}"
	AlternateTestString        = "{value: 7}"
	AnotherAlternateTestString = "{value: 9}"
)

func TestCreateGetProphecy(t *testing.T) {
	ctx, _, _, _, oracleKeeper, _, whitelist := test.CreateTestKeepers(t, 0.7, []int64{3, 7}, "")
	validatorAddresses := whitelist.GetAllValidators()
	validator1Pow3 := validatorAddresses[0]

	//Test normal Creation
	oracleClaim := types.NewClaim(TestID, validator1Pow3.String(), TestString)
	status, err := oracleKeeper.ProcessClaim(ctx, 1, oracleClaim)
	require.NoError(t, err)
	require.Equal(t, status.Text, types.StatusText_STATUS_TEXT_PENDING)

	//Test bad Creation with blank id
	oracleClaim = types.NewClaim("", validator1Pow3.String(), TestString)
	status, err = oracleKeeper.ProcessClaim(ctx, 1, oracleClaim)
	require.Error(t, err)

	//Test bad Creation with blank claim
	oracleClaim = types.NewClaim(TestID, validator1Pow3.String(), "")
	status, err = oracleKeeper.ProcessClaim(ctx, 1, oracleClaim)
	require.Error(t, err)

	//Test retrieval
	prophecy, found := oracleKeeper.GetProphecy(ctx, TestID)
	require.True(t, found)
	require.Equal(t, prophecy.ID, TestID)
	require.Equal(t, prophecy.Status.Text, types.StatusText_STATUS_TEXT_PENDING)
	require.Equal(t, prophecy.ClaimValidators[TestString][0], validator1Pow3)
	require.Equal(t, prophecy.ValidatorClaims[validator1Pow3.String()], TestString)
}

func TestBadConsensusForOracle(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()
	_, _, _, _, _, _, _ = test.CreateTestKeepers(t, 0, []int64{10}, "")
	_, _, _, _, _, _, _ = test.CreateTestKeepers(t, 1.2, []int64{10}, "")
}

func TestBadMsgs(t *testing.T) {
	ctx, _, _, _, oracleKeeper, _, whitelist := test.CreateTestKeepers(t, 0.6, []int64{3, 3}, "")
	validatorAddresses := whitelist.GetAllValidators()
	validator1Pow3 := validatorAddresses[0]

	//Test empty claim
	oracleClaim := types.NewClaim(TestID, validator1Pow3.String(), "")
	status, err := oracleKeeper.ProcessClaim(ctx, 1, oracleClaim)
	require.Error(t, err)
	require.Equal(t, status.FinalClaim, "")
	require.True(t, strings.Contains(err.Error(), "claim cannot be empty string"))

	//Test normal Creation
	oracleClaim = types.NewClaim(TestID, validator1Pow3.String(), TestString)
	status, err = oracleKeeper.ProcessClaim(ctx, 1, oracleClaim)
	require.NoError(t, err)
	require.Equal(t, status.Text, types.StatusText_STATUS_TEXT_PENDING)

	//Test duplicate message
	oracleClaim = types.NewClaim(TestID, validator1Pow3.String(), TestString)
	status, err = oracleKeeper.ProcessClaim(ctx, 1, oracleClaim)
	require.Error(t, err)
	require.True(t, strings.Contains(err.Error(), "already processed message from validator for this id"))

	//Test second but non duplicate message
	oracleClaim = types.NewClaim(TestID, validator1Pow3.String(), AlternateTestString)
	status, err = oracleKeeper.ProcessClaim(ctx, 1, oracleClaim)
	require.Error(t, err)
	require.True(t, strings.Contains(err.Error(), "already processed message from validator for this id"))
}

func TestSuccessfulProphecy(t *testing.T) {
	ctx, _, _, _, oracleKeeper, _, whitelist := test.CreateTestKeepers(t, 0.6, []int64{3, 3, 4}, "")

	validatorAddresses := whitelist.GetAllValidators()
	validator1Pow3 := validatorAddresses[0]
	validator2Pow3 := validatorAddresses[1]
	validator3Pow4 := validatorAddresses[2]

	//Test first claim
	oracleClaim := types.NewClaim(TestID, validator1Pow3.String(), TestString)
	status, err := oracleKeeper.ProcessClaim(ctx, 1, oracleClaim)
	require.NoError(t, err)
	require.Equal(t, status.Text, types.StatusText_STATUS_TEXT_PENDING)

	//Test second claim completes and finalizes to success
	oracleClaim = types.NewClaim(TestID, validator2Pow3.String(), TestString)
	status, err = oracleKeeper.ProcessClaim(ctx, 1, oracleClaim)
	require.NoError(t, err)
	require.Equal(t, status.Text, types.StatusText_STATUS_TEXT_SUCCESS)
	require.Equal(t, status.FinalClaim, TestString)

	//Test third claim not possible
	oracleClaim = types.NewClaim(TestID, validator3Pow4.String(), TestString)
	status, err = oracleKeeper.ProcessClaim(ctx, 1, oracleClaim)
	require.Error(t, err)
	require.True(t, strings.Contains(err.Error(), "prophecy already finalized"))
}

func TestSuccessfulProphecyWithDisagreement(t *testing.T) {
	ctx, _, _, _, oracleKeeper, _, whitelist := test.CreateTestKeepers(t, 0.6, []int64{3, 3, 4}, "")

	validatorAddresses := whitelist.GetAllValidators()
	validator1Pow3 := validatorAddresses[0]
	validator2Pow3 := validatorAddresses[1]
	validator3Pow4 := validatorAddresses[2]

	//Test first claim
	oracleClaim := types.NewClaim(TestID, validator1Pow3.String(), TestString)
	status, err := oracleKeeper.ProcessClaim(ctx, 1, oracleClaim)
	require.NoError(t, err)
	require.Equal(t, status.Text, types.StatusText_STATUS_TEXT_PENDING)

	//Test second disagreeing claim processed fine
	oracleClaim = types.NewClaim(TestID, validator2Pow3.String(), AlternateTestString)
	status, err = oracleKeeper.ProcessClaim(ctx, 1, oracleClaim)
	require.NoError(t, err)
	require.Equal(t, status.Text, types.StatusText_STATUS_TEXT_PENDING)

	//Test third claim agrees and finalizes to success
	oracleClaim = types.NewClaim(TestID, validator3Pow4.String(), TestString)
	status, err = oracleKeeper.ProcessClaim(ctx, 1, oracleClaim)
	require.NoError(t, err)
	require.Equal(t, status.Text, types.StatusText_STATUS_TEXT_SUCCESS)
	require.Equal(t, status.FinalClaim, TestString)
}

func TestFailedProphecy(t *testing.T) {
	ctx, _, _, _, oracleKeeper, _, whitelist := test.CreateTestKeepers(t, 0.6, []int64{3, 3, 4}, "")

	validatorAddresses := whitelist.GetAllValidators()
	validator1Pow3 := validatorAddresses[0]
	validator2Pow3 := validatorAddresses[1]
	validator3Pow4 := validatorAddresses[2]

	//Test first claim
	oracleClaim := types.NewClaim(TestID, validator1Pow3.String(), TestString)
	status, err := oracleKeeper.ProcessClaim(ctx, 1, oracleClaim)
	require.NoError(t, err)
	require.Equal(t, status.Text, types.StatusText_STATUS_TEXT_PENDING)

	//Test second disagreeing claim processed fine
	oracleClaim = types.NewClaim(TestID, validator2Pow3.String(), AlternateTestString)
	status, err = oracleKeeper.ProcessClaim(ctx, 1, oracleClaim)
	require.NoError(t, err)
	require.Equal(t, status.Text, types.StatusText_STATUS_TEXT_PENDING)
	require.Equal(t, status.FinalClaim, "")

	//Test third disagreeing claim processed fine and prophecy fails
	oracleClaim = types.NewClaim(TestID, validator3Pow4.String(), AnotherAlternateTestString)
	status, err = oracleKeeper.ProcessClaim(ctx, 1, oracleClaim)
	require.NoError(t, err)
	require.Equal(t, status.Text, types.StatusText_STATUS_TEXT_FAILED)
	require.Equal(t, status.FinalClaim, "")
}

func TestPowerOverrule(t *testing.T) {
	//Testing with 2 validators but one has high enough power to overrule
	ctx, _, _, _, oracleKeeper, _, whitelist := test.CreateTestKeepers(t, 0.7, []int64{3, 7}, "")

	validatorAddresses := whitelist.GetAllValidators()
	validator1Pow3 := validatorAddresses[0]
	validator2Pow7 := validatorAddresses[1]

	//Test first claim
	oracleClaim := types.NewClaim(TestID, validator1Pow3.String(), TestString)
	status, err := oracleKeeper.ProcessClaim(ctx, 1, oracleClaim)
	require.NoError(t, err)
	require.Equal(t, status.Text, types.StatusText_STATUS_TEXT_PENDING)

	//Test second disagreeing claim processed fine and finalized to its bytes
	oracleClaim = types.NewClaim(TestID, validator2Pow7.String(), AlternateTestString)
	status, err = oracleKeeper.ProcessClaim(ctx, 1, oracleClaim)
	require.NoError(t, err)
	require.Equal(t, status.Text, types.StatusText_STATUS_TEXT_SUCCESS)
	require.Equal(t, status.FinalClaim, AlternateTestString)
}
func TestPowerAternate(t *testing.T) {
	//Test alternate power setup with validators of 5/4/3/9 and total power 22 and 12/21 required
	ctx, _, _, _, oracleKeeper, _, whitelist := test.CreateTestKeepers(t, 0.571, []int64{5, 4, 3, 9}, "")
	validatorAddresses := whitelist.GetAllValidators()
	validator1Pow5 := validatorAddresses[0]
	validator2Pow4 := validatorAddresses[1]
	validator3Pow3 := validatorAddresses[2]
	validator4Pow9 := validatorAddresses[3]

	//Test claim by v1
	oracleClaim := types.NewClaim(TestID, validator1Pow5.String(), TestString)
	status, err := oracleKeeper.ProcessClaim(ctx, 1, oracleClaim)
	require.NoError(t, err)
	require.Equal(t, status.Text, types.StatusText_STATUS_TEXT_PENDING)

	//Test claim by v2
	oracleClaim = types.NewClaim(TestID, validator2Pow4.String(), TestString)
	status, err = oracleKeeper.ProcessClaim(ctx, 1, oracleClaim)
	require.NoError(t, err)
	require.Equal(t, status.Text, types.StatusText_STATUS_TEXT_PENDING)

	//Test alternate claim by v4
	oracleClaim = types.NewClaim(TestID, validator4Pow9.String(), AlternateTestString)
	status, err = oracleKeeper.ProcessClaim(ctx, 1, oracleClaim)
	require.NoError(t, err)
	require.Equal(t, status.Text, types.StatusText_STATUS_TEXT_PENDING)

	//Test finalclaim by v3
	oracleClaim = types.NewClaim(TestID, validator3Pow3.String(), TestString)
	status, err = oracleKeeper.ProcessClaim(ctx, 1, oracleClaim)
	require.NoError(t, err)
	require.Equal(t, status.Text, types.StatusText_STATUS_TEXT_SUCCESS)
	require.Equal(t, status.FinalClaim, TestString)
}

func TestMultipleProphecies(t *testing.T) {
	//Test multiple prophecies running in parallel work fine as expected
	ctx, _, _, _, oracleKeeper, _, whitelist := test.CreateTestKeepers(t, 0.7, []int64{3, 7}, "")
	validatorAddresses := whitelist.GetAllValidators()

	validator1Pow3 := validatorAddresses[0]
	validator2Pow7 := validatorAddresses[1]

	//Test claim on first id with first validator
	oracleClaim := types.NewClaim(TestID, validator1Pow3.String(), TestString)
	status, err := oracleKeeper.ProcessClaim(ctx, 1, oracleClaim)
	require.NoError(t, err)
	require.Equal(t, status.Text, types.StatusText_STATUS_TEXT_PENDING)

	//Test claim on second id with second validator
	oracleClaim = types.NewClaim(AlternateTestID, validator2Pow7.String(), AlternateTestString)
	status, err = oracleKeeper.ProcessClaim(ctx, 1, oracleClaim)
	require.NoError(t, err)
	require.Equal(t, status.Text, types.StatusText_STATUS_TEXT_SUCCESS)
	require.Equal(t, status.FinalClaim, AlternateTestString)

	//Test claim on first id with second validator
	oracleClaim = types.NewClaim(TestID, validator2Pow7.String(), TestString)
	status, err = oracleKeeper.ProcessClaim(ctx, 1, oracleClaim)
	require.NoError(t, err)
	require.Equal(t, status.Text, types.StatusText_STATUS_TEXT_SUCCESS)
	require.Equal(t, status.FinalClaim, TestString)

	//Test claim on second id with first validator
	oracleClaim = types.NewClaim(AlternateTestID, validator1Pow3.String(), AlternateTestString)
	status, err = oracleKeeper.ProcessClaim(ctx, 1, oracleClaim)
	require.Error(t, err)
	require.True(t, strings.Contains(err.Error(), "prophecy already finalized"))
}

func TestNonValidator(t *testing.T) {
	//Test multiple prophecies running in parallel work fine as expected
	ctx, _, _, _, oracleKeeper, _, _ := test.CreateTestKeepers(t, 0.7, []int64{3, 7}, "")

	addresses := app.CreateRandomAccounts(10)
	testValidatorAddresses := app.ConvertAddrsToValAddrs(addresses)
	inActiveValidatorAddress := testValidatorAddresses[9]

	//Test claim on first id with first validator
	oracleClaim := types.NewClaim(TestID, inActiveValidatorAddress.String(), TestString)
	_, err := oracleKeeper.ProcessClaim(ctx, 1, oracleClaim)
	require.Error(t, err)
	require.True(t, strings.Contains(err.Error(), "validator not in white list"))
}
