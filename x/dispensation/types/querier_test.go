package types_test

import (
	"testing"

	"github.com/Sifchain/sifnode/x/dispensation/types"
	"github.com/stretchr/testify/assert"
)

func TestQueryRecordsByDistributionName(t *testing.T) {
	distributionName := types.QueryRecordsByDistrName
	status := types.DistributionStatus_DISTRIBUTION_STATUS_PENDING
	result := types.NewQueryRecordsByDistributionName(distributionName, status)
	assert.Equal(t, distributionName, result.DistributionName)
	assert.Equal(t, status, result.Status)
}

func TestQueryRecordsByRecipientAddr(t *testing.T) {
	address := types.QueryRecordsByRecipient
	result := types.NewQueryRecordsByRecipientAddr(address)
	assert.Equal(t, address, result.Address)
}

func TestQueryUserClaims(t *testing.T) {

	userClaimtype := types.DistributionType_DISTRIBUTION_TYPE_UNSPECIFIED
	result := types.NewQueryUserClaims(userClaimtype)
	assert.Equal(t, userClaimtype, result.UserClaimType)
}
