package types_test

import (
	"fmt"
	"testing"

	"github.com/Sifchain/sifnode/x/dispensation/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/stretchr/testify/assert"
)

func TestGetDistributionRecordKey_statusCompleted(t *testing.T) {
	status := types.DistributionStatus_DISTRIBUTION_STATUS_COMPLETED
	name := types.AttributeKeyDistributionName
	recipient := types.AttributeKeyDistributionRecordAddress
	distributionType := types.DistributionType_DISTRIBUTION_TYPE_LIQUIDITY_MINING
	key := []byte(fmt.Sprintf("%s_%d_%s", name, distributionType, recipient))
	DistributionRecordPrefixCompleted := []byte{0x011}
	result := types.GetDistributionRecordKey(status, name, recipient, distributionType)
	output := append(DistributionRecordPrefixCompleted, key...)
	assert.Equal(t, result, output)

}

func TestGetDistributionRecordKey_statusPending(t *testing.T) {
	status := types.DistributionStatus_DISTRIBUTION_STATUS_PENDING
	name := types.AttributeKeyDistributionName
	recipient := types.AttributeKeyDistributionRecordAddress
	distributionType := types.DistributionType_DISTRIBUTION_TYPE_LIQUIDITY_MINING
	DistributionRecordPrefixPending := []byte{0x000}
	key := []byte(fmt.Sprintf("%s_%d_%s", name, distributionType, recipient))
	result := types.GetDistributionRecordKey(status, name, recipient, distributionType)
	output := append(DistributionRecordPrefixPending, key...)
	assert.Equal(t, result, output)

}

func TestGetDistributionRecordKey_statusFailed(t *testing.T) {
	status := types.DistributionStatus_DISTRIBUTION_STATUS_FAILED
	name := types.AttributeKeyDistributionName
	recipient := types.AttributeKeyDistributionRecordAddress
	distributionType := types.DistributionType_DISTRIBUTION_TYPE_LIQUIDITY_MINING
	DistributionRecordPrefixFailed := []byte{0x012}
	key := []byte(fmt.Sprintf("%s_%d_%s", name, distributionType, recipient))
	result := types.GetDistributionRecordKey(status, name, recipient, distributionType)
	output := append(DistributionRecordPrefixFailed, key...)
	assert.Equal(t, result, output)
}

func TestGetDistributionRecordKey_statusDefault(t *testing.T) {
	status := types.DistributionStatus_DISTRIBUTION_STATUS_UNSPECIFIED
	name := types.AttributeKeyDistributionName
	recipient := types.AttributeKeyDistributionRecordAddress
	distributionType := types.DistributionType_DISTRIBUTION_TYPE_LIQUIDITY_MINING
	DistributionRecordPrefixCompleted := []byte{0x011}
	key := []byte(fmt.Sprintf("%s_%d_%s", name, distributionType, recipient))
	result := types.GetDistributionRecordKey(status, name, recipient, distributionType)
	output := append(DistributionRecordPrefixCompleted, key...)
	assert.Equal(t, result, output)
}

func TestGetDistributionsKey(t *testing.T) {
	distributionType := types.DistributionType_DISTRIBUTION_TYPE_LIQUIDITY_MINING
	authorizedRunner := types.AttributeKeyDistributionRunner
	name := types.AttributeKeyDistributionName
	DistributionsPrefix := []byte{0x01}
	key := []byte(fmt.Sprintf("%s_%d_%s", name, distributionType, authorizedRunner))
	result := types.GetDistributionsKey(name, distributionType, authorizedRunner)
	output := append(DistributionsPrefix, key...)
	assert.Equal(t, result, output)
}

func TestGetUserClaimKey(t *testing.T) {
	userClaimType := types.DistributionType_DISTRIBUTION_TYPE_LIQUIDITY_MINING
	userAddress := types.AttributeKeyDistributionRecordAddress
	UserClaimPrefix := []byte{0x02}
	key := []byte(fmt.Sprintf("%s_%d", userAddress, userClaimType))
	result := types.GetUserClaimKey(userAddress, userClaimType)
	output := append(UserClaimPrefix, key...)

	assert.Equal(t, result, output)

}

func TestGetDistributionModuleAddress(t *testing.T) {
	moduleAddress := authtypes.NewModuleAddress(types.ModuleName)
	output := types.GetDistributionModuleAddress()
	assert.Equal(t, moduleAddress, output)

}
