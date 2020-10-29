import { Ref, ref } from "@vue/reactivity";
import { AssetAmount, IAssetAmount, Network, Pair, Token } from "../entities";
import { useSwapCalculator } from "./swapCalculator";

const TOKENS = {
  atk: Token({
    decimals: 6,
    symbol: "atk",
    name: "AppleToken",
    address: "123",
    network: Network.ETHEREUM,
  }),
  btk: Token({
    decimals: 18,
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
  const marketPairFinder = jest.fn();

  // output
  let nextStepMessage: Ref<string>;
  let priceMessage: Ref<string | null>;

  beforeEach(() => {
    ({ nextStepMessage, priceMessage } = useSwapCalculator({
      balances,
      fromAmount,
      toAmount,
      fromSymbol,
      selectedField,
      toSymbol,
      marketPairFinder,
    }));
  });

  test("calculates wallet not attached", () => {
    selectedField.value = "from";
    expect(nextStepMessage.value).toBe("Select tokens");
    marketPairFinder.mockImplementationOnce(() =>
      Pair(AssetAmount(TOKENS.atk, "1000"), AssetAmount(TOKENS.btk, "2000"))
    );
    fromSymbol.value = "atk";
    toSymbol.value = "btk";
    expect(nextStepMessage.value).toBe("You cannot swap 0");

    balances.value = [
      AssetAmount(TOKENS.eth, "1234"),
      AssetAmount(TOKENS.btk, "1000"),
      AssetAmount(TOKENS.atk, "1000"),
    ];
    fromAmount.value = "100";
    expect(toAmount.value).toBe("200.000000");
    expect(nextStepMessage.value).toBe("Swap"); // Should be something else if values are 0

    selectedField.value = "to";
    toAmount.value = "100";
    expect(fromAmount.value).toBe("50.000000000000000000");
    expect(toAmount.value).toBe("100");
    selectedField.value = "from";
    expect(toAmount.value).toBe("100.000000");

    selectedField.value = "from";
    fromAmount.value = "10000";
    expect(nextStepMessage.value).toBe("Insufficient funds");

    expect(priceMessage.value).toBe("2.000000000000000000 BTK per ATK");
  });
});
