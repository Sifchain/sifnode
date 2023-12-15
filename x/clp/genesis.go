package clp

import (
	"fmt"
	"math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/query"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/Sifchain/sifnode/x/clp/keeper"
	"github.com/Sifchain/sifnode/x/clp/types"
)

func InitGenesis(ctx sdk.Context, k keeper.Keeper, data types.GenesisState) (res []abci.ValidatorUpdate) {
	// Set initial CLP parameters
	k.SetParams(ctx, data.Params)

	// Set initial address whitelist
	if data.AddressWhitelist == nil || len(data.AddressWhitelist) == 0 {
		panic("AddressWhiteList must be set.")
	}
	wl := make([]sdk.AccAddress, len(data.AddressWhitelist))
	for i, entry := range data.AddressWhitelist {
		wlAddress, err := sdk.AccAddressFromBech32(entry)
		if err != nil {
			panic(err)
		}
		wl[i] = wlAddress
	}
	k.SetClpWhiteList(ctx, wl)

	// Set all the pools
	for _, pool := range data.PoolList {
		err := k.SetPool(ctx, pool)
		if err != nil {
			panic(fmt.Sprintf("Pool could not be set : %s", pool.String()))
		}
	}

	// Set all the liquidity providers
	for _, lp := range data.LiquidityProviders {
		k.SetLiquidityProvider(ctx, lp)
	}

	// Set all the rewardsBucket
	for _, elem := range data.RewardsBucketList {
		k.SetRewardsBucket(ctx, elem)
	}

	// Set initial reward states
	k.SetRewardParams(ctx, &data.RewardParams)

	// Set initial pmtp states
	k.SetPmtpRateParams(ctx, data.PmtpRateParams)
	k.SetPmtpEpoch(ctx, data.PmtpEpoch)
	k.SetPmtpParams(ctx, &data.PmtpParams)
	k.SetPmtpRateParams(ctx, data.PmtpRateParams)

	// Set initial liquidity protection states
	k.SetLiquidityProtectionParams(ctx, &data.LiquidityProtectionParams)
	k.SetLiquidityProtectionRateParams(ctx, data.LiquidityProtectionRateParams)

	// Set initial swap fee states
	k.SetSwapFeeParams(ctx, &data.SwapFeeParams)

	// Set initial provider distribution states
	k.SetProviderDistributionParams(ctx, &data.ProviderDistributionParams)

	return []abci.ValidatorUpdate{}
}

func ExportGenesis(ctx sdk.Context, keeper keeper.Keeper) types.GenesisState {
	// Retrieve CLP parameters
	params := keeper.GetParams(ctx)

	// Retrieve all the pools
	var poolList []*types.Pool
	poolList, _, _ = keeper.GetPoolsPaginated(ctx, &query.PageRequest{
		Limit: uint64(math.MaxUint64),
	})

	// Retrieve all the liquidity providers
	liquidityProviders, _, _ := keeper.GetAllLiquidityProvidersPaginated(ctx, &query.PageRequest{
		Limit: uint64(math.MaxUint64),
	})

	// Retrieve all the whitelist addresses
	whiteList := keeper.GetClpWhiteList(ctx)
	wl := make([]string, len(whiteList))
	for i, entry := range whiteList {
		wl[i] = entry.String()
	}

	// Retrieve all the rewardsBucket
	rewardsBucketList := keeper.GetAllRewardsBucket(ctx)

	// Retrieve all the reward states
	rewardParams := keeper.GetRewardsParams(ctx)
	if rewardParams == nil {
		rewardParams = types.GetDefaultRewardParams()
	}

	// Retrieve all the pmtp states
	pmtpParams := keeper.GetPmtpParams(ctx)
	if pmtpParams == nil {
		pmtpParams = types.GetDefaultPmtpParams()
	}
	pmtpEpoch := keeper.GetPmtpEpoch(ctx)
	pmtpRateParams := keeper.GetPmtpRateParams(ctx)

	// Retrieve all the liquidity protection states
	liquidityProtectionParams := keeper.GetLiquidityProtectionParams(ctx)
	if liquidityProtectionParams == nil {
		liquidityProtectionParams = types.GetDefaultLiquidityProtectionParams()
	}
	liquidityProtectionRateParams := keeper.GetLiquidityProtectionRateParams(ctx)

	// Retrieve all the swap fee states
	swapFeeParams := keeper.GetSwapFeeParams(ctx)

	// Retrieve all the provider distribution states
	providerDistributionParams := keeper.GetProviderDistributionParams(ctx)
	if providerDistributionParams == nil {
		providerDistributionParams = types.GetDefaultProviderDistributionParams()
	}

	return types.GenesisState{
		Params:                        params,
		AddressWhitelist:              wl,
		PoolList:                      poolList,
		LiquidityProviders:            liquidityProviders,
		RewardsBucketList:             rewardsBucketList,
		RewardParams:                  *rewardParams,
		PmtpParams:                    *pmtpParams,
		PmtpEpoch:                     pmtpEpoch,
		PmtpRateParams:                pmtpRateParams,
		LiquidityProtectionParams:     *liquidityProtectionParams,
		LiquidityProtectionRateParams: liquidityProtectionRateParams,
		SwapFeeParams:                 swapFeeParams,
		ProviderDistributionParams:    *providerDistributionParams,
	}
}

// ValidateGenesis validates the clp genesis parameters
func ValidateGenesis(data types.GenesisState) error {
	if err := data.Params.Validate(); err != nil {
		return sdkerrors.Wrap(types.ErrInvalid, fmt.Sprintf("clp: params are invalid : %s \n %s", err.Error(), data.Params.String()))
	}
	for _, pool := range data.PoolList {
		if !pool.Validate() {
			return sdkerrors.Wrap(types.ErrInvalid, fmt.Sprintf("clp: pool is invalid : %s", pool.String()))
		}
	}
	for _, lp := range data.LiquidityProviders {
		if !lp.Validate() {
			return sdkerrors.Wrap(types.ErrInvalid, fmt.Sprintf("clp: liquidityProvider is invalid : %s", lp.String()))
		}
	}
	return nil
}
