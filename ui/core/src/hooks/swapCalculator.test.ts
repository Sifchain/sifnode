import { Ref, ref } from "@vue/reactivity";
import { AssetAmount, IAssetAmount, Network, Pair, Token } from "../entities";
import { SwapState, useSwapCalculator } from "./swapCalculator";

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
  let state: Ref<SwapState>;
  let priceMessage: Ref<string | null>;

  beforeEach(() => {
    ({ state, priceMessage } = useSwapCalculator({
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
    expect(state.value).toBe(SwapState.SELECT_TOKENS);

    marketPairFinder.mockImplementationOnce(() =>
      Pair(AssetAmount(TOKENS.atk, "1000"), AssetAmount(TOKENS.btk, "2000"))
    );
    fromSymbol.value = "atk";
    toSymbol.value = "btk";
    expect(state.value).toBe(SwapState.ZERO_AMOUNTS);

    balances.value = [
      AssetAmount(TOKENS.eth, "1234"),
      AssetAmount(TOKENS.btk, "1000"),
      AssetAmount(TOKENS.atk, "1000"),
    ];
    fromAmount.value = "100";
    expect(toAmount.value).toBe("200.0");
    expect(state.value).toBe(SwapState.VALID_INPUT);

    selectedField.value = "to";
    toAmount.value = "100";
    expect(fromAmount.value).toBe("50.0");
    expect(toAmount.value).toBe("100");
    selectedField.value = "from";
    expect(toAmount.value).toBe("100.0");

    selectedField.value = "from";
    fromAmount.value = "10000";

    expect(state.value).toBe(SwapState.INSUFFICIENT_FUNDS);

    expect(priceMessage.value).toBe("2.0 BTK per ATK");
  });
});
