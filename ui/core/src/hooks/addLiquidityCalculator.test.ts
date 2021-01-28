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

const ZERO = new Fraction("0");

describe("addLiquidityCalculator", () => {
  // input
  const tokenAAmount: Ref<string> = ref("0");
  const tokenASymbol: Ref<string | null> = ref(null);
  const tokenBAmount: Ref<string> = ref("0");
  const tokenBSymbol: Ref<string | null> = ref(null);
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
    LiquidityProvider(ATK, ZERO, akasha.address, ZERO, ZERO)
  ) as Ref<LiquidityProvider | null>; // ? not sure why we need to cast

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
      tokenAAmount,
      tokenBAmount,
      tokenASymbol,
      tokenBSymbol,
      poolFinder,
      liquidityProvider,
    }));

    balances.value = [];
    tokenAAmount.value = "0";
    tokenBAmount.value = "0";
    tokenASymbol.value = null;
    tokenBSymbol.value = null;
  });

  afterEach(() => {
    poolFinder.mockReset();
  });

  test("poolCalculator ratio messages", () => {
    tokenAAmount.value = "1000";
    tokenBAmount.value = "500";
    tokenASymbol.value = "atk";
    tokenBSymbol.value = "rowan";

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
            new Fraction("1000000")
          )
        ) as Ref<Pool>
    );

    // Liquidity provider already owns 1000 pool units (1000000 from another investor)
    liquidityProvider.value = LiquidityProvider(
      ATK,
      new Fraction("500000"),
      akasha.address,
      new Fraction("500000"),
      new Fraction("500000")
    );

    // Add 1000 of both tokens
    tokenAAmount.value = "500000";
    tokenBAmount.value = "500000";
    tokenASymbol.value = "atk";
    tokenBSymbol.value = "rowan";

    expect(aPerBRatioMessage.value).toBe("1.00000000");
    expect(bPerARatioMessage.value).toBe("1.00000000");

    // New shareOfPoolPercent for liquidity provider (inc prev liquidity)
    //2000/1002000 = 0.001996007984031936 so roughtly 0.2%
    expect(shareOfPoolPercent.value).toBe("66.67%");

    // New pool units for liquidity provider (inc prev liquidity)
    expect(totalLiquidityProviderUnits.value).toBe("1000000");

    expect(totalPoolUnits.value).toBe("1500000");
  });

  test("poolCalculator with preexisting pool but no preexisting liquidity", () => {
    // Pool exists with 1001000 preexisting units 1000 of which are owned by this lp
    poolFinder.mockImplementation(
      () =>
        ref(
          Pool(
            AssetAmount(ATK, "1000000"),
            AssetAmount(ROWAN, "1000000"),
            new Fraction("1000000")
          )
        ) as Ref<Pool>
    );

    // Liquidity provider is null
    liquidityProvider.value = null;

    // Add 1000 of both tokens
    tokenAAmount.value = "500000";
    tokenBAmount.value = "500000";
    tokenASymbol.value = "atk";
    tokenBSymbol.value = "rowan";

    expect(aPerBRatioMessage.value).toBe("1.00000000");
    expect(bPerARatioMessage.value).toBe("1.00000000");

    // New shareOfPoolPercent for liquidity provider (inc prev liquidity)
    //2000/1002000 = 0.001996007984031936 so roughtly 0.2%
    expect(shareOfPoolPercent.value).toBe("33.33%");

    // New pool units for liquidity provider (inc prev liquidity)
    expect(totalLiquidityProviderUnits.value).toBe("500000");

    expect(totalPoolUnits.value).toBe("1500000");
  });

  test("Can handle division by zero", () => {
    liquidityProvider.value = LiquidityProvider(ATK, ZERO, "", ZERO, ZERO);
    tokenAAmount.value = "0";
    tokenBAmount.value = "0";
    tokenASymbol.value = "atk";
    tokenBSymbol.value = "rowan";
    expect(state.value).toBe(PoolState.ZERO_AMOUNTS);
    expect(aPerBRatioMessage.value).toBe("");
    expect(bPerARatioMessage.value).toBe("");
    expect(shareOfPoolPercent.value).toBe("0.00%");
  });

  test("Don't allow rowan === 0 when creating new pool", () => {
    tokenAAmount.value = "1000";
    tokenBAmount.value = "0";
    tokenASymbol.value = "atk";
    tokenBSymbol.value = "rowan";
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
    tokenAAmount.value = "1000";
    tokenBAmount.value = "0";
    tokenASymbol.value = "atk";
    tokenBSymbol.value = "rowan";
    expect(state.value).toBe(PoolState.VALID_INPUT);
    expect(aPerBRatioMessage.value).toBe("N/A");
    expect(bPerARatioMessage.value).toBe("N/A");
  });

  test("insufficient funds", () => {
    balances.value = [AssetAmount(ATK, "100"), AssetAmount(ROWAN, "100")];
    tokenAAmount.value = "1000";
    tokenBAmount.value = "500";
    tokenASymbol.value = "atk";
    tokenBSymbol.value = "rowan";

    expect(state.value).toBe(PoolState.INSUFFICIENT_FUNDS);
  });

  test("valid funds below limit", () => {
    balances.value = [AssetAmount(ATK, "1000"), AssetAmount(ROWAN, "500")];
    tokenAAmount.value = "999";
    tokenBAmount.value = "499";
    tokenASymbol.value = "atk";
    tokenBSymbol.value = "rowan";
    expect(state.value).toBe(PoolState.VALID_INPUT);
  });

  test("valid funds at limit", () => {
    balances.value = [AssetAmount(ATK, "1000"), AssetAmount(ROWAN, "500")];
    tokenAAmount.value = "1000";
    tokenBAmount.value = "500";
    tokenASymbol.value = "atk";
    tokenBSymbol.value = "rowan";

    expect(state.value).toBe(PoolState.VALID_INPUT);
  });

  test("invalid funds above limit", () => {
    balances.value = [AssetAmount(ATK, "1000"), AssetAmount(ROWAN, "500")];
    tokenAAmount.value = "1001";
    tokenBAmount.value = "501";
    tokenASymbol.value = "atk";
    tokenBSymbol.value = "rowan";

    expect(state.value).toBe(PoolState.INSUFFICIENT_FUNDS);
  });
});
