// TODO remove refs dependency and move to `actions/clp/calculateAddLiquidity`

import { computed, Ref } from "@vue/reactivity";
import {
  Asset,
  AssetAmount,
  IAssetAmount,
  LiquidityProvider,
  Pool,
} from "../entities";
import { Fraction } from "../entities/fraction/Fraction";
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
    return pool?.value ?? null;
  });

  const liquidityPool = computed(() => {
    if (
      !tokenAField.fieldAmount.value ||
      !tokenBField.fieldAmount.value ||
      !tokenAField.asset.value ||
      !tokenBField.asset.value
    )
      return null;

    return (
      preExistingPool.value ||
      Pool(
        AssetAmount(tokenAField.asset.value, "0"),
        AssetAmount(tokenBField.asset.value, "0")
      )
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
    return `${shareOfPool.value.multiply("100").toFixed(2)}%`;
  });

  const aPerBRatio = computed(() => {
    if (!liquidityPool.value || !tokenAField.asset.value) {
      return null;
    }
    const amount = AssetAmount(tokenAField.asset.value, "1");
    const swapVal = liquidityPool.value.calcSwapResult(amount);
    const providerFee = liquidityPool.value.calcProviderFee(amount);
    const swapAmountNoFee = swapVal.add(providerFee);
    return amount.divide(swapAmountNoFee);
  });

  const bPerARatio = computed(() => {
    if (!aPerBRatio.value || !tokenBField.asset.value) {
      return null;
    }
    const amount = AssetAmount(tokenBField.asset.value, "1");
    return amount.divide(aPerBRatio.value);
  });

  const aPerBRatioMessage = computed(() => {
    if (!preExistingPool.value || !aPerBRatio.value) {
      return "N/A";
    }

    return aPerBRatio.value.toFixed(8);
  });

  const bPerARatioMessage = computed(() => {
    if (!preExistingPool.value || !bPerARatio.value) {
      return "N/A";
    }

    return bPerARatio.value.toFixed(8);
  });

  const state = computed(() => {
    const aAmount = tokenAField.fieldAmount.value;
    const bAmount = tokenBField.fieldAmount.value;

    if (!input.tokenASymbol.value || !input.tokenBSymbol.value)
      return PoolState.SELECT_TOKENS;

    if (!aAmount || aAmount.equalTo("0")) return PoolState.ZERO_AMOUNTS;

    if (!bAmount || bAmount.equalTo("0"))
      // if rowan is 0 or empty ...
      // allow if the pool exists
      // invalid if the pool doesn't exist - ""
      return preExistingPool.value
        ? PoolState.VALID_INPUT
        : PoolState.ZERO_AMOUNTS;

    if (fromBalanceOverdrawn.value || toBalanceOverdrawn.value)
      return PoolState.INSUFFICIENT_FUNDS;

    return PoolState.VALID_INPUT;
  });

  return {
    state,
    aPerBRatioMessage,
    bPerARatioMessage,
    shareOfPool,
    shareOfPoolPercent,
    preExistingPool,
    totalLiquidityProviderUnits,
    totalPoolUnits,
    tokenAFieldAmount: tokenAField.fieldAmount,
    tokenBFieldAmount: tokenBField.fieldAmount,
  };
}
