import { computed, ComputedRef, effect, ref, Ref } from "@vue/reactivity";

import { Asset, AssetAmount, LiquidityProvider, Pool } from "../entities";
import { calculateWithdrawal } from "../entities/formulae";
import { Fraction, IFraction } from "../entities/fraction/Fraction";
import { PoolState } from "./addLiquidityCalculator";
import { buildAsset } from "./utils";

export function useRemoveLiquidityCalculator(input: {
  externalAssetSymbol: Ref<string | null>;
  nativeAssetSymbol: Ref<string | null>;
  wBasisPoints: Ref<string | null>;
  asymmetry: Ref<string | null>;
  marketPairFinder: (a: Asset | string, b: Asset | string) => Pool | null;
  liquidityProvider: Ref<LiquidityProvider | null>;
  sifAddress: Ref<string>;
}) {
  const externalAsset = computed(() => {
    if (!input.externalAssetSymbol.value) return null;
    return buildAsset(input.externalAssetSymbol.value);
  });

  const nativeAsset = computed(() => {
    if (!input.nativeAssetSymbol.value) return null;
    return buildAsset(input.nativeAssetSymbol.value);
  });

  const liquidityPool = computed(() => {
    if (!nativeAsset.value || !externalAsset.value) return null;

    // Find pool from marketPairFinder
    return input.marketPairFinder(nativeAsset.value, externalAsset.value);
  });

  const poolUnits = computed(() => {
    if (!liquidityPool.value) return null;
    return liquidityPool.value.poolUnits;
  });

  const wBasisPoints = computed(() => {
    if (!input.wBasisPoints.value) return null;
    return new Fraction(input.wBasisPoints.value);
  });

  const asymmetry = computed(() => {
    if (!input.asymmetry.value) return null;
    return new Fraction(input.asymmetry.value);
  });

  const nativeAssetBalance = computed(() => {
    if (!liquidityPool.value) return null;
    return (
      liquidityPool.value.amounts.find(
        (a) => a.asset.symbol === input.nativeAssetSymbol.value
      ) ?? null
    );
  });

  const externalAssetBalance = computed(() => {
    if (!liquidityPool.value) return null;
    return (
      liquidityPool.value.amounts.find(
        (a) => a.asset.symbol === input.externalAssetSymbol.value
      ) ?? null
    );
  });

  const lpUnits = computed(() => {
    if (!input.liquidityProvider.value) return null;

    return input.liquidityProvider.value.units as IFraction;
  });

  const hasLiquidity = computed(() => {
    if (!lpUnits.value) return false;
    return lpUnits.value.greaterThan("0");
  });

  const withdrawalAmounts = computed(() => {
    if (
      !poolUnits.value ||
      !nativeAssetBalance.value ||
      !externalAssetBalance.value ||
      !lpUnits.value ||
      !wBasisPoints.value ||
      !asymmetry.value ||
      !externalAsset.value ||
      !nativeAsset.value
    )
      return null;

    const inputs = {
      poolUnits: poolUnits.value,
      nativeAssetBalance: nativeAssetBalance.value,
      externalAssetBalance: externalAssetBalance.value,
      lpUnits: lpUnits.value,
      wBasisPoints: wBasisPoints.value,
      asymmetry: asymmetry.value,
    };
    const {
      withdrawExternalAssetAmount,
      withdrawNativeAssetAmount,
    } = calculateWithdrawal(inputs);

    return {
      hasLiquidity,
      withdrawExternalAssetAmount: AssetAmount(
        externalAsset.value,
        withdrawExternalAssetAmount
      ),
      withdrawNativeAssetAmount: AssetAmount(
        nativeAsset.value,
        withdrawNativeAssetAmount
      ),
    };
  });

  const state = computed(() => {
    if (!input.externalAssetSymbol.value || !input.nativeAssetSymbol.value)
      return PoolState.SELECT_TOKENS;

    if (!wBasisPoints.value?.greaterThan("0")) return PoolState.ZERO_AMOUNTS;

    if (!hasLiquidity.value) return PoolState.NO_LIQUIDITY;
    if (!lpUnits.value) {
      return PoolState.INSUFFICIENT_FUNDS;
    }

    return PoolState.VALID_INPUT;
  });

  const withdrawExternalAssetAmountMessage = computed(() => {
    return (
      withdrawalAmounts.value?.withdrawExternalAssetAmount.toFormatted({
        decimals: 1,
      }) || ""
    );
  });

  const withdrawNativeAssetAmountMessage = computed(() => {
    return (
      withdrawalAmounts.value?.withdrawNativeAssetAmount.toFormatted({
        decimals: 1,
      }) || ""
    );
  });

  return {
    withdrawExternalAssetAmount: withdrawExternalAssetAmountMessage,
    withdrawNativeAssetAmount: withdrawNativeAssetAmountMessage,
    state,
  };
}
