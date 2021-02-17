// TODO remove refs dependency and move to `actions/clp/calculateAddLiquidity`

import { computed, Ref } from "@vue/reactivity";
import {
  Asset,
  AssetAmount,
  IAssetAmount,
  LiquidityProvider,
  Pool,
} from "../entities";
import { Fraction } from "../entities";
import { useField } from "./useField";
import { useBalances } from "./utils";

export enum PoolState {
  SELECT_TOKENS,
  ZERO_AMOUNTS,
  INSUFFICIENT_FUNDS,
  VALID_INPUT,
  NO_LIQUIDITY,
}

export function usePoolCalculator(input: {
  tokenAAmount: Ref<string>;
  tokenASymbol: Ref<string | null>;
  tokenBAmount: Ref<string>;
  tokenBSymbol: Ref<string | null>;
  balances: Ref<IAssetAmount[]>;
  liquidityProvider: Ref<LiquidityProvider | null>;
  poolFinder: (a: Asset | string, b: Asset | string) => Ref<Pool> | null;
}) {
  const tokenAField = useField(input.tokenAAmount, input.tokenASymbol);
  const tokenBField = useField(input.tokenBAmount, input.tokenBSymbol);
  const balanceMap = useBalances(input.balances);

  const tokenABalance = computed(() => {
    return input.tokenASymbol.value
      ? balanceMap.value.get(input.tokenASymbol.value) ?? null
      : null;
  });

  const tokenBBalance = computed(() => {
    return input.tokenBSymbol.value
      ? balanceMap.value.get(input.tokenBSymbol.value) ?? null
      : null;
  });

  const fromBalanceOverdrawn = computed(() => {
    return !tokenABalance.value?.greaterThanOrEqual(
      tokenAField.fieldAmount.value || "0"
    );
  });

  const toBalanceOverdrawn = computed(
    () =>
      !tokenBBalance.value?.greaterThanOrEqual(
        tokenBField.fieldAmount.value || "0"
      )
  );

  const preExistingPool = computed(() => {
    if (!tokenAField.asset.value || !tokenBField.asset.value) return null;

    // Find pool from poolFinder
    const pool = input.poolFinder(
      tokenAField.asset.value.symbol,
      tokenBField.asset.value.symbol
    );

    return pool?.value || null;
  });

  const liquidityPool = computed(() => {
    if (preExistingPool.value) return preExistingPool.value;
    if (
      !tokenAField.fieldAmount.value ||
      !tokenBField.fieldAmount.value ||
      !tokenAField.asset.value ||
      !tokenBField.asset.value
    )
      return null;

    return Pool(
      AssetAmount(tokenAField.asset.value, "0"),
      AssetAmount(tokenBField.asset.value, "0")
    );
  });

  // pool units for this prospective transaction [total, newUnits]
  const provisionedPoolUnitsArray = computed(() => {
    if (
      !liquidityPool.value ||
      !tokenBField.fieldAmount.value ||
      !tokenAField.fieldAmount.value
    ) {
      return [new Fraction("0"), new Fraction("0")];
    }

    return liquidityPool.value.calculatePoolUnits(
      tokenBField.fieldAmount.value,
      tokenAField.fieldAmount.value
    );
  });

  // pool units from the perspective of the liquidity provider
  const liquidityProviderPoolUnitsArray = computed(() => {
    if (!provisionedPoolUnitsArray.value)
      return [new Fraction("0"), new Fraction("0")];

    const [totalPoolUnits, newUnits] = provisionedPoolUnitsArray.value;

    // if this user already has pool units include those in the newUnits
    const totalLiquidityProviderUnits = input.liquidityProvider.value
      ? input.liquidityProvider.value.units.add(newUnits)
      : newUnits;

    return [totalPoolUnits, totalLiquidityProviderUnits];
  });

  const totalPoolUnits = computed(() =>
    liquidityProviderPoolUnitsArray.value[0].toFixed(0)
  );

  const totalLiquidityProviderUnits = computed(() =>
    liquidityProviderPoolUnitsArray.value[1].toFixed(0)
  );

  const shareOfPool = computed(() => {
    if (!liquidityProviderPoolUnitsArray.value) return new Fraction("0");

    const [units, lpUnits] = liquidityProviderPoolUnitsArray.value;

    // shareOfPool should be 0 if units and lpUnits are zero
    if (units.equalTo("0") && lpUnits.equalTo("0")) return new Fraction("0");

    // if no units lp owns 100% of pool
    return units.equalTo("0") ? new Fraction("1") : lpUnits.divide(units);
  });

  const shareOfPoolPercent = computed(() => {
    if (shareOfPool.value.multiply("10000").lessThan("1")) return "< 0.01%";
    return `${shareOfPool.value.multiply("100").toFixed(2)}%`;
  });

  const poolAmounts = computed(() => {
    if (!preExistingPool.value || !tokenAField.asset.value) {
      return null;
    }
    if (!preExistingPool.value.contains(tokenAField.asset.value)) return null;
    const externalBalance = preExistingPool.value.getAmount(
      tokenAField.asset.value
    );
    const nativeBalance = preExistingPool.value.getAmount("rowan");
    return [nativeBalance, externalBalance];
  });

  // external_balance / native_balance
  const aPerBRatio = computed(() => {
    if (!poolAmounts.value) return null;
    const [native, external] = poolAmounts.value;
    return external.divide(native);
  });

  const aPerBRatioMessage = computed(() => {
    if (!aPerBRatio.value) {
      return "N/A";
    }

    return aPerBRatio.value.toFixed(8);
  });

  // native_balance / external_balance
  const bPerARatio = computed(() => {
    if (!poolAmounts.value) return null;
    const [native, external] = poolAmounts.value;
    return native.divide(external);
  });

  const bPerARatioMessage = computed(() => {
    if (!bPerARatio.value) {
      return "N/A";
    }

    return bPerARatio.value.toFixed(8);
  });

  // Price Impact and Pool Share:
  // (external_balance + external_added) / (native_balance + native_added)
  const aPerBRatioProjected = computed(() => {
    if (
      !poolAmounts.value ||
      !tokenAField.fieldAmount.value ||
      !tokenBField.fieldAmount.value
    )
      return null;

    const [native, external] = poolAmounts.value;
    const externalAdded = tokenAField.fieldAmount.value;
    const nativeAdded = tokenBField.fieldAmount.value;

    return external.add(externalAdded).divide(native.add(nativeAdded));
  });

  const aPerBRatioProjectedMessage = computed(() => {
    if (!aPerBRatioProjected.value) {
      return "N/A";
    }

    return aPerBRatioProjected.value.toFixed(8);
  });

  // Price Impact and Pool Share:
  // (native_balance + native_added)/(external_balance + external_added)
  const bPerARatioProjected = computed(() => {
    if (
      !poolAmounts.value ||
      !tokenAField.fieldAmount.value ||
      !tokenBField.fieldAmount.value
    )
      return null;

    const [native, external] = poolAmounts.value;
    const externalAdded = tokenAField.fieldAmount.value;
    const nativeAdded = tokenBField.fieldAmount.value;
    return native.add(nativeAdded).divide(external.add(externalAdded));
  });

  const bPerARatioProjectedMessage = computed(() => {
    if (!bPerARatioProjected.value) {
      return "N/A";
    }

    return bPerARatioProjected.value.toFixed(8);
  });

  const state = computed(() => {
    // Select Tokens
    const aSymbolNotSelected = !input.tokenASymbol.value;
    const bSymbolNotSelected = !input.tokenBSymbol.value;
    if (aSymbolNotSelected || bSymbolNotSelected)
      return PoolState.SELECT_TOKENS;

    // Zero amounts
    const aAmount = tokenAField.fieldAmount.value;
    const bAmount = tokenBField.fieldAmount.value;
    const aAmountIsZeroOrFalsy = !aAmount || aAmount.equalTo("0");
    const bAmountIsZeroOrFalsy = !bAmount || bAmount.equalTo("0");
    const noPreexistingPool = !preExistingPool.value;
    if (noPreexistingPool || (bAmountIsZeroOrFalsy && aAmountIsZeroOrFalsy))
      return PoolState.ZERO_AMOUNTS;

    // Insufficient Funds
    if (fromBalanceOverdrawn.value || toBalanceOverdrawn.value)
      return PoolState.INSUFFICIENT_FUNDS;

    // Valid yay!
    return PoolState.VALID_INPUT;
  });

  return {
    state,
    aPerBRatioMessage,
    bPerARatioMessage,
    aPerBRatioProjectedMessage,
    bPerARatioProjectedMessage,
    shareOfPool,
    shareOfPoolPercent,
    preExistingPool,
    totalLiquidityProviderUnits,
    totalPoolUnits,
    tokenAFieldAmount: tokenAField.fieldAmount,
    tokenBFieldAmount: tokenBField.fieldAmount,
  };
}
