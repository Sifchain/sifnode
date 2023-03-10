package types_test

import (
	"testing"

	"github.com/Sifchain/sifnode/x/dispensation/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

func TestQueryAllDistributionsResponse(t *testing.T) {
	d := types.Distribution{
		DistributionType: types.DistributionType_DISTRIBUTION_TYPE_LIQUIDITY_MINING,
		DistributionName: types.AttributeKeyDistributionName,
		Runner:           types.AttributeKeyDistributionRunner,
	}
	dp := &d
	distr := types.Distributions{Distributions: []*types.Distribution{dp}}
	height := int64(2)
	result := types.NewQueryAllDistributionsResponse(distr, height)
	assert.Equal(t, []*types.Distribution{dp}, result.Distributions)

}

func TestQueryRecordsByDistributionNameResponse(t *testing.T) {

	d := types.DistributionRecord{
		DistributionStatus: types.DistributionStatus_DISTRIBUTION_STATUS_COMPLETED,
		DistributionType:   types.DistributionType_DISTRIBUTION_TYPE_AIRDROP,
		DistributionName:   types.AttributeKeyDistributionName,
		RecipientAddress:   types.AttributeKeyDistributionRecordAddress,
		Coins: sdk.Coins{sdk.Coin{Denom: "rowan",
			Amount: sdk.NewInt(20)}},
		DistributionStartHeight:     int64(0),
		DistributionCompletedHeight: int64(10),
		AuthorizedRunner:            types.AttributeKeyDistributionRunner,
	}
	dp := &d
	distr := types.DistributionRecords{DistributionRecords: []*types.DistributionRecord{dp}}
	height := int64(2)
	result := types.NewQueryRecordsByDistributionNameResponse(distr, height)
	assert.Equal(t, []*types.DistributionRecord{dp}, result.DistributionRecords.DistributionRecords)
}

func TestQueryRecordsByRecipientAddrResponse(t *testing.T) {

	d := types.DistributionRecord{
		DistributionStatus: types.DistributionStatus_DISTRIBUTION_STATUS_COMPLETED,
		DistributionType:   types.DistributionType_DISTRIBUTION_TYPE_AIRDROP,
		DistributionName:   types.AttributeKeyDistributionName,
		RecipientAddress:   types.AttributeKeyDistributionRecordAddress,
		Coins: sdk.Coins{sdk.Coin{Denom: "rowan",
			Amount: sdk.NewInt(20)}},
		DistributionStartHeight:     int64(0),
		DistributionCompletedHeight: int64(10),
		AuthorizedRunner:            types.AttributeKeyDistributionRunner,
	}
	dp := &d
	distr := types.DistributionRecords{DistributionRecords: []*types.DistributionRecord{dp}}
	height := int64(2)
	result := types.NewQueryRecordsByRecipientAddrResponse(distr, height)
	assert.Equal(t, []*types.DistributionRecord{dp}, result.DistributionRecords.DistributionRecords)
}
