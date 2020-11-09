import { computed, ref } from "@vue/reactivity";
import { Ref } from "vue";
import { Asset, IAssetAmount, Pool } from "../entities";

import { useField } from "./useField";

export enum PoolState {
  SELECT_TOKENS,
  ZERO_AMOUNTS,
  INSUFFICIENT_FUNDS,
  VALID_INPUT,
}

export function usePoolCalculator(input: {
  fromAmount: Ref<string>;
  fromSymbol: Ref<string | null>;
  toAmount: Ref<string>;
  toSymbol: Ref<string | null>;
  balances: Ref<IAssetAmount[]>;
  selectedField: Ref<"from" | "to" | null>;
  marketPairFinder: (a: Asset | string, b: Asset | string) => Pool | null;
}) {
  const fromField = useField(input.fromAmount, input.fromSymbol);
  const toField = useField(input.toAmount, input.toSymbol);

  // const liquidityPool = computed(() => {
  //   if (!fromField.fieldAmount.value || !toField.fieldAmount.value) return null;
  //   return Pool(fromField.fieldAmount.value, toField.fieldAmount.value);
  // });

  // const balanceMap = useBalances(input.balances);

  // const fromBalance = computed(() => {
  //   return input.fromSymbol.value
  //     ? balanceMap.value.get(input.fromSymbol.value) ?? null
  //     : null;
  // });

  // const toBalance = computed(() => {
  //   return input.toSymbol.value
  //     ? balanceMap.value.get(input.toSymbol.value) ?? null
  //     : null;
  // });

  // const aPerBRatioMessage = computed(() => {
  //   const asset = fromField.asset.value;
  //   const pair = liquidityPool.value;
  //   return assetPriceMessage(asset, pair);
  // });

  // const bPerARatioMessage = computed(() => {
  //   const asset = toField.asset.value;
  //   const pair = liquidityPool.value;
  //   return assetPriceMessage(asset, pair);
  // });

  // const shareOfPool = computed(() => {
  //   if (!liquidityPool.value) return "";
  //   const [ama, amb] = liquidityPool.value.amounts;
  //   const marketPool = input.marketPoolFinder(ama.asset, amb.asset);

  //   // TODO: Naive calculation need to check this is correct
  //   // get the sum of the market pair
  //   const marketPoolSum = marketPool
  //     ? marketPool.amounts.reduce(
  //         (acc, amount) => amount.add(acc),
  //         new Fraction("0")
  //       )
  //     : new Fraction("0");

  //   // TODO: Naive calculation need to check this is correct
  //   // get the sum of the liquidity pair being created
  //   const liquidityPoolSum = liquidityPool.value.amounts.reduce(
  //     (acc, amount) => amount.add(acc),
  //     new Fraction("0")
  //   );

  //   // TODO: Naive calculation need to check this is correct
  //   // Work out the total share of the pool by adding
  //   // all the amounts up and dividing by the liquidity pair
  //   if (!liquidityPoolSum || liquidityPoolSum.equalTo("0")) return "";
  //   return `${liquidityPoolSum
  //     .divide(marketPoolSum.add(liquidityPoolSum))
  //     .multiply(new Fraction("100"))
  //     .toFixed(2)}%`;
  // });

  // const fromBalanceOverdrawn = computed(
  //   () => !fromBalance.value?.greaterThan(fromField.fieldAmount.value || "0")
  // );
  // const toBalanceOverdrawn = computed(
  //   () => !toBalance.value?.greaterThan(toField.fieldAmount.value || "0")
  // );

  // const state = computed(() => {
  //   if (!input.fromSymbol.value || !input.toSymbol.value)
  //     return PoolState.SELECT_TOKENS;
  //   if (
  //     fromField.fieldAmount.value?.equalTo("0") &&
  //     toField.fieldAmount.value?.equalTo("0")
  //   )
  //     return PoolState.ZERO_AMOUNTS;
  //   if (fromBalanceOverdrawn.value || toBalanceOverdrawn.value) {
  //     return PoolState.INSUFFICIENT_FUNDS;
  //   }

  //   return PoolState.VALID_INPUT;
  // });

  // const nextStepAllowed = computed(() => {
  //   state.value === PoolState.VALID_INPUT;
  // });

  // effect(() => {
  //   // Deselect a field formats all values
  //   if (input.selectedField.value === null) {
  //     const fromAsset = fromField.asset.value;
  //     if (fromAsset) {
  //       input.fromAmount.value = AssetAmount(
  //         fromAsset,
  //         input.fromAmount.value
  //       ).toFixed();
  //     }

  //     const toAsset = fromField.asset.value;
  //     if (toAsset) {
  //       input.toAmount.value = AssetAmount(
  //         toAsset,
  //         input.toAmount.value
  //       ).toFixed();
  //     }
  //   }
  // });

  return {
    aPerBRatioMessage: "",
    bPerARatioMessage: "",
    shareOfPool: "",
    state: ref(PoolState.SELECT_TOKENS),
    nextStepAllowed: true,
    fromFieldAmount: fromField.fieldAmount,
    toFieldAmount: toField.fieldAmount,
    toAmount: input.toAmount,
    fromAmount: input.fromAmount,
  };
}

export function removeLiquidity(input: {
  fromAmount: Ref<string>;
  fromSymbol: Ref<string | null>;
  toAmount: Ref<string>;
  toSymbol: Ref<string | null>;
  balances: Ref<IAssetAmount[]>;
  // Get a pair that represents the balance the current user has contibuted
  userPoolFinder: (a: Asset | string, b: Asset | string) => Pool | null;
}) {
  const fromField = useField(input.fromAmount, input.fromSymbol);
  const toField = useField(input.toAmount, input.toSymbol);

  const userPool = computed(() => {
    if (!fromField.asset.value) return null;
    if (!toField.asset.value) return null;

    return input.userPoolFinder(fromField.asset.value, toField.asset.value);
  });

  const fromBalanceOverdrawn = computed(() => {
    if (!fromField.fieldAmount.value) return null;
    return userPool.value?.priceAsset(fromField.fieldAmount.value.asset);
  });

  const toBalanceOverdrawn = computed(() => {
    if (!toField.fieldAmount.value) return null;
    return userPool.value?.priceAsset(toField.fieldAmount.value.asset);
  });

  // Bit hard to work out how much of this works without getting clarity
  const state = computed(() => {
    if (!input.fromSymbol.value || !input.toSymbol.value)
      return PoolState.SELECT_TOKENS;
    if (
      fromField.fieldAmount.value?.equalTo("0") &&
      toField.fieldAmount.value?.equalTo("0")
    )
      return PoolState.ZERO_AMOUNTS;
    if (fromBalanceOverdrawn.value || toBalanceOverdrawn.value) {
      return PoolState.INSUFFICIENT_FUNDS;
    }

    return PoolState.VALID_INPUT;
  });
}
