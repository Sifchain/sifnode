package chain

import (
	clptypes "github.com/Sifchain/sifnode/x/clp/types"
	keeper "github.com/Sifchain/sifnode/x/margin/keeper"
	margintypes "github.com/Sifchain/sifnode/x/margin/types"
	"github.com/cosmos/cosmos-sdk/client"
)

/* Query margin params from chain. */

func QueryMarginParams(clientCtx client.Context) (*margintypes.Params, error) {
	marginTypes := &margintypes.Params{
		"LeverageMax":          keeper.GetLeverageParam(clientCtx),
		"InterestRateMax":      keeper.GetInterestRateMax(clientCtx),
		"InterestRateMin":      keeper.GetInterestRateMin(clientCtx),
		"InterestRateIncrease": keeper.GetInterestRateIncrease(clientCtx),
		"InterestRateDecrease": keeper.GetInterestRateDecrease(clientCtx),
		"HealthGainFactor":     keeper.GetHealthGainFactor(clientCtx),
		"EpochLength":          keeper.GetEpochLength,
	}
	return marginTypes, nil
}

/* Query pool data from chain. */

func QueryPools(clientCtx client.Context) ([]*clptypes.Pool, error) {
	pools := keeper.GetEnabledPools(clientCtx)
	return pools, nil
}

/* Returns events from block_results?height=height */

func BlockEvents(clientCtx client.Context, height int64) ([]*Event, error) {
	return nil, nil
}
