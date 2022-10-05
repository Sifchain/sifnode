package keeper

import (
	"fmt"

	"github.com/Sifchain/sifnode/x/clp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func RegisterInvariants(registry sdk.InvariantRegistry, k Keeper) {
	// registry.RegisterRoute(types.ModuleName, "balance-module-account-check", k.BalanceModuleAccountCheck())
}

func (k Keeper) BalanceModuleAccountCheck() sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		// Get Rowan Balance from CLP Module
		clpModuleTotalNativeBalance := k.GetBankKeeper().GetBalance(ctx, types.GetCLPModuleAddress(), types.GetSettlementAsset().Symbol)
		clpModuleTotalNativeBalanceUint := sdk.NewUintFromString(clpModuleTotalNativeBalance.Amount.String())

		pools := k.GetPools(ctx)
		poolsTotalNativeBalanceUint := sdk.ZeroUint()
		poolsTotalNativeCustodyUint := sdk.ZeroUint()
		for _, pool := range pools {
			poolsTotalNativeBalanceUint = poolsTotalNativeBalanceUint.Add(pool.NativeAssetBalance)
			poolsTotalNativeCustodyUint = poolsTotalNativeCustodyUint.Add(pool.NativeCustody)

			clpModuleTotalExternalBalance := k.GetBankKeeper().GetBalance(ctx, types.GetCLPModuleAddress(), pool.ExternalAsset.Symbol)
			clpModuleTotalExternalBalanceUint := sdk.NewUintFromString(clpModuleTotalExternalBalance.Amount.String())

			ok := pool.ExternalAssetBalance.Add(pool.ExternalCustody).Equal(clpModuleTotalExternalBalanceUint)
			if !ok {
				return fmt.Sprintf("external balance mismatch in pool %s (module: %s != pool: %s)",
					pool.ExternalAsset.Symbol,
					clpModuleTotalExternalBalanceUint.String(),
					pool.ExternalAssetBalance.String()), true
			}
		}

		ok := poolsTotalNativeBalanceUint.Add(poolsTotalNativeCustodyUint).Equal(clpModuleTotalNativeBalanceUint)
		if !ok {
			return fmt.Sprintf("native balance mismatch across all pools (module: %s != pools: %s)",
				clpModuleTotalNativeBalanceUint.String(),
				poolsTotalNativeBalanceUint.String()), true
		}

		return "pool and module account balances match", false
	}
}

func (k Keeper) SingleExternalBalanceModuleAccountCheck(externalAsset string) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		pool, err := k.GetPool(ctx, externalAsset)
		if err == nil {
			clpModuleTotalExternalBalance := k.GetBankKeeper().GetBalance(ctx, types.GetCLPModuleAddress(), pool.ExternalAsset.Symbol)
			clpModuleTotalExternalBalanceUint := sdk.NewUintFromString(clpModuleTotalExternalBalance.Amount.String())

			ok := pool.ExternalAssetBalance.Add(pool.ExternalCustody).Equal(clpModuleTotalExternalBalanceUint)
			if !ok {
				return fmt.Sprintf("external balance mismatch in pool %s (module: %s != pool: %s)",
					pool.ExternalAsset.Symbol,
					clpModuleTotalExternalBalanceUint.String(),
					pool.ExternalAssetBalance.String()), true
			}
		}

		// Get Rowan Balance from CLP Module
		clpModuleTotalNativeBalance := k.GetBankKeeper().GetBalance(ctx, types.GetCLPModuleAddress(), types.GetSettlementAsset().Symbol)
		clpModuleTotalNativeBalanceUint := sdk.NewUintFromString(clpModuleTotalNativeBalance.Amount.String())

		pools := k.GetPools(ctx)
		poolsTotalNativeBalanceUint := sdk.ZeroUint()
		poolsTotalNativeCustodyUint := sdk.ZeroUint()
		for _, pool := range pools {
			poolsTotalNativeBalanceUint = poolsTotalNativeBalanceUint.Add(pool.NativeAssetBalance)
			poolsTotalNativeCustodyUint = poolsTotalNativeCustodyUint.Add(pool.NativeCustody)
		}

		ok := poolsTotalNativeBalanceUint.Add(poolsTotalNativeCustodyUint).Equal(clpModuleTotalNativeBalanceUint)
		if !ok {
			return fmt.Sprintf("native balance mismatch across all pools (module: %s != pools: %s)",
				clpModuleTotalNativeBalanceUint.String(),
				poolsTotalNativeBalanceUint.String()), true
		}

		return "pool and module account balances match", false
	}
}

func (k Keeper) UnitsCheck() sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		pools := k.GetPools(ctx)
		for _, pool := range pools {
			totalLPUnits := sdk.ZeroUint()

			lps, err := k.GetAllLiquidityProvidersForAsset(ctx, *pool.ExternalAsset)
			if err != nil {
				ctx.Logger().Error("error looking up LPs for pool during UnitsCheck",
					"pool", pool.ExternalAsset.Symbol,
				)
				continue
			}

			for _, lp := range lps {
				totalLPUnits = totalLPUnits.Add(lp.LiquidityProviderUnits)
			}

			ok := pool.PoolUnits.Equal(totalLPUnits)
			if !ok {
				return fmt.Sprintf("pool units vs total lp units mismatch in pool %s (pool: %s != lps: %s)",
					pool.ExternalAsset.Symbol,
					pool.PoolUnits.String(),
					totalLPUnits.String(),
				), true
			}
		}
		return "all pool units vs total lp units match", false
	}
}
