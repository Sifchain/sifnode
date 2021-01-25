import { computed, Ref } from "@vue/reactivity";

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
  poolFinder: (a: Asset | string, b: Asset | string) => Ref<Pool> | null;
  liquidityProvider: Ref<LiquidityProvider | null>;
  sifAddress: Ref<string>;
}) {

  if (!input.externalAssetSymbol.value) return null;
  if (!input.nativeAssetSymbol.value) return null;
  if (!input.wBasisPoints.value) return null;
  if (!input.asymmetry.value) return null;
  if (!input.liquidityProvider.value) return null;

  const externalAsset = buildAsset(input.externalAssetSymbol.value);
  const nativeAsset = buildAsset(input.nativeAssetSymbol.value);

  if (!nativeAsset || !externalAsset) return null;

  // Find pool from poolFinder
  const pool = input.poolFinder(externalAsset, nativeAsset);
  const liquidityPool = pool?.value ?? null;
  
  if (!liquidityPool) return null;

  const poolUnits = liquidityPool.poolUnits;
  const wBasisPoints = new Fraction(input.wBasisPoints.value);
  const asymmetry = new Fraction(input.asymmetry.value);

  const nativeAssetBalance = (
    liquidityPool.amounts.find(
      (a) => a.asset.symbol === input.nativeAssetSymbol.value
    ) ?? null
  )

  const externalAssetBalance = (
    liquidityPool.amounts.find(
      (a) => a.asset.symbol === input.externalAssetSymbol.value
    ) ?? null
  )

  const lpUnits = input.liquidityProvider.value.units as IFraction;
  if (!lpUnits) return null;

  const hasLiquidity = lpUnits.greaterThan("0");

  const withdrawalAmounts = (() => {
    if (
      !poolUnits ||
      !nativeAssetBalance ||
      !externalAssetBalance ||
      !lpUnits ||
      !wBasisPoints ||
      !asymmetry ||
      !externalAsset ||
      !nativeAsset  || 
      !hasLiquidity
    ) {
      return null
    }

    const inputs = {
      poolUnits: poolUnits,
      nativeAssetBalance: nativeAssetBalance,
      externalAssetBalance: externalAssetBalance,
      lpUnits: lpUnits,
      wBasisPoints: wBasisPoints,
      asymmetry: asymmetry,
    };
    const {
      withdrawExternalAssetAmount,
      withdrawNativeAssetAmount,
    } = calculateWithdrawal(inputs);

    return {
      hasLiquidity,
      withdrawExternalAssetAmount: AssetAmount(
        externalAsset,
        withdrawExternalAssetAmount
      ),
      withdrawNativeAssetAmount: AssetAmount(
        nativeAsset,
        withdrawNativeAssetAmount
      ),
    };
  })()

  const state = (() => {
    if (!input.externalAssetSymbol.value || !input.nativeAssetSymbol.value)
      return PoolState.SELECT_TOKENS;

    if (!wBasisPoints?.greaterThan("0")) return PoolState.ZERO_AMOUNTS;

    if (!hasLiquidity) return PoolState.NO_LIQUIDITY;
    if (!lpUnits) {
      return PoolState.INSUFFICIENT_FUNDS;
    }

    return PoolState.VALID_INPUT;
  })()

  const withdrawExternalAssetAmountMessage = (() => {
    return (
      withdrawalAmounts?.withdrawExternalAssetAmount.toFormatted({
        decimals: 18,
      }) || ""
    );
  })()

  const withdrawNativeAssetAmountMessage = (() => {
    return (
      withdrawalAmounts?.withdrawNativeAssetAmount.toFormatted({
        decimals: 18,
      }) || ""
    );
  })()

  return {
    withdrawExternalAssetAmount: withdrawExternalAssetAmountMessage,
    withdrawNativeAssetAmount: withdrawNativeAssetAmountMessage,
    state,
  };
}
