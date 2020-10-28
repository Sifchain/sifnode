import { Ref, ref } from "@vue/reactivity";
import {
  Asset,
  AssetAmount,
  IAssetAmount,
  Network,
  Pair,
  Token,
} from "../../../../core";
import { usePoolCalculator, useField } from "./usePoolCalculator";

const TOKENS = {
  atk: Token({
    decimals: 6,
    symbol: "atk",
    name: "AppleToken",
    address: "123",
    network: Network.ETHEREUM,
  }),
  btk: Token({
    decimals: 6,
    symbol: "btk",
    name: "BananaToken",
    address: "1234",
    network: Network.ETHEREUM,
  }),
  eth: Token({
    decimals: 18,
    symbol: "eth",
    name: "Ethereum",
    address: "1234",
    network: Network.ETHEREUM,
  }),
};

describe("swapCalculator", () => {
  // input
  const fromAmount: Ref<string> = ref("0");
  const fromSymbol: Ref<string | null> = ref(null);
  const toAmount: Ref<string> = ref("0");
  const toSymbol: Ref<string | null> = ref(null);
  const balances = ref([]) as Ref<IAssetAmount[]>;
  const selectedField: Ref<"from" | "to" | null> = ref("from");

  // output

  let aPerBRatioMessage: Ref<string>;
  let bPerARatioMessage: Ref<string>;
  beforeEach(() => {
    ({ aPerBRatioMessage, bPerARatioMessage } = usePoolCalculator({
      balances,
      fromAmount,
      toAmount,
      fromSymbol,
      selectedField,
      toSymbol,
    }));
  });

  test("poolCalculator ratio messages", () => {
    fromAmount.value = "1000";
    toAmount.value = "500";
    fromSymbol.value = "atk";
    toSymbol.value = "btk";

    expect(aPerBRatioMessage.value).toBe("0.500000 BTK per ATK");
    expect(bPerARatioMessage.value).toBe("2.000000 ATK per BTK");
  });
});
