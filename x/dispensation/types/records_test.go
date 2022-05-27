package types_test

import (
	"testing"
	"time"

	"github.com/Sifchain/sifnode/x/dispensation/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	gogotypes "github.com/gogo/protobuf/types"
	"github.com/stretchr/testify/assert"
)

func TestNewDistributionRecord(t *testing.T) {
	distributionType := types.DistributionType_DISTRIBUTION_TYPE_LIQUIDITY_MINING
	status := types.DistributionStatus_DISTRIBUTION_STATUS_PENDING
	distributionName := types.AttributeKeyDistributionName
	runner := types.AttributeKeyDistributionRunner
	recipientAddress := types.AttributeKeyDistributionRecordAddress
	Coins := sdk.Coins{sdk.Coin{
		Denom:  "rowan",
		Amount: sdk.NewInt(20),
	}}
	distributionStartHeight := int64(1)
	distributionCompletedHeight := int64(10)
	result := types.NewDistributionRecord(status, distributionType, distributionName, recipientAddress, Coins, distributionStartHeight, distributionCompletedHeight, runner)

	assert.Equal(t, distributionName, result.DistributionName)
	assert.Equal(t, distributionType, result.DistributionType)
	assert.Equal(t, status, result.DistributionStatus)
	assert.Equal(t, Coins, result.Coins)
	assert.Equal(t, runner, result.AuthorizedRunner)
	assert.Equal(t, recipientAddress, result.RecipientAddress)
	assert.Equal(t, distributionStartHeight, result.DistributionStartHeight)
	assert.Equal(t, distributionCompletedHeight, result.DistributionCompletedHeight)

}

func TestNewDistributionRecord_validateTrue(t *testing.T) {
	distributionType := types.DistributionType_DISTRIBUTION_TYPE_LIQUIDITY_MINING
	status := types.DistributionStatus_DISTRIBUTION_STATUS_PENDING
	distributionName := types.AttributeKeyDistributionName
	runner := types.AttributeKeyDistributionRunner
	recipientAddress := types.AttributeKeyDistributionRecordAddress
	Coins := sdk.Coins{sdk.Coin{
		Denom:  "rowan",
		Amount: sdk.NewInt(20)}}
	distributionStartHeight := int64(1)
	distributionCompletedHeight := int64(10)
	result := types.NewDistributionRecord(status, distributionType, distributionName, recipientAddress, Coins, distributionStartHeight, distributionCompletedHeight, runner)
	bool := result.Validate()
	assert.True(t, bool)
}

func TestNewDistributionRecord_validateEmptyAddress(t *testing.T) {
	distributionType := types.DistributionType_DISTRIBUTION_TYPE_LIQUIDITY_MINING
	status := types.DistributionStatus_DISTRIBUTION_STATUS_PENDING
	distributionName := types.AttributeKeyDistributionName
	runner := types.AttributeKeyDistributionRunner
	recipientAddress := ""
	Coins := sdk.Coins{sdk.Coin{
		Denom:  "denom",
		Amount: sdk.NewInt(20)}}
	distributionStartHeight := int64(1)
	distributionCompletedHeight := int64(10)
	result := types.NewDistributionRecord(status, distributionType, distributionName, recipientAddress, Coins, distributionStartHeight, distributionCompletedHeight, runner)
	bool := result.Validate()
	assert.False(t, bool)
}

func TestNewDistributionRecord_validateCoinIsInvalid(t *testing.T) {
	distributionType := types.DistributionType_DISTRIBUTION_TYPE_LIQUIDITY_MINING
	status := types.DistributionStatus_DISTRIBUTION_STATUS_PENDING
	distributionName := types.AttributeKeyDistributionName
	runner := types.AttributeKeyDistributionRunner
	recipientAddress := types.AttributeKeyDistributionRecordAddress
	Coins := sdk.Coins{sdk.Coin{
		Denom:  "",
		Amount: sdk.NewInt(20)}}
	distributionStartHeight := int64(1)
	distributionCompletedHeight := int64(10)
	result := types.NewDistributionRecord(status, distributionType, distributionName, recipientAddress, Coins, distributionStartHeight, distributionCompletedHeight, runner)
	bool := result.Validate()
	assert.False(t, bool)
}

func TestNewDistributionRecord_Add(t *testing.T) {
	distributionType := types.DistributionType_DISTRIBUTION_TYPE_LIQUIDITY_MINING
	status := types.DistributionStatus_DISTRIBUTION_STATUS_PENDING
	distributionName := types.AttributeKeyDistributionName
	runner := types.AttributeKeyDistributionRunner
	recipientAddress := types.AttributeKeyDistributionRecordAddress
	Coins := sdk.Coins{sdk.Coin{
		Denom:  "",
		Amount: sdk.NewInt(20)}}
	distributionStartHeight := int64(1)
	distributionCompletedHeight := int64(10)

	result := types.NewDistributionRecord(status, distributionType, distributionName, recipientAddress, Coins, distributionStartHeight, distributionCompletedHeight, runner)
	result2 := types.NewDistributionRecord(status, distributionType, distributionName, recipientAddress, Coins, distributionStartHeight, distributionCompletedHeight, runner)
	result.Coins = result.Coins.Add(result2.Coins...)
	output := result2.Add(result2)

	assert.Equal(t, result, output)
	assert.Equal(t, result.Coins, output.Coins)
}

func TestNewDistributionRecord_DoesTypeSupportClaim(t *testing.T) {
	distributionType := types.DistributionType_DISTRIBUTION_TYPE_LIQUIDITY_MINING
	status := types.DistributionStatus_DISTRIBUTION_STATUS_PENDING
	distributionName := types.AttributeKeyDistributionName
	runner := types.AttributeKeyDistributionRunner
	recipientAddress := types.AttributeKeyDistributionRecordAddress
	Coins := sdk.Coins{sdk.Coin{
		Denom:  "rowan",
		Amount: sdk.NewInt(20)}}
	distributionStartHeight := int64(1)
	distributionCompletedHeight := int64(10)
	result := types.NewDistributionRecord(status, distributionType, distributionName, recipientAddress, Coins, distributionStartHeight, distributionCompletedHeight, runner)
	bool := result.DoesTypeSupportClaim()
	assert.True(t, bool)

}

func TestNewDistributionRecord_DoesTypeSupportClaim_False(t *testing.T) {
	distributionType := types.DistributionType_DISTRIBUTION_TYPE_AIRDROP
	status := types.DistributionStatus_DISTRIBUTION_STATUS_PENDING
	distributionName := types.AttributeKeyDistributionName
	runner := types.AttributeKeyDistributionRunner
	recipientAddress := types.AttributeKeyDistributionRecordAddress
	Coins := sdk.Coins{sdk.Coin{
		Denom:  "rowan",
		Amount: sdk.NewInt(20)}}
	distributionStartHeight := int64(1)
	distributionCompletedHeight := int64(10)
	result := types.NewDistributionRecord(status, distributionType, distributionName, recipientAddress, Coins, distributionStartHeight, distributionCompletedHeight, runner)
	bool := result.DoesTypeSupportClaim()
	assert.False(t, bool)

}

func TestNewUserClaim(t *testing.T) {
	UserAddress := types.AttributeKeyDistributionRecordAddress
	userClaimType := types.DistributionType_DISTRIBUTION_TYPE_LIQUIDITY_MINING
	tm := time.Time{}
	tProto, err := gogotypes.TimestampProto(tm)

	UserClaim, error := types.NewUserClaim(UserAddress, userClaimType, tm)
	assert.Equal(t, UserAddress, UserClaim.UserAddress)
	assert.Equal(t, userClaimType, UserClaim.UserClaimType)
	assert.Equal(t, tProto, UserClaim.UserClaimTime)
	assert.Equal(t, err, error)

}

func TestNewUserClaim_ValidateFalse(t *testing.T) {
	UserAddress := ""
	userClaimType := types.DistributionType_DISTRIBUTION_TYPE_LIQUIDITY_MINING
	tm := time.Time{}
	UserClaim, error := types.NewUserClaim(UserAddress, userClaimType, tm)
	bool := UserClaim.Validate()
	assert.False(t, bool)
	assert.Empty(t, error)
}
func TestNewUserClaim_ValidateTrue(t *testing.T) {
	UserAddress := types.AttributeKeyDistributionRecordAddress
	userClaimType := types.DistributionType_DISTRIBUTION_TYPE_LIQUIDITY_MINING
	tm := time.Time{}
	UserClaim, error := types.NewUserClaim(UserAddress, userClaimType, tm)
	bool := UserClaim.Validate()
	assert.True(t, bool)
	assert.Empty(t, error)
}
func TestNewDistribution(t *testing.T) {
	distributiontype := types.DistributionType_DISTRIBUTION_TYPE_LIQUIDITY_MINING
	name := types.AttributeKeyDistributionName
	authorizedRunner := types.AttributeKeyDistributionRunner
	result := types.NewDistribution(distributiontype, name, authorizedRunner)
	assert.Equal(t, name, result.DistributionName)
	assert.Equal(t, distributiontype, result.DistributionType)
	assert.Equal(t, authorizedRunner, result.Runner)
}

func TestNewDistribution_validateTrue(t *testing.T) {
	distributiontype := types.DistributionType_DISTRIBUTION_TYPE_LIQUIDITY_MINING
	name := types.AttributeKeyDistributionName
	authorizedRunner := types.AttributeKeyDistributionRunner
	result := types.NewDistribution(distributiontype, name, authorizedRunner)
	bool := result.Validate()
	assert.True(t, bool)

}

func TestNewDistribution_validateFalse(t *testing.T) {
	distributiontype := types.DistributionType_DISTRIBUTION_TYPE_LIQUIDITY_MINING
	name := ""
	authorizedRunner := types.AttributeKeyDistributionRunner
	result := types.NewDistribution(distributiontype, name, authorizedRunner)
	bool := result.Validate()
	assert.False(t, bool)

}

func TestGetDistributionStatus_Completed(t *testing.T) {
	status := "Completed"
	result, output := types.GetDistributionStatus(status)
	Distributionstatus := int32(result)
	assert.NotEmpty(t, Distributionstatus)
	assert.True(t, output)

}

func TestGetDistributionStatus_Pending(t *testing.T) {
	status := "Pending"
	result, output := types.GetDistributionStatus(status)
	Distributionstatus := int32(result)
	assert.NotEmpty(t, Distributionstatus)
	assert.True(t, output)
}

func TestGetDistributionStatus_Failed(t *testing.T) {
	status := "Failed"
	result, output := types.GetDistributionStatus(status)
	Distributionstatus := int32(result)
	assert.NotEmpty(t, Distributionstatus)
	assert.True(t, output)
}

func TestGetDistributionStatus_Default(t *testing.T) {
	status := ""
	result, output := types.GetDistributionStatus(status)
	Distributionstatus := int32(result)
	assert.Empty(t, Distributionstatus)
	assert.False(t, output)
}

func TestGetClaimType_ValidatorSubsidy(t *testing.T) {
	claimType := "ValidatorSubsidy"
	result, output := types.GetClaimType(claimType)
	DistributionType := int32(result)
	assert.NotEmpty(t, DistributionType)
	assert.True(t, output)
}

func TestGetClaimType_LiquidityMining(t *testing.T) {
	claimType := "LiquidityMining"
	result, output := types.GetClaimType(claimType)
	DistributionType := int32(result)
	assert.NotEmpty(t, DistributionType)
	assert.True(t, output)
}

func TestGetClaimType_Default(t *testing.T) {
	claimType := ""
	result, output := types.GetClaimType(claimType)
	DistributionType := int32(result)
	assert.Empty(t, DistributionType)
	assert.False(t, output)
}

func TestGetDistributionTypeFromShortString_LiquidityMining(t *testing.T) {
	distributionType := "LiquidityMining"
	result, output := types.GetDistributionTypeFromShortString(distributionType)
	DistributionType := int32(result)
	assert.NotEmpty(t, DistributionType)
	assert.True(t, output)
}

func TestGetDistributionTypeFromShortString_Airdrop(t *testing.T) {
	distributionType := "Airdrop"
	result, output := types.GetDistributionTypeFromShortString(distributionType)
	DistributionType := int32(result)
	assert.NotEmpty(t, DistributionType)
	assert.True(t, output)
}

func TestGetDistributionTypeFromShortString_ValidatorSubsidy(t *testing.T) {
	distributionType := "ValidatorSubsidy"
	result, output := types.GetDistributionTypeFromShortString(distributionType)
	DistributionType := int32(result)
	assert.NotEmpty(t, DistributionType)
	assert.True(t, output)
}

func TestGetDistributionTypeFromShortString_default(t *testing.T) {
	distributionType := ""
	result, output := types.GetDistributionTypeFromShortString(distributionType)
	DistributionType := int32(result)
	assert.Empty(t, DistributionType)
	assert.False(t, output)
}

func TestIsValidDistributionType_airdrop(t *testing.T) {
	distributiontype := "DISTRIBUTION_TYPE_AIRDROP"

	result, output := types.IsValidDistributionType(distributiontype)
	returndistribution := types.DistributionType_DISTRIBUTION_TYPE_AIRDROP
	assert.NotEmpty(t, distributiontype)
	assert.Equal(t, result, returndistribution)
	assert.True(t, output)
}

func TestIsValidDistributionType_mining(t *testing.T) {
	distributiontype := "DISTRIBUTION_TYPE_LIQUIDITY_MINING"

	result, output := types.IsValidDistributionType(distributiontype)
	returndistribution := types.DistributionType_DISTRIBUTION_TYPE_LIQUIDITY_MINING
	assert.NotEmpty(t, distributiontype)
	assert.Equal(t, result, returndistribution)
	assert.True(t, output)
}

func TestIsValidDistributionType_ValidatorSubsidy(t *testing.T) {
	distributiontype := "DISTRIBUTION_TYPE_VALIDATOR_SUBSIDY"

	result, output := types.IsValidDistributionType(distributiontype)
	returndistribution := types.DistributionType_DISTRIBUTION_TYPE_VALIDATOR_SUBSIDY
	assert.NotEmpty(t, distributiontype)
	assert.Equal(t, result, returndistribution)
	assert.True(t, output)
}

func TestIsValidDistributionType_default(t *testing.T) {
	distributiontype := ""

	result, output := types.IsValidDistributionType(distributiontype)
	returndistribution := types.DistributionType_DISTRIBUTION_TYPE_UNSPECIFIED
	assert.Empty(t, distributiontype)
	assert.Equal(t, result, returndistribution)
	assert.False(t, output)
}

func Test_IsValidClaimType_liquidity_mining(t *testing.T) {
	claimtype := "DISTRIBUTION_TYPE_LIQUIDITY_MINING"

	result, output := types.IsValidClaimType(claimtype)
	returnclaimtype := types.DistributionType_DISTRIBUTION_TYPE_LIQUIDITY_MINING
	assert.NotEmpty(t, claimtype)
	assert.Equal(t, result, returnclaimtype)
	assert.True(t, output)
}

func Test_IsValidClaimType_ValidatorSubsidy(t *testing.T) {
	claimtype := "DISTRIBUTION_TYPE_VALIDATOR_SUBSIDY"

	result, output := types.IsValidClaimType(claimtype)
	returnclaimtype := types.DistributionType_DISTRIBUTION_TYPE_VALIDATOR_SUBSIDY
	assert.NotEmpty(t, claimtype)
	assert.Equal(t, result, returnclaimtype)
	assert.True(t, output)
}

func Test_IsValidClaimType_default(t *testing.T) {
	claimtype := ""

	result, output := types.IsValidClaimType(claimtype)
	returnclaimtype := types.DistributionType_DISTRIBUTION_TYPE_UNSPECIFIED
	assert.Empty(t, claimtype)
	assert.Equal(t, result, returnclaimtype)
	assert.False(t, output)
}
