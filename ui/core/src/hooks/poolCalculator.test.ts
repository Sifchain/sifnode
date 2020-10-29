import { Ref, ref } from "@vue/reactivity";
import { AssetAmount, IAssetAmount, Network, Pair, Token } from "../entities";
import { usePoolCalculator } from "./poolCalculator";

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

describe("usePoolCalculator", () => {
  // input
  const fromAmount: Ref<string> = ref("0");
  const fromSymbol: Ref<string | null> = ref(null);
  const toAmount: Ref<string> = ref("0");
  const toSymbol: Ref<string | null> = ref(null);
  const balances = ref([]) as Ref<IAssetAmount[]>;
  const selectedField: Ref<"from" | "to" | null> = ref("from");
  const marketPairFinder = jest.fn<Pair | null, any>(() => null);

  // output

  let aPerBRatioMessage: Ref<string>;
  let bPerARatioMessage: Ref<string>;
  let shareOfPool: Ref<string>;
  beforeEach(() => {
    ({ aPerBRatioMessage, bPerARatioMessage, shareOfPool } = usePoolCalculator({
      balances,
      fromAmount,
      toAmount,
      fromSymbol,
      selectedField,
      toSymbol,
      marketPairFinder,
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

  test("Can handle division by zero", () => {
    fromAmount.value = "0";
    toAmount.value = "0";
    fromSymbol.value = "atk";
    toSymbol.value = "btk";
    expect(aPerBRatioMessage.value).toBe("");
    expect(bPerARatioMessage.value).toBe("");
  });

  test("Calculate against a given pool", () => {
    marketPairFinder.mockImplementationOnce(() =>
      Pair(AssetAmount(TOKENS.atk, "2000"), AssetAmount(TOKENS.btk, "2000"))
    );

    fromAmount.value = "1000";
    toAmount.value = "1000";
    fromSymbol.value = "atk";
    toSymbol.value = "btk";

    // TODO: All the maths here are pretty naive need to double check with blockscience
    expect(aPerBRatioMessage.value).toBe("1.000000 BTK per ATK");
    expect(bPerARatioMessage.value).toBe("1.000000 ATK per BTK");
    expect(shareOfPool.value).toBe("33.33%");
  });
});
