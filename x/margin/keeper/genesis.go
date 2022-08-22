//go:build FEATURE_TOGGLE_MARGIN_CLI_ALPHA
// +build FEATURE_TOGGLE_MARGIN_CLI_ALPHA

package keeper

import (
	"github.com/Sifchain/sifnode/x/margin/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func (k Keeper) InitGenesis(ctx sdk.Context, data types.GenesisState) []abci.ValidatorUpdate {
	k.SetParams(ctx, data.Params)

	return []abci.ValidatorUpdate{}
}

func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	return &types.GenesisState{
		Params: &types.Params{
			LeverageMax:                                    k.GetMaxLeverageParam(ctx),
			InterestRateMax:                                k.GetInterestRateMax(ctx),
			InterestRateMin:                                k.GetInterestRateMin(ctx),
			InterestRateIncrease:                           k.GetInterestRateIncrease(ctx),
			InterestRateDecrease:                           k.GetInterestRateDecrease(ctx),
			HealthGainFactor:                               k.GetHealthGainFactor(ctx),
			EpochLength:                                    k.GetEpochLength(ctx),
			ForceCloseThreshold:                            k.GetForceCloseThreshold(ctx),
			ForceCloseFundPercentage:                       k.GetForceCloseFundPercentage(ctx),
			ForceCloseInsuranceFundAddress:                 k.GetForceCloseInsuranceFundAddress(ctx).String(),
			IncrementalInterestPaymentFundPercentage:       k.GetIncrementalInterestPaymentFundPercentage(ctx),
			IncrementalInterestPaymentInsuranceFundAddress: k.GetIncrementalInterestPaymentInsuranceFundAddress(ctx).String(),
			PoolOpenThreshold:                              k.GetPoolOpenThreshold(ctx),
			RemovalQueueThreshold:                          k.GetRemovalQueueThreshold(ctx),
			MaxOpenPositions:                               k.GetMaxOpenPositions(ctx),
			Pools:                                          k.GetEnabledPools(ctx),
			SqModifier:                                     k.GetSqModifier(ctx),
			SafetyFactor:                                   k.GetSafetyFactor(ctx),
		},
	}
}
