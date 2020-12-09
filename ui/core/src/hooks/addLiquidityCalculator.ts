import { computed, Ref } from "@vue/reactivity";
import { Asset, AssetAmount, IAssetAmount, Pool } from "../entities";
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
  fromAmount: Ref<string>;
  fromSymbol: Ref<string | null>;
  toAmount: Ref<string>;
  toSymbol: Ref<string | null>;
  balances: Ref<IAssetAmount[]>;
  selectedField: Ref<"from" | "to" | null>;
  poolFinder: (a: Asset | string, b: Asset | string) => Ref<Pool> | null;
}) {
  const fromField = useField(input.fromAmount, input.fromSymbol);
  const toField = useField(input.toAmount, input.toSymbol);
  const balanceMap = useBalances(input.balances);

  const fromBalance = computed(() => {
    return input.fromSymbol.value
      ? balanceMap.value.get(input.fromSymbol.value) ?? null
      : null;
  });

  const toBalance = computed(() => {
    return input.toSymbol.value
      ? balanceMap.value.get(input.toSymbol.value) ?? null
      : null;
  });

  const fromBalanceOverdrawn = computed(
    () => !fromBalance.value?.greaterThan(fromField.fieldAmount.value || "0")
  );

  const toBalanceOverdrawn = computed(
    () => !toBalance.value?.greaterThan(toField.fieldAmount.value || "0")
  );

  const preExistingPool = computed(() => {
    if (
      !fromField.fieldAmount.value ||
      !toField.fieldAmount.value ||
      !fromField.asset.value ||
      !toField.asset.value
    )
      return null;

    // Find pool from poolFinder
    const pool = input.poolFinder(
      fromField.asset.value.symbol,
      toField.asset.value.symbol
    );
    return pool?.value ?? null;
  });

  const liquidityPool = computed(() => {
    if (
      !fromField.fieldAmount.value ||
      !toField.fieldAmount.value ||
      !fromField.asset.value ||
      !toField.asset.value
    )
      return null;

    return (
      preExistingPool.value ||
      Pool(
        AssetAmount(fromField.asset.value, "0"),
        AssetAmount(toField.asset.value, "0")
      )
    );
  });

  const shareOfPool = computed(() => {
    if (
      !liquidityPool.value ||
      !toField.fieldAmount.value ||
      !fromField.fieldAmount.value
    )
      return new Fraction("0");

    const [units, lpUnits] = liquidityPool.value.calculatePoolUnits(
      toField.fieldAmount.value,
      fromField.fieldAmount.value
    );

    // if no units lp owns 100% of pool
    return units.equalTo("0") ? new Fraction("1") : lpUnits.divide(units);
  });

  const shareOfPoolPercent = computed(() => {
    return `${shareOfPool.value.multiply("100").toFixed(2)}%`;
  });

  const aPerBRatioMessage = computed(() => {
    const aAmount = fromField.fieldAmount.value;
    const bAmount = toField.fieldAmount.value;
    if (!aAmount || !bAmount) return "";
    if (bAmount.equalTo("0")) return "";
    return `${aAmount
      .divide(bAmount)
      .toFixed(
        8
      )} ${aAmount.asset.symbol.toUpperCase()} per ${bAmount.asset.symbol.toUpperCase()}`;
  });

  const bPerARatioMessage = computed(() => {
    const aAmount = fromField.fieldAmount.value;
    const bAmount = toField.fieldAmount.value;
    if (!aAmount || !bAmount) return "";
    if (aAmount.equalTo("0")) return "";

    return `${bAmount
      .divide(aAmount)
      .toFixed(
        8
      )} ${bAmount.asset.symbol.toUpperCase()} per ${aAmount.asset.symbol.toUpperCase()}`;
  });

  const state = computed(() => {
    if (!input.fromSymbol.value || !input.toSymbol.value)
      return PoolState.SELECT_TOKENS;
    if (
      fromField.fieldAmount.value?.equalTo("0") ||
      toField.fieldAmount.value?.equalTo("0")
    )
      return PoolState.ZERO_AMOUNTS;
    if (fromBalanceOverdrawn.value || toBalanceOverdrawn.value) {
      return PoolState.INSUFFICIENT_FUNDS;
    }

    return PoolState.VALID_INPUT;
  });
  return {
    state,
    aPerBRatioMessage,
    bPerARatioMessage,
    shareOfPool,
    shareOfPoolPercent,
    preExistingPool,
    fromFieldAmount: fromField.fieldAmount,
    toFieldAmount: toField.fieldAmount,
  };
}
