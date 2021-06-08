package keeper_test

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"path/filepath"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/Sifchain/sifnode/x/dispensation/keeper"
	"github.com/Sifchain/sifnode/x/dispensation/test"
	"github.com/Sifchain/sifnode/x/dispensation/types"
	"github.com/Sifchain/sifnode/x/dispensation/types/legacy"
)

func TestUpgradeDistributionRecords(t *testing.T) {
	output := test.CreatOutputList(2, "100")

	legacyRecords := []legacy.DistributionRecord084{
		{
			ClaimStatus:                 1,
			DistributionName:            "first",
			RecipientAddress:            output[0].Address,
			Coins:                       sdk.NewCoins(sdk.NewCoin("rowan", sdk.NewInt(100))),
			DistributionStartHeight:     2,
			DistributionCompletedHeight: 3,
		},
	}

	upgradedRecords := []types.DistributionRecord{
		{
			DistributionStatus:          1,
			DistributionName:            "2_first",
			DistributionType:            types.Airdrop,
			RecipientAddress:            output[0].Address,
			Coins:                       sdk.NewCoins(sdk.NewCoin("rowan", sdk.NewInt(100))),
			DistributionStartHeight:     2,
			DistributionCompletedHeight: 3,
		},
	}

	var tt = []struct {
		name     string
		records  []legacy.DistributionRecord084
		upgraded []types.DistributionRecord
	}{
		{"success", legacyRecords, upgradedRecords},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			sifapp, ctx := test.CreateTestApp(false)

			for _, dr := range tc.records {
				sifapp.DispensationKeeper.Set(ctx,
					types.GetDistributionRecordKey(dr.DistributionName, dr.RecipientAddress.String(), types.DistributionTypeUnknown.String()),
					sifapp.Codec().MustMarshalBinaryBare(dr),
				)
			}

			keeper.UpgradeDistributionRecords(ctx, sifapp.DispensationKeeper)

			var got []types.DistributionRecord
			iterator := sifapp.DispensationKeeper.GetDistributionRecordsIterator(ctx)
			defer iterator.Close()
			for ; iterator.Valid(); iterator.Next() {
				var dr types.DistributionRecord
				bytesValue := iterator.Value()
				sifapp.Codec().MustUnmarshalBinaryBare(bytesValue, &dr)
				got = append(got, dr)
			}

			require.Equal(t, tc.upgraded, got)
		})
	}
}

func TestUpgradeDistributions(t *testing.T) {
	legacyDistributions := []legacy.Distribution084{
		{
			DistributionName: "first",
			DistributionType: types.LiquidityMining,
		},
	}
	upgradedDistributions := []types.Distribution{
		{
			DistributionName: "first",
			DistributionType: types.LiquidityMining,
		},
	}
	var tt = []struct {
		name     string
		records  []legacy.Distribution084
		upgraded []types.Distribution
	}{
		{"success", legacyDistributions, upgradedDistributions},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			sifapp, ctx := test.CreateTestApp(false)

			for _, dr := range tc.records {
				sifapp.DispensationKeeper.Set(ctx,
					types.GetDistributionsKey(dr.DistributionName, dr.DistributionType),
					sifapp.Codec().MustMarshalBinaryBare(dr),
				)
			}

			keeper.UpgradeDistributions(ctx, sifapp.DispensationKeeper)

			var got []types.Distribution
			iterator := sifapp.DispensationKeeper.GetDistributionIterator(ctx)
			defer iterator.Close()
			for ; iterator.Valid(); iterator.Next() {
				var d types.Distribution
				bytesValue := iterator.Value()
				sifapp.Codec().MustUnmarshalBinaryBare(bytesValue, &d)
				got = append(got, d)
			}

			require.Equal(t, tc.upgraded, got)
		})
	}
}
func SetConfig() {
	const (
		AccountAddressPrefix = "sif"
	)
	var (
		AccountPubKeyPrefix    = AccountAddressPrefix + "pub"
		ValidatorAddressPrefix = AccountAddressPrefix + "valoper"
		ValidatorPubKeyPrefix  = AccountAddressPrefix + "valoperpub"
		ConsNodeAddressPrefix  = AccountAddressPrefix + "valcons"
		ConsNodePubKeyPrefix   = AccountAddressPrefix + "valconspub"
	)
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(AccountAddressPrefix, AccountPubKeyPrefix)
	config.SetBech32PrefixForValidator(ValidatorAddressPrefix, ValidatorPubKeyPrefix)
	config.SetBech32PrefixForConsensusNode(ConsNodeAddressPrefix, ConsNodePubKeyPrefix)
	config.Seal()
}
func TestUpgradeDistributionRecords_BetanetData(t *testing.T) {
	SetConfig()
	legacyRecords := []legacy.DistributionRecord084{}
	file, err := filepath.Abs("../../../betanetRecords-complete.json")
	assert.NoError(t, err)
	o, err := ioutil.ReadFile(file)
	assert.NoError(t, err)
	err = json.Unmarshal(o, &legacyRecords)
	assert.NoError(t, err)
	upgradedRecords := []types.DistributionRecord{}

	for _, legacyRecord := range legacyRecords {
		distributionName := fmt.Sprintf("%d_%s", legacyRecord.DistributionStartHeight, legacyRecord.DistributionName)
		upgradedRecord := types.DistributionRecord{
			DistributionStatus:          types.Completed,
			DistributionName:            distributionName,
			DistributionType:            types.Airdrop,
			RecipientAddress:            legacyRecord.RecipientAddress,
			Coins:                       legacyRecord.Coins,
			DistributionStartHeight:     legacyRecord.DistributionStartHeight,
			DistributionCompletedHeight: legacyRecord.DistributionCompletedHeight,
		}
		upgradedRecords = append(upgradedRecords, upgradedRecord)
	}
	var tt = []struct {
		name            string
		legacyRecords   []legacy.DistributionRecord084
		upgradedRecords []types.DistributionRecord
	}{
		{"betanet-records", legacyRecords, upgradedRecords},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			sifapp, ctx := test.CreateTestApp(false)

			for _, dr := range tc.legacyRecords {
				sifapp.DispensationKeeper.Set(ctx,
					types.GetDistributionRecordKey(dr.DistributionName, dr.RecipientAddress.String(), types.DistributionTypeUnknown.String()),
					sifapp.Codec().MustMarshalBinaryBare(dr),
				)
			}
			keeper.UpgradeDistributionRecords(ctx, sifapp.DispensationKeeper)
			var got []types.DistributionRecord
			iterator := sifapp.DispensationKeeper.GetDistributionRecordsIterator(ctx)
			defer iterator.Close()
			for ; iterator.Valid(); iterator.Next() {
				var dr types.DistributionRecord
				bytesValue := iterator.Value()
				sifapp.Codec().MustUnmarshalBinaryBare(bytesValue, &dr)
				got = append(got, dr)
			}
			require.Equal(t, tc.upgradedRecords, got)
		})
	}
}
