import { Ref } from "@vue/reactivity";

import { Asset, AssetAmount, LiquidityProvider, Pool } from "../entities";
import { calculateWithdrawal } from "../entities/formulae";
import { Fraction, IFraction } from "../entities";
import { PoolState } from "./addLiquidityCalculator";
import { buildAsset } from "./utils";

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
    return buildAsset(input.externalAssetSymbol.value);
  })();

  const nativeAsset = (() => {
    if (!input.nativeAssetSymbol.value) return null;
    return buildAsset(input.nativeAssetSymbol.value);
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
    return new Fraction(input.wBasisPoints.value);
  })();

  const asymmetry = (() => {
    if (!input.asymmetry.value) return null;
    return new Fraction(input.asymmetry.value);
  })();

  const nativeAssetBalance = (() => {
    if (!liquidityPool) return null;
    return (
      liquidityPool.amounts.find(
        (a) => a.asset.symbol === input.nativeAssetSymbol.value,
      ) ?? null
    );
  })();

  const externalAssetBalance = (() => {
    if (!liquidityPool) return null;
    return (
      liquidityPool.amounts.find(
        (a) => a.asset.symbol === input.externalAssetSymbol.value,
      ) ?? null
    );
  })();

  const lpUnits = (() => {
    if (!input.liquidityProvider.value) return null;

    return input.liquidityProvider.value.units as IFraction;
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
    return (
      withdrawalAmounts?.withdrawExternalAssetAmount.toFormatted({
        decimals: 6,
        symbol: false,
      }) || ""
    );
  })();

  const withdrawNativeAssetAmountMessage = (() => {
    return (
      withdrawalAmounts?.withdrawNativeAssetAmount.toFormatted({
        decimals: 6,
        symbol: false,
      }) || ""
    );
  })();

  return {
    withdrawExternalAssetAmount: withdrawExternalAssetAmountMessage,
    withdrawNativeAssetAmount: withdrawNativeAssetAmountMessage,
    state,
  };
}
