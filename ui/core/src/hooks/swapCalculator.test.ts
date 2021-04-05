import { Ref, ref } from "@vue/reactivity";
import { AssetAmount, IAssetAmount, Network, Pool, Token } from "../entities";
import { getTestingTokens } from "../test/utils/getTestingToken";
import { SwapState, useSwapCalculator } from "./swapCalculator";

const [ATK, BTK, ROWAN, ETH] = getTestingTokens(["ATK", "BTK", "ROWAN", "ETH"]);

describe("swapCalculator", () => {
  // input
  const fromAmount: Ref<string> = ref("0");
  const fromSymbol: Ref<string | null> = ref(null);
  const toAmount: Ref<string> = ref("0");
  const toSymbol: Ref<string | null> = ref(null);
  const balances = ref([]) as Ref<IAssetAmount[]>;
  const selectedField: Ref<"from" | "to" | null> = ref("from");
  const slippage = ref("0.5");

  // output
  let state: Ref<SwapState>;
  let priceMessage: Ref<string | null>;
  let priceImpact: Ref<string | null>;
  let providerFee: Ref<string | null>;
  let minimumReceived: Ref<IAssetAmount | null>;

  test("calculate swap usecase", () => {
    const pool1 = ref(
      Pool(
        AssetAmount(ATK, "2000000000000"),
        AssetAmount(ROWAN, "1000000000000"),
      ),
    ) as Ref<Pool | null>;

    const pool2 = ref(
      Pool(
        AssetAmount(BTK, "1000000000000"),
        AssetAmount(ROWAN, "1000000000000"),
      ),
    ) as Ref<Pool | null>;

    const poolFinder: any = jest.fn((a: string, b: string) => {
      if (a === "atk" && b === "rowan") {
        return pool1;
      } else {
        return pool2;
      }
    });

    ({
      state,
      priceMessage,
      priceImpact,
      providerFee,
      minimumReceived,
    } = useSwapCalculator({
      balances,
      fromAmount,
      toAmount,
      fromSymbol,
      selectedField,
      toSymbol,
      poolFinder,
      slippage,
    }));

    selectedField.value = "from";
    expect(state.value).toBe(SwapState.SELECT_TOKENS);

    balances.value = [
      AssetAmount(ATK, "1000"),
      AssetAmount(BTK, "1000"),
      AssetAmount(ETH, "1234"),
    ];

    fromSymbol.value = "atk";
    toSymbol.value = "btk";

    expect(state.value).toBe(SwapState.ZERO_AMOUNTS);

    fromAmount.value = "100";

    expect(toAmount.value).toBe("49.99999999"); // 1 ATK ~= 0.5 BTK
    expect(state.value).toBe(SwapState.VALID_INPUT);
    expect(minimumReceived.value?.toString()).toBe("49.749999990050000000 BTK");

    selectedField.value = null; // deselect

    expect(fromAmount.value).toBe("100.0");

    // Check background update
    pool1.value = Pool(
      AssetAmount(ATK, "1000000000000"),
      AssetAmount(ROWAN, "1000000000000"),
    );

    selectedField.value = "from";
    fromAmount.value = "1000";
    selectedField.value = null;

    expect(toAmount.value).toBe("999.999996");

    pool1.value = Pool(
      AssetAmount(ATK, "2000000000000"),
      AssetAmount(ROWAN, "1000000000000"),
    );

    selectedField.value = "from";
    fromAmount.value = "100";

    selectedField.value = null;

    selectedField.value = "to"; // select to field

    toAmount.value = "50"; // set to amount to 100
    expect(fromAmount.value).toBe("100.00000004");
    expect(toAmount.value).toBe("50");

    selectedField.value = null; // deselect
    selectedField.value = "from"; // select from field
    expect(toAmount.value).toBe("50.0");

    fromAmount.value = "10000";

    expect(state.value).toBe(SwapState.INSUFFICIENT_FUNDS);
    expect(toAmount.value).toBe("4999.9999");
    expect(priceMessage.value).toBe("0.500000 cTK per cTK");
    expect(priceImpact.value).toBe("0.000001");
    expect(providerFee.value).toBe("0.00005");
  });

  test("Avoid division by zero", () => {
    const pool1 = ref(
      Pool(AssetAmount(ATK, "1000000"), AssetAmount(ROWAN, "1000000")),
    ) as Ref<Pool | null>;

    const pool2 = ref(
      Pool(AssetAmount(BTK, "2000000"), AssetAmount(ROWAN, "1000000")),
    ) as Ref<Pool | null>;

    const poolFinder: any = jest.fn((a: string, b: string) => {
      if (a === "atk" && b === "rowan") {
        return pool1;
      } else {
        return pool2;
      }
    });

    ({ state, priceMessage, priceImpact, providerFee } = useSwapCalculator({
      balances,
      fromAmount,
      toAmount,
      fromSymbol,
      selectedField,
      toSymbol,
      poolFinder,
      slippage,
    }));

    selectedField.value = "from";
    fromAmount.value = "0";
    toAmount.value = "0";
    fromSymbol.value = "atk";
    toSymbol.value = "btk";
    expect(priceMessage.value).toBe("");
    expect(priceImpact.value).toBe("0.0");
    expect(providerFee.value).toBe("0.0");
  });

  test("insufficient funds", () => {
    balances.value = [AssetAmount(ATK, "100"), AssetAmount(ROWAN, "100")];
    fromAmount.value = "1000";
    toAmount.value = "500";
    fromSymbol.value = "atk";
    toSymbol.value = "rowan";

    expect(state.value).toBe(SwapState.INSUFFICIENT_FUNDS);
  });

  test("valid funds below limit", () => {
    balances.value = [AssetAmount(ATK, "1000"), AssetAmount(ROWAN, "500")];
    fromAmount.value = "999";
    toAmount.value = "499";
    fromSymbol.value = "atk";
    toSymbol.value = "rowan";
    expect(state.value).toBe(SwapState.VALID_INPUT);
  });

  test("valid funds at limit", () => {
    balances.value = [AssetAmount(ATK, "1000"), AssetAmount(ROWAN, "500")];
    fromAmount.value = "1000";
    toAmount.value = "500";
    fromSymbol.value = "atk";
    toSymbol.value = "rowan";

    expect(state.value).toBe(SwapState.VALID_INPUT);
  });

  test("invalid funds above limit", () => {
    balances.value = [AssetAmount(ATK, "1000"), AssetAmount(ROWAN, "500")];
    fromAmount.value = "1001";
    toAmount.value = "501";
    fromSymbol.value = "atk";
    toSymbol.value = "rowan";

    expect(state.value).toBe(SwapState.INSUFFICIENT_FUNDS);
  });
});
