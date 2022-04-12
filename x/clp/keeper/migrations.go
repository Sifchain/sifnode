package keeper

import (
	"github.com/Sifchain/sifnode/x/clp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Migrator struct {
	keeper Keeper
}

func NewMigrator(keeper Keeper) Migrator {
	return Migrator{keeper: keeper}
}

func (m Migrator) MigrateToVer2(ctx sdk.Context) error {
	// Initiate Rewards
	m.keeper.SetRewardParams(ctx, types.GetDefaultRewardParams())
	// Initiate PMTP
	m.keeper.SetPmtpRateParams(ctx, types.PmtpRateParams{
		PmtpPeriodBlockRate:    sdk.ZeroDec(),
		PmtpCurrentRunningRate: sdk.ZeroDec(),
		PmtpInterPolicyRate:    sdk.ZeroDec(),
	})
	m.keeper.SetPmtpEpoch(ctx, types.PmtpEpoch{
		EpochCounter: 0,
		BlockCounter: 0,
	})
	m.keeper.SetPmtpParams(ctx, types.GetDefaultPmtpParams())
	m.keeper.SetPmtpInterPolicyRate(ctx, sdk.NewDec(0))
	pools := m.keeper.GetPools(ctx)
	for _, pool := range pools {
		spe := sdk.ZeroDec()
		spn := sdk.ZeroDec()
		rpnd := sdk.ZeroUint()
		pool.SwapPriceExternal = &spe
		pool.SwapPriceNative = &spn
		pool.RewardPeriodNativeDistributed = rpnd
		err := m.keeper.SetPool(ctx, pool)
		if err != nil {
			panic(err)
		}
	}
	return nil
}

func (m Migrator) MigrateToVer3(ctx sdk.Context) error {
	// Initiate Rewards
	m.keeper.SetRewardParams(ctx, types.GetDefaultRewardParams())
	// Initiate PMTP
	m.keeper.SetPmtpRateParams(ctx, types.PmtpRateParams{
		PmtpPeriodBlockRate:    sdk.ZeroDec(),
		PmtpCurrentRunningRate: sdk.ZeroDec(),
		PmtpInterPolicyRate:    sdk.ZeroDec(),
	})
	m.keeper.SetPmtpEpoch(ctx, types.PmtpEpoch{
		EpochCounter: 0,
		BlockCounter: 0,
	})
	m.keeper.SetPmtpParams(ctx, types.GetDefaultPmtpParams())
	m.keeper.SetPmtpInterPolicyRate(ctx, sdk.NewDec(0))
	pools := m.keeper.GetPools(ctx)
	for _, pool := range pools {
		spe := sdk.ZeroDec()
		spn := sdk.ZeroDec()
		rpnd := sdk.ZeroUint()
		pool.SwapPriceExternal = &spe
		pool.SwapPriceNative = &spn
		pool.RewardPeriodNativeDistributed = rpnd
		err := m.keeper.SetPool(ctx, pool)
		if err != nil {
			panic(err)
		}
	}
	return nil
}

func (m Migrator) MigrateToVer4(ctx sdk.Context) error {
	// Initiate Rewards
	m.keeper.SetRewardParams(ctx, types.GetDefaultRewardParams())
	// Initiate PMTP
	m.keeper.SetPmtpRateParams(ctx, types.PmtpRateParams{
		PmtpPeriodBlockRate:    sdk.ZeroDec(),
		PmtpCurrentRunningRate: sdk.ZeroDec(),
		PmtpInterPolicyRate:    sdk.ZeroDec(),
	})
	m.keeper.SetPmtpEpoch(ctx, types.PmtpEpoch{
		EpochCounter: 0,
		BlockCounter: 0,
	})
	m.keeper.SetPmtpParams(ctx, types.GetDefaultPmtpParams())
	m.keeper.SetPmtpInterPolicyRate(ctx, sdk.NewDec(0))
	pools := m.keeper.GetPools(ctx)
	for _, pool := range pools {
		spe := sdk.ZeroDec()
		spn := sdk.ZeroDec()
		rpnd := sdk.ZeroUint()
		pool.SwapPriceExternal = &spe
		pool.SwapPriceNative = &spn
		pool.RewardPeriodNativeDistributed = rpnd
		err := m.keeper.SetPool(ctx, pool)
		if err != nil {
			panic(err)
		}
	}
	return nil
}
