import { Ref, ref } from "@vue/reactivity";
import { AssetAmount, IAssetAmount, Network, Pool, Token } from "../entities";
import { SwapState, useSwapCalculator } from "./swapCalculator";

const TOKENS = {
  atk: Token({
    decimals: 18,
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
  rwn: Token({
    decimals: 18,
    symbol: "rwn",
    name: "Rowan",
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

  test("calculate swap usecase", () => {
    selectedField.value = "from";
    expect(state.value).toBe(SwapState.SELECT_TOKENS);

    marketPairFinder
      .mockImplementationOnce(() => {
        return Pool(
          AssetAmount(TOKENS.atk, "2000000000000"),
          AssetAmount(TOKENS.rwn, "1000000000000")
        );
      })
      .mockImplementationOnce(() => {
        return Pool(
          AssetAmount(TOKENS.btk, "1000000000000"),
          AssetAmount(TOKENS.rwn, "1000000000000")
        );
      });

    balances.value = [
      AssetAmount(TOKENS.atk, "1000"),
      AssetAmount(TOKENS.btk, "1000"),
      AssetAmount(TOKENS.eth, "1234"),
    ];
    fromSymbol.value = "atk";
    toSymbol.value = "btk";

    expect(state.value).toBe(SwapState.ZERO_AMOUNTS);

    fromAmount.value = "100";

    expect(toAmount.value).toBe("49.99999999"); // 1 ATK ~= 0.5 BTK
    expect(state.value).toBe(SwapState.VALID_INPUT);

    selectedField.value = null; // deselect

    expect(fromAmount.value).toBe("100.0");

    selectedField.value = "to"; // select to field

    toAmount.value = "50"; // set to amount to 100
    expect(fromAmount.value).toBe("100.00000004");
    expect(toAmount.value).toBe("50");

    selectedField.value = null; // deselect
    selectedField.value = "from"; // select from field
    expect(toAmount.value).toBe("50.0");

    fromAmount.value = "10000";

    expect(state.value).toBe(SwapState.INSUFFICIENT_FUNDS);

    expect(priceMessage.value).toBe("0.500000 BTK per ATK");
  });

  test("Avoid division by zero", () => {
    selectedField.value = "from";
    fromAmount.value = "0";
    toAmount.value = "0";
    marketPairFinder
      .mockImplementationOnce(() =>
        Pool(
          AssetAmount(TOKENS.atk, "1000000"),
          AssetAmount(TOKENS.rwn, "1000000")
        )
      )
      .mockImplementationOnce(() =>
        Pool(
          AssetAmount(TOKENS.btk, "2000000"),
          AssetAmount(TOKENS.rwn, "1000000")
        )
      );
    fromSymbol.value = "atk";
    toSymbol.value = "btk";
    expect(priceMessage.value).toBe("");
  });
});
