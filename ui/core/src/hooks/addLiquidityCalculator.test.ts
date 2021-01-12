import { ref, Ref } from "@vue/reactivity";
import { AssetAmount, Coin, IAssetAmount, Network, Pool } from "../entities";
import { Fraction } from "../entities/fraction/Fraction";
import { getTestingTokens } from "../test/utils/getTestingToken";
import { PoolState, usePoolCalculator } from "./addLiquidityCalculator";

const [ATK, ROWAN] = getTestingTokens(["ATK", "ROWAN"]);

describe("usePoolCalculator", () => {
  // input
  const fromAmount: Ref<string> = ref("0");
  const fromSymbol: Ref<string | null> = ref(null);
  const toAmount: Ref<string> = ref("0");
  const toSymbol: Ref<string | null> = ref(null);
  const balances = ref([]) as Ref<IAssetAmount[]>;
  const selectedField: Ref<"from" | "to" | null> = ref("from");
  const poolFinder = jest.fn<Ref<Pool> | null, any>(() => null);

  // output
  let aPerBRatioMessage: Ref<string>;
  let bPerARatioMessage: Ref<string>;
  let shareOfPool: Ref<Fraction>;
  let shareOfPoolPercent: Ref<string>;
  let state: Ref<PoolState>;

  beforeEach(() => {
    ({
      state,
      aPerBRatioMessage,
      bPerARatioMessage,
      shareOfPool,
      shareOfPoolPercent,
    } = usePoolCalculator({
      balances,
      fromAmount,
      toAmount,
      fromSymbol,
      selectedField,
      toSymbol,
      poolFinder,
    }));

    balances.value = [];
    fromAmount.value = "0";
    toAmount.value = "0";
    fromSymbol.value = null;
    toSymbol.value = null;
  });

  test("poolCalculator ratio messages", () => {
    fromAmount.value = "1000";
    toAmount.value = "500";
    fromSymbol.value = "atk";
    toSymbol.value = "rowan";

    expect(aPerBRatioMessage.value).toBe("2.00000000");
    expect(bPerARatioMessage.value).toBe("0.50000000");
    expect(shareOfPoolPercent.value).toBe("100.00%");
  });

  test("poolCalculator with preexisting pool", () => {
    poolFinder.mockImplementation(
      () =>
        ref(
          Pool(AssetAmount(ATK, "1000000"), AssetAmount(ROWAN, "1000000"))
        ) as Ref<Pool>
    );
    fromAmount.value = "1000";
    toAmount.value = "500";
    fromSymbol.value = "atk";
    toSymbol.value = "rowan";

    expect(aPerBRatioMessage.value).toBe("2.00000000");
    expect(bPerARatioMessage.value).toBe("0.50000000");
    expect(shareOfPoolPercent.value).toBe("0.07%");
  });

  test("Can handle division by zero", () => {
    fromAmount.value = "0";
    toAmount.value = "0";
    fromSymbol.value = "atk";
    toSymbol.value = "rowan";
    expect(state.value).toBe(PoolState.ZERO_AMOUNTS);
    expect(aPerBRatioMessage.value).toBe("");
    expect(bPerARatioMessage.value).toBe("");
    expect(shareOfPoolPercent.value).toBe("0.00%");
  });
});
