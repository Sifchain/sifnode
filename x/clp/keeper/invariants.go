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
