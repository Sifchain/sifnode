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

	// compute swap prices for each pool
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
		pool.SwapPriceExternal = &spe
		pool.SwapPriceNative = &spn
		err := m.keeper.SetPool(ctx, pool)
		if err != nil {
			panic(err)
		}
	}
	// SET ADMIN ACCOUNT
	// m.keeper.SetAdminAccount()

	return nil
}
