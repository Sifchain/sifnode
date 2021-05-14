import { ref, Ref } from "@vue/reactivity";

import {
  Amount,
  AssetAmount,
  LiquidityProvider,
  Pool,
} from "../../../entities";
import { getTestingTokens } from "../../../test/utils/getTestingToken";
import { PoolState } from "./addLiquidityCalculator";
import { useRemoveLiquidityCalculator } from "./removeLiquidityCalculator";

const [CATK, ROWAN] = getTestingTokens(["CATK", "ROWAN"]);

const ZERO = Amount("0");

describe("useRemoveLiquidityCalculator", () => {
  // input
  const asymmetry: Ref<string> = ref("0");
  const externalAssetSymbol: Ref<string | null> = ref(null);
  const nativeAssetSymbol: Ref<string | null> = ref(null);
  const sifAddress: Ref<string> = ref("12345678asFDSghkjg");
  const wBasisPoints: Ref<string> = ref("5000");
  const liquidityProvider: Ref<LiquidityProvider | null> = ref(null);
  const poolFinder = jest.fn<Ref<Pool> | null, any>(() => null);

  // output
  let withdrawExternalAssetAmount: Ref<string | null> = ref(null);
  let withdrawNativeAssetAmount: Ref<string | null> = ref(null);
  let state: Ref<PoolState> = ref(0);

  // watch fires when certain wBasisPoints, asymmetry, or liquidityProvider changes
  function simulateWatch() {
    const calcData = useRemoveLiquidityCalculator({
      asymmetry,
      externalAssetSymbol,
      liquidityProvider,
      poolFinder,
      nativeAssetSymbol,
      sifAddress,
      wBasisPoints,
    });
    state.value = calcData.state;
    withdrawExternalAssetAmount.value = calcData.withdrawExternalAssetAmount;
    withdrawNativeAssetAmount.value = calcData.withdrawNativeAssetAmount;
  }

  beforeEach(() => {
    simulateWatch();
  });

  test("displays the correct withdrawal amounts", async () => {
    liquidityProvider.value = LiquidityProvider(
      CATK,
      Amount("100000000000000000000000"),
      "sif123456876512341234",
      ZERO,
      ZERO,
    );

    poolFinder.mockImplementation(
      () =>
        ref(
          Pool(
            AssetAmount(CATK, "1000000000000000000000000"),
            AssetAmount(ROWAN, "1000000000000000000000000"),
            Amount("1000000000000000000000000"),
          ),
        ) as Ref<Pool>,
    );

    expect(state.value).toBe(PoolState.SELECT_TOKENS);
    asymmetry.value = "0";
    externalAssetSymbol.value = "catk";
    nativeAssetSymbol.value = "rowan";
    sifAddress.value = "sif123456876512341234";
    wBasisPoints.value = "0";
    simulateWatch();

    expect(state.value).toBe(PoolState.ZERO_AMOUNTS);
    wBasisPoints.value = "10000";
    simulateWatch();

    expect(state.value).toBe(PoolState.VALID_INPUT);

    expect(withdrawExternalAssetAmount.value).toEqual("100000.000000");
    expect(withdrawNativeAssetAmount.value).toEqual("100000.000000");
    asymmetry.value = "10000";
    simulateWatch();

    expect(withdrawExternalAssetAmount.value).toEqual("181000.000000");
    expect(withdrawNativeAssetAmount.value).toEqual("0.000000");
    wBasisPoints.value = "5000";
    simulateWatch();

    expect(withdrawExternalAssetAmount.value).toEqual("95125.000000");
    expect(withdrawNativeAssetAmount.value).toEqual("0.000000");
  });
});
