import { computed } from "@vue/reactivity";
import { Ref } from "vue";
import { IAssetAmount, Pair } from "../entities";
import { useField } from "./useField";

export function usePoolCalculator(input: {
  fromAmount: Ref<string>;
  fromSymbol: Ref<string | null>;
  toAmount: Ref<string>;
  toSymbol: Ref<string | null>;
  balances: Ref<IAssetAmount[]>;
  selectedField: Ref<"from" | "to" | null>;
}) {
  const fromField = useField(input.fromAmount, input.fromSymbol);
  const toField = useField(input.toAmount, input.toSymbol);
  const poolPair = computed(() => {
    if (!fromField.fieldAmount.value || !toField.fieldAmount.value) return null;

    return Pair(fromField.fieldAmount.value, toField.fieldAmount.value);
  });

  const aPerBRatioMessage = computed(() => {
    const asset = fromField.asset.value;
    const pair = poolPair.value;
    if (!asset || !pair) return "";

    return `${pair
      .priceAsset(asset)
      .toFormatted()} per ${asset?.symbol.toUpperCase()}`;
  });

  const bPerARatioMessage = computed(() => {
    const asset = toField.asset.value;
    const pair = poolPair.value;
    if (!asset || !pair) return "";

    return `${pair
      .priceAsset(asset)
      .toFormatted()} per ${asset?.symbol.toUpperCase()}`;
  });

  const nextStepMessage = computed(() => {
    return "";
  });

  return {
    aPerBRatioMessage,
    bPerARatioMessage,
    nextStepMessage,
    fromFieldAmount: fromField.fieldAmount,
    toFieldAmount: toField.fieldAmount,
    toAmount: input.toAmount,
    fromAmount: input.fromAmount,
  };
}
