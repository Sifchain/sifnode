package v42

import (
	gogotypes "github.com/gogo/protobuf/types"

	v039dispensation "github.com/Sifchain/sifnode/x/dispensation/legacy/v39"
	"github.com/Sifchain/sifnode/x/dispensation/types"
)

func Migrate(state v039dispensation.GenesisState) *types.GenesisState {

	records := make([]*types.DistributionRecord, 0)
	for _, d := range state.DistributionRecords {
		records = append(records, &types.DistributionRecord{
			DistributionStatus:          migrateDistributionRecordStatus(d.DistributionStatus),
			DistributionType:            migrateDistributionType(d.DistributionType),
			DistributionName:            d.DistributionName,
			RecipientAddress:            d.RecipientAddress.String(),
			Coins:                       d.Coins,
			DistributionStartHeight:     d.DistributionStartHeight,
			DistributionCompletedHeight: d.DistributionCompletedHeight,
			AuthorizedRunner:            d.AuthorizedRunner.String(),
		})
	}

	distributions := make([]*types.Distribution, 0)
	for _, d := range state.Distributions {
		distributions = append(distributions, &types.Distribution{
			DistributionType: migrateDistributionType(d.DistributionType),
			DistributionName: d.DistributionName,
			// Note: Distribution.Runner is only defined and stored in 39, it was not defined in 42
			Runner: d.Runner.String(),
		})
	}

	claims := make([]*types.UserClaim, 0)
	for _, c := range state.Claims {
		t, err := gogotypes.TimestampProto(c.UserClaimTime)
		if err != nil {
			panic(err)
		}
		claims = append(claims, &types.UserClaim{
			UserAddress:   c.UserAddress.String(),
			UserClaimType: migrateDistributionType(c.UserClaimType),
			UserClaimTime: t,
		})
	}

	return &types.GenesisState{
		DistributionRecords: &types.DistributionRecords{DistributionRecords: records},
		Distributions:       &types.Distributions{Distributions: distributions},
		Claims:              &types.UserClaims{UserClaims: claims},
	}
}

func migrateDistributionRecordStatus(status v039dispensation.DistributionStatus) types.DistributionStatus {
	// Note: Failed status does not exist on 39
	if status == v039dispensation.Pending {
		return types.DistributionStatus_DISTRIBUTION_STATUS_PENDING
	} else if status == v039dispensation.Completed {
		return types.DistributionStatus_DISTRIBUTION_STATUS_COMPLETED
	}

	return types.DistributionStatus_DISTRIBUTION_STATUS_UNSPECIFIED
}

func migrateDistributionType(t v039dispensation.DistributionType) types.DistributionType {
	if t == v039dispensation.Airdrop {
		return types.DistributionType_DISTRIBUTION_TYPE_AIRDROP
	} else if t == v039dispensation.LiquidityMining {
		return types.DistributionType_DISTRIBUTION_TYPE_LIQUIDITY_MINING
	} else if t == v039dispensation.ValidatorSubsidy {
		return types.DistributionType_DISTRIBUTION_TYPE_VALIDATOR_SUBSIDY
	}

	return types.DistributionType_DISTRIBUTION_TYPE_UNSPECIFIED
}
