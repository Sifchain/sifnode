import { ref, Ref } from "@vue/reactivity";
import {
  AssetAmount,
  Coin,
  IAssetAmount,
  LiquidityProvider,
  Network,
  Pool,
} from "../entities";
import { Fraction, IFraction } from "../entities/fraction/Fraction";
import { akasha } from "../test/utils/accounts";
import { getTestingTokens } from "../test/utils/getTestingToken";
import { PoolState, usePoolCalculator } from "./addLiquidityCalculator";

const [ATK, ROWAN] = getTestingTokens(["ATK", "ROWAN"]);

describe("addLiquidityCalculator", () => {
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
  let totalLiquidityProviderUnits: Ref<string>;
  let totalPoolUnits: Ref<string>;
  let shareOfPoolPercent: Ref<string>;
  let state: Ref<PoolState>;
  let liquidityProvider = ref(
    LiquidityProvider(ATK, new Fraction("0"), akasha.address)
  ) as Ref<LiquidityProvider>; // ? not sure why we need to cast

  beforeEach(() => {
    ({
      state,
      aPerBRatioMessage,
      bPerARatioMessage,
      shareOfPool,
      shareOfPoolPercent,
      totalLiquidityProviderUnits,
      totalPoolUnits,
    } = usePoolCalculator({
      balances,
      fromAmount,
      toAmount,
      fromSymbol,
      selectedField,
      toSymbol,
      poolFinder,
      liquidityProvider,
    }));

    balances.value = [];
    fromAmount.value = "0";
    toAmount.value = "0";
    fromSymbol.value = null;
    toSymbol.value = null;
  });

  afterEach(() => {
    poolFinder.mockReset();
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
    // Pool exists with 1001000 preexisting units 1000 of which are owned by this lp
    poolFinder.mockImplementation(
      () =>
        ref(
          Pool(
            AssetAmount(ATK, "1000000"),
            AssetAmount(ROWAN, "1000000"),
            new Fraction("1001000")
          )
        ) as Ref<Pool>
    );

    // Liquidity provider already owns 1000 pool units (1000000 from another investor)
    liquidityProvider.value = LiquidityProvider(
      ATK,
      new Fraction("1000"),
      akasha.address
    );

    // Add 1000 of both tokens
    fromAmount.value = "1000";
    toAmount.value = "1000";
    fromSymbol.value = "atk";
    toSymbol.value = "rowan";

    expect(aPerBRatioMessage.value).toBe("1.00000000");
    expect(bPerARatioMessage.value).toBe("1.00000000");

    // New shareOfPoolPercent for liquidity provider (inc prev liquidity)
    //2000/1002000 = 0.001996007984031936 so roughtly 0.2%
    expect(shareOfPoolPercent.value).toBe("0.20%");

    // New pool units for liquidity provider (inc prev liquidity)
    expect(totalLiquidityProviderUnits.value).toBe("2001");

    expect(totalPoolUnits.value).toBe("1002001");
  });

  test("Can handle division by zero", () => {
    liquidityProvider.value = LiquidityProvider(ATK, new Fraction("0"), "");
    fromAmount.value = "0";
    toAmount.value = "0";
    fromSymbol.value = "atk";
    toSymbol.value = "rowan";
    expect(state.value).toBe(PoolState.ZERO_AMOUNTS);
    expect(aPerBRatioMessage.value).toBe("");
    expect(bPerARatioMessage.value).toBe("");
    expect(shareOfPoolPercent.value).toBe("0.00%");
  });

  test("Don't allow rowan === 0 when creating new pool", () => {
    fromAmount.value = "1000";
    toAmount.value = "0";
    fromSymbol.value = "atk";
    toSymbol.value = "rowan";
    expect(state.value).toBe(PoolState.ZERO_AMOUNTS);
    expect(aPerBRatioMessage.value).toBe("");
    expect(bPerARatioMessage.value).toBe("");
  });

  test("Allow rowan === 0 when adding to preExistingPool", () => {
    poolFinder.mockImplementation(
      () =>
        ref(
          Pool(AssetAmount(ATK, "1000000"), AssetAmount(ROWAN, "1000000"))
        ) as Ref<Pool>
    );
    fromAmount.value = "1000";
    toAmount.value = "0";
    fromSymbol.value = "atk";
    toSymbol.value = "rowan";
    expect(state.value).toBe(PoolState.VALID_INPUT);
    expect(aPerBRatioMessage.value).toBe("N/A");
    expect(bPerARatioMessage.value).toBe("N/A");
  });

  test("insufficient funds", () => {
    balances.value = [AssetAmount(ATK, "100"), AssetAmount(ROWAN, "100")];
    fromAmount.value = "1000";
    toAmount.value = "500";
    fromSymbol.value = "atk";
    toSymbol.value = "rowan";

    expect(state.value).toBe(PoolState.INSUFFICIENT_FUNDS);
  });

  test("valid funds below limit", () => {
    balances.value = [AssetAmount(ATK, "1000"), AssetAmount(ROWAN, "500")];
    fromAmount.value = "999";
    toAmount.value = "499";
    fromSymbol.value = "atk";
    toSymbol.value = "rowan";
    expect(state.value).toBe(PoolState.VALID_INPUT);
  });

  test("valid funds at limit", () => {
    balances.value = [AssetAmount(ATK, "1000"), AssetAmount(ROWAN, "500")];
    fromAmount.value = "1000";
    toAmount.value = "500";
    fromSymbol.value = "atk";
    toSymbol.value = "rowan";

    expect(state.value).toBe(PoolState.VALID_INPUT);
  });

  test("invalid funds above limit", () => {
    balances.value = [AssetAmount(ATK, "1000"), AssetAmount(ROWAN, "500")];
    fromAmount.value = "1001";
    toAmount.value = "501";
    fromSymbol.value = "atk";
    toSymbol.value = "rowan";

    expect(state.value).toBe(PoolState.INSUFFICIENT_FUNDS);
  });
});
