package keeper

import (
	"github.com/Sifchain/sifnode/x/margin/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Migrator struct {
	keeper Keeper
}

func NewMigrator(keeper Keeper) Migrator {
	return Migrator{keeper: keeper}
}

func (m Migrator) MigrateToVer2(ctx sdk.Context) error {
	m.keeper.SetParams(ctx, &types.Params{
		LeverageMax:              sdk.NewDec(2),
		HealthGainFactor:         sdk.NewDec(1),
		InterestRateMin:          sdk.NewDecWithPrec(5, 3),
		InterestRateMax:          sdk.NewDec(3),
		InterestRateDecrease:     sdk.NewDecWithPrec(1, 1),
		InterestRateIncrease:     sdk.NewDecWithPrec(1, 1),
		ForceCloseThreshold:      sdk.NewDecWithPrec(1, 1),
		ForceCloseFundPercentage: sdk.NewDecWithPrec(1, 1),
		InsuranceFundAddress:     "",
		PoolOpenThreshold:        sdk.NewDecWithPrec(1, 1),
		RemovalQueueThreshold:    sdk.NewDecWithPrec(1, 1),
		EpochLength:              1,
		MaxOpenPositions:         10000,
		Pools:                    []string{},
		SqModifier:               sdk.MustNewDecFromStr("10000000000000000000000000"),
		SafetyFactor:             sdk.MustNewDecFromStr("1.05"),
	})

	return nil
}
