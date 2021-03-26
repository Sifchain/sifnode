import { Ref } from "@vue/reactivity";

import {
  Amount,
  Asset,
  AssetAmount,
  LiquidityProvider,
  Pool,
} from "../entities";
import { calculateWithdrawal } from "../entities/formulae";
import { format } from "../utils/format";
import { PoolState } from "./addLiquidityCalculator";

export function useRemoveLiquidityCalculator(input: {
  externalAssetSymbol: Ref<string | null>;
  nativeAssetSymbol: Ref<string | null>;
  wBasisPoints: Ref<string | null>;
  asymmetry: Ref<string | null>;
  poolFinder: (a: Asset | string, b: Asset | string) => Ref<Pool> | null;
  liquidityProvider: Ref<LiquidityProvider | null>;
  sifAddress: Ref<string>;
}) {
  // this function needs to be refactored so
  const externalAsset = (() => {
    if (!input.externalAssetSymbol.value) return null;
    return Asset(input.externalAssetSymbol.value);
  })();

  const nativeAsset = (() => {
    if (!input.nativeAssetSymbol.value) return null;
    return Asset(input.nativeAssetSymbol.value);
  })();

  const liquidityPool = (() => {
    if (!nativeAsset || !externalAsset) return null;

    // Find pool from poolFinder
    const pool = input.poolFinder(externalAsset, nativeAsset);
    return pool?.value ?? null;
  })();

  const poolUnits = (() => {
    if (!liquidityPool) return null;
    return liquidityPool.poolUnits;
  })();

  const wBasisPoints = (() => {
    if (!input.wBasisPoints.value) return null;
    return Amount(input.wBasisPoints.value);
  })();

  const asymmetry = (() => {
    if (!input.asymmetry.value) return null;
    return Amount(input.asymmetry.value);
  })();

  const nativeAssetBalance = (() => {
    if (!liquidityPool) return null;
    return (
      liquidityPool.amounts.find(
        (a) => a.symbol === input.nativeAssetSymbol.value,
      ) ?? null
    );
  })();

  const externalAssetBalance = (() => {
    if (!liquidityPool) return null;
    return (
      liquidityPool.amounts.find(
        (a) => a.symbol === input.externalAssetSymbol.value,
      ) ?? null
    );
  })();

  const lpUnits = (() => {
    if (!input.liquidityProvider.value) return null;

    return input.liquidityProvider.value.units;
  })();

  const hasLiquidity = (() => {
    if (!lpUnits) return false;
    return lpUnits.greaterThan("0");
  })();

  const withdrawalAmounts = (() => {
    if (
      !poolUnits ||
      !nativeAssetBalance ||
      !externalAssetBalance ||
      !lpUnits ||
      !wBasisPoints ||
      !asymmetry ||
      !externalAsset ||
      !nativeAsset
    )
      return null;

    const {
      withdrawExternalAssetAmount,
      withdrawNativeAssetAmount,
    } = calculateWithdrawal({
      poolUnits,
      nativeAssetBalance,
      externalAssetBalance,
      lpUnits,
      wBasisPoints,
      asymmetry: asymmetry,
    });

    return {
      hasLiquidity,
      withdrawExternalAssetAmount: AssetAmount(
        externalAsset,
        withdrawExternalAssetAmount,
      ),
      withdrawNativeAssetAmount: AssetAmount(
        nativeAsset,
        withdrawNativeAssetAmount,
      ),
    };
  })();

  const state = (() => {
    if (!input.externalAssetSymbol.value || !input.nativeAssetSymbol.value)
      return PoolState.SELECT_TOKENS;

    if (!wBasisPoints?.greaterThan("0")) return PoolState.ZERO_AMOUNTS;

    if (!hasLiquidity) return PoolState.NO_LIQUIDITY;
    if (!lpUnits) {
      return PoolState.INSUFFICIENT_FUNDS;
    }

    return PoolState.VALID_INPUT;
  })();

  const withdrawExternalAssetAmountMessage = (() => {
    if (!withdrawalAmounts) return "";
    const assetAmount = withdrawalAmounts?.withdrawExternalAssetAmount;
    return format(assetAmount.amount, assetAmount.asset, {
      mantissa: 6,
    });
  })();

  const withdrawNativeAssetAmountMessage = (() => {
    if (!withdrawalAmounts) return "";
    const assetAmount = withdrawalAmounts?.withdrawNativeAssetAmount;
    return format(assetAmount.amount, assetAmount.asset, {
      mantissa: 6,
    });
  })();

  return {
    withdrawExternalAssetAmount: withdrawExternalAssetAmountMessage,
    withdrawNativeAssetAmount: withdrawNativeAssetAmountMessage,
    state,
  };
}
