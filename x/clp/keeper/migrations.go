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
	type LPS struct {
		address string
		units   sdk.Uint
	}
	lps := []LPS{
		{"sif1smknl4uf89ef84kg020ff7ask7l0sxz3s93gva", sdk.NewUintFromString("49661617604396299845632")},
		{"sif1gy2ne7m62uer4h5s4e7xlfq7aeem5zpwx6nu9q", sdk.NewUintFromString("57436791960686469054464")},
		{"sif1uk77p6he26undp9wjjav6ygtu53kswl60cd5va", sdk.NewUintFromString("54301402049261093257216")},
		{"sif1hspkfnzexvn4drk9dlfpg8n0ppw8sxsl00t65a", sdk.NewUintFromString("55396192200751201648640")},
		{"sif1y2rfrgh374gd40gj0yusjasw2paaysahu9qk5j", sdk.NewUintFromString("138696029054663414251520")},
		{"sif1ra9563z5tn2lmqhydt2atrgzftk2d7umyr4vqw", sdk.NewUintFromString("139123862445111950442496")}}

	symbol := "cusdc"

	for _, lp := range lps {
		l, err := m.keeper.GetLiquidityProvider(ctx, symbol, lp.address)
		if err != nil {
			ctx.Logger().Error("error looking up LP for migration", "lp", lp.address)
			continue
		}
		pool, err := m.keeper.GetPool(ctx, symbol)
		if err != nil {
			panic(err)
		}

		l.LiquidityProviderUnits = l.LiquidityProviderUnits.Sub(lp.units)
		m.keeper.SetLiquidityProvider(ctx, &l)

		pool.PoolUnits = pool.PoolUnits.Sub(lp.units)
		// nolint:errcheck
		m.keeper.SetPool(ctx, &pool)
	}

	return nil
}

func (m Migrator) MigrateToVer4(ctx sdk.Context) error {
	m.keeper.SetParams(ctx, types.DefaultParams())

	return nil
}

func (m Migrator) MigrateToVer5(ctx sdk.Context) error {
	// set rewards params
	m.keeper.SetRewardParams(ctx, types.GetDefaultRewardParams())

	// loop over all the pools and set the reward amount external field to zero
	pools := m.keeper.GetPools(ctx)
	for _, pool := range pools {
		pool.RewardAmountExternal = sdk.ZeroUint()
		// nolint:errcheck
		m.keeper.SetPool(ctx, pool)
	}

	// loop over all the LPs and set their last_updated_block to the current block
	// this will ensure that they get their rewards
	// and initialize the reward amount array
	all, err := m.keeper.GetAllLiquidityProviders(ctx)
	if err != nil {
		return err
	}
	for _, lp := range all {
		lp.LastUpdatedBlock = ctx.BlockHeight()
		lp.RewardAmount = sdk.NewCoins()
		m.keeper.SetLiquidityProvider(ctx, lp)
	}

	return nil
}
