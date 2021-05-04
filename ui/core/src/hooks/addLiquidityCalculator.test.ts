import { ref, Ref } from "@vue/reactivity";
import {
  Amount,
  Asset,
  AssetAmount,
  IAssetAmount,
  LiquidityProvider,
  Pool,
} from "../entities";
import { IAmount } from "../entities/Amount";
import { akasha } from "../test/utils/accounts";
import { getTestingTokens } from "../test/utils/getTestingToken";
import { PoolState, usePoolCalculator } from "./addLiquidityCalculator";

const [ATK, ROWAN, CTEST] = getTestingTokens(["ATK", "ROWAN", "CTEST"]);

const ZERO = Amount("0");

describe("addLiquidityCalculator", () => {
  // input
  const tokenAAmount: Ref<string> = ref("0");
  const tokenASymbol: Ref<string | null> = ref(null);
  const tokenBAmount: Ref<string> = ref("0");
  const tokenBSymbol: Ref<string | null> = ref(null);
  const asyncPooling: Ref<boolean> = ref(false);
  const lastFocusedTokenField: Ref<"A" | "B" | null> = ref(null);
  const balances = ref([]) as Ref<IAssetAmount[]>;
  const selectedField: Ref<"from" | "to" | null> = ref("from");
  const poolFinder = jest.fn<Ref<Pool> | null, any>(() => null);

  // output
  let aPerBRatioMessage: Ref<string>;
  let bPerARatioMessage: Ref<string>;
  let shareOfPool: Ref<IAmount>;
  let aPerBRatioProjectedMessage: Ref<string>;
  let bPerARatioProjectedMessage: Ref<string>;
  let totalLiquidityProviderUnits: Ref<string>;
  let totalPoolUnits: Ref<string>;
  let shareOfPoolPercent: Ref<string>;
  let state: Ref<PoolState>;
  let liquidityProvider = ref(
    LiquidityProvider(ATK, ZERO, akasha.address, ZERO, ZERO),
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
      aPerBRatioProjectedMessage,
      bPerARatioProjectedMessage,
    } = usePoolCalculator({
      balances,
      tokenAAmount,
      tokenBAmount,
      tokenASymbol,
      tokenBSymbol,
      asyncPooling,
      lastFocusedTokenField,
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
    liquidityProvider.value = null;
  });

  const ratios: any[] = [
    {
      poolExternal: "1000000000000000000000000000",
      poolNative: "1000000000000000000000000000",
      poolUnits: "1000000000000000000000000000",
      addedExternal: "10000000",
      addedNative: "10000000",
      externalSymbol: "atk",
      nativeSymbol: "rowan",
      preexistingLiquidity: {
        native: "0",
        external: "0",
        units: "0",
      },
      expected: {
        aPerBRatioMessage: "1.00000000",
        bPerARatioMessage: "1.00000000",
        aPerBRatioProjectedMessage: "1.00000000",
        bPerARatioProjectedMessage: "1.00000000",
        shareOfPool: "0.99%",
        state: PoolState.VALID_INPUT,
      },
    },
    {
      poolExternal: "1000000000000000000000000000",
      poolNative: "1000000000000000000000000000",
      poolUnits: "1000000000000000000000000000",
      addedExternal: "10000000",
      addedNative: "10000000",
      externalSymbol: "atk",
      nativeSymbol: "rowan",
      preexistingLiquidity: {
        native: "500000000000000000000000000",
        external: "500000000000000000000000000",
        units: "500000000000000000000000000",
      },
      expected: {
        aPerBRatioMessage: "1.00000000",
        bPerARatioMessage: "1.00000000",
        aPerBRatioProjectedMessage: "1.00000000",
        bPerARatioProjectedMessage: "1.00000000",
        shareOfPool: "50.50%",
        state: PoolState.VALID_INPUT,
      },
    },
    // Test for small decimals coin
    {
      poolExternal: "1000000000000000",
      poolNative: "1000000000000000000000000000",
      poolUnits: "1000000000000000000000000000",
      addedExternal: "10000",
      addedNative: "10000",
      externalSymbol: "ctest",
      nativeSymbol: "rowan",
      preexistingLiquidity: {
        native: "0",
        external: "0",
        units: "0",
      },
      expected: {
        aPerBRatioMessage: "1.00000000",
        bPerARatioMessage: "1.00000000",
        aPerBRatioProjectedMessage: "1.00000000",
        bPerARatioProjectedMessage: "1.00000000",
        shareOfPool: "33.33%", // This seems puzzling to me...
        state: PoolState.VALID_INPUT,
      },
    },
    {
      poolExternal: "1000000000000000000000000000",
      poolNative: "1000000000000000000000000000",
      poolUnits: "1000000000000000000000000000",
      addedExternal: "10000000",
      addedNative: "10000000",
      externalSymbol: "atk",
      nativeSymbol: "rowan",
      preexistingLiquidity: {
        native: "1000000000000000000000000000",
        external: "1000000000000000000000000000",
        units: "1000000000000000000000000000",
      },
      expected: {
        aPerBRatioMessage: "1.00000000",
        bPerARatioMessage: "1.00000000",
        aPerBRatioProjectedMessage: "1.00000000",
        bPerARatioProjectedMessage: "1.00000000",
        shareOfPool: "100.00%",
        state: PoolState.VALID_INPUT,
      },
    },
    {
      poolExternal: "1000000000000000000000000000",
      poolNative: "1000000000000000000000000000",
      poolUnits: "1000000000000000000000000000",
      addedExternal: "10000000",
      addedNative: "10000000",
      externalSymbol: "atk",
      nativeSymbol: "rowan",
      expected: {
        aPerBRatioMessage: "1.00000000",
        bPerARatioMessage: "1.00000000",
        aPerBRatioProjectedMessage: "1.00000000",
        bPerARatioProjectedMessage: "1.00000000",
        shareOfPool: "0.99%",
        state: PoolState.VALID_INPUT,
      },
    },
    {
      poolExternal: "20000000000000000000000000",
      poolNative: "10000000000000000000000000",
      poolUnits: "20000000000000000000000000",
      addedExternal: "10000000",
      addedNative: "40000000",
      externalSymbol: "atk",
      nativeSymbol: "rowan",
      expected: {
        aPerBRatioMessage: "2.00000000",
        bPerARatioMessage: "0.50000000",
        aPerBRatioProjectedMessage: "0.60000000",
        bPerARatioProjectedMessage: "1.66666667",
        shareOfPool: "55.56%",
        state: PoolState.VALID_INPUT,
      },
    },
    {
      poolExternal: "40000000000000000000000000",
      poolNative: "10000000000000000000000000",
      poolUnits: "40000000000000000000000000",
      addedExternal: "10000000",
      addedNative: "40000000",
      externalSymbol: "atk",
      nativeSymbol: "rowan",
      expected: {
        aPerBRatioMessage: "4.00000000",
        bPerARatioMessage: "0.25000000",
        aPerBRatioProjectedMessage: "1.00000000",
        bPerARatioProjectedMessage: "1.00000000",
        shareOfPool: "50.00%",
        state: PoolState.VALID_INPUT,
      },
    },

    {
      poolExternal: "100000000000000000000000000",
      poolNative: "100000000000000000000000000",
      poolUnits: "10000000000000000000000000000000000000000000",
      addedExternal: "100000000000",
      addedNative: "1",
      externalSymbol: "atk",
      nativeSymbol: "rowan",
      expected: {
        aPerBRatioMessage: "1.00000000", // 100000000 / 100000000
        bPerARatioMessage: "1.00000000",
        aPerBRatioProjectedMessage: "1000.99998999", // 100100000000/100000001
        bPerARatioProjectedMessage: "0.00099900",
        shareOfPool: "33.31%",
        state: PoolState.VALID_INPUT,
      },
    },
    {
      poolExternal: "100000000000000000000000000",
      poolNative: "100000000000000000000000000",
      poolUnits: "10000000000000000000000000000000000000000000",
      addedExternal: "100000000000000", // more than balance
      addedNative: "1",
      externalSymbol: "atk",
      nativeSymbol: "rowan",
      expected: {
        aPerBRatioMessage: "1.00000000", // 100000000 / 100000000
        bPerARatioMessage: "1.00000000",
        aPerBRatioProjectedMessage: "1000000.98999999", // 100100000000/100000001
        bPerARatioProjectedMessage: "0.00000100",
        shareOfPool: "33.55%",
        state: PoolState.INSUFFICIENT_FUNDS,
      },
    },
  ];
  describe("ratios", () => {
    ratios.forEach(
      (
        {
          only = false,
          skip = false,
          poolExternal,
          poolNative,
          poolUnits,
          addedExternal,
          addedNative,
          externalSymbol,
          nativeSymbol,
          preexistingLiquidity,
          expected,
        },
        index,
      ) => {
        const tester = only ? test.only : skip ? test.skip : test;
        tester(`Ratios #${index + 1}`, () => {
          balances.value = [
            AssetAmount(ATK, "100000000000000000000000000000"),
            AssetAmount(ROWAN, "100000000000000000000000000000"),
            AssetAmount(CTEST, "100000000000000000000000000000"),
          ];
          liquidityProvider.value = !preexistingLiquidity
            ? null
            : LiquidityProvider(
                ATK,
                Amount(preexistingLiquidity.units),
                akasha.address,
                Amount(preexistingLiquidity.native),
                Amount(preexistingLiquidity.external),
              );

          poolFinder.mockImplementation(() => {
            const pool = Pool(
              AssetAmount(Asset.get(externalSymbol), poolExternal),
              AssetAmount(ROWAN, poolNative),
              Amount(poolUnits),
            );

            return ref(pool) as Ref<Pool>;
          });

          tokenAAmount.value = addedExternal;
          tokenBAmount.value = addedNative;
          tokenASymbol.value = externalSymbol;
          tokenBSymbol.value = nativeSymbol;
          expect(aPerBRatioMessage.value).toBe(expected.aPerBRatioMessage);
          expect(bPerARatioMessage.value).toBe(expected.bPerARatioMessage);
          expect(aPerBRatioProjectedMessage.value).toBe(
            expected.aPerBRatioProjectedMessage,
          );
          expect(bPerARatioProjectedMessage.value).toBe(
            expected.bPerARatioProjectedMessage,
          );
          expect(shareOfPoolPercent.value).toBe(expected.shareOfPool);
          expect(state.value).toBe(expected.state);
        });
      },
    );
  });

  test("poolCalculator ratio messages", () => {
    poolFinder.mockImplementation(
      () =>
        ref(
          Pool(
            AssetAmount(ATK, "2000000000000000000000000"),
            AssetAmount(ROWAN, "1000000000000000000000000"),
          ),
        ) as Ref<Pool>,
    );

    tokenAAmount.value = "100000";
    tokenBAmount.value = "500000";
    tokenASymbol.value = "atk";
    tokenBSymbol.value = "rowan";

    expect(aPerBRatioMessage.value).toBe("2.00000000");
    expect(bPerARatioMessage.value).toBe("0.50000000");
    expect(aPerBRatioProjectedMessage.value).toBe("1.40000000");
    expect(bPerARatioProjectedMessage.value).toBe("0.71428571");
    expect(shareOfPoolPercent.value).toBe("13.49%");
  });

  test("poolCalculator with preexisting pool", () => {
    // Pool exists with 1001000 preexisting units 1000 of which are owned by this lp
    poolFinder.mockImplementation(
      () =>
        ref(
          Pool(
            AssetAmount(ATK, "1000000000000000000000000"),
            AssetAmount(ROWAN, "1000000000000000000000000"),
            Amount("1000000000000000000000000"),
          ),
        ) as Ref<Pool>,
    );

    // Liquidity provider already owns 1000 pool units (1000000 from another investor)
    liquidityProvider.value = LiquidityProvider(
      ATK,
      Amount("500000000000000000000000"),
      akasha.address,
      Amount("500000000000000000000000"),
      Amount("500000000000000000000000"),
    );

    // Add 1000 of both tokens
    tokenAAmount.value = "500000";
    tokenBAmount.value = "500000";
    tokenASymbol.value = "atk";
    tokenBSymbol.value = "rowan";

    expect(aPerBRatioMessage.value).toBe("1.00000000");
    expect(bPerARatioMessage.value).toBe("1.00000000");
    expect(aPerBRatioProjectedMessage.value).toBe("1.00000000");
    expect(bPerARatioProjectedMessage.value).toBe("1.00000000");
    // New shareOfPoolPercent for liquidity provider (inc prev liquidity)
    //2000/1002000 = 0.001996007984031936 so roughtly 0.2%
    expect(shareOfPoolPercent.value).toBe("66.67%");

    // New pool units for liquidity provider (inc prev liquidity)
    expect(totalLiquidityProviderUnits.value).toBe("1000000000000000000000000");

    expect(totalPoolUnits.value).toBe("1500000000000000000000000");
  });

  test("poolCalculator with preexisting pool but no preexisting liquidity", () => {
    // Pool exists with 1001000 preexisting units 1000 of which are owned by this lp
    poolFinder.mockImplementation(
      () =>
        ref(
          Pool(
            AssetAmount(ATK, "1000000000000000000000000"),
            AssetAmount(ROWAN, "1000000000000000000000000"),
            Amount("1000000000000000000000000"),
          ),
        ) as Ref<Pool>,
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
    expect(aPerBRatioProjectedMessage.value).toBe("1.00000000");
    expect(bPerARatioProjectedMessage.value).toBe("1.00000000");

    // New shareOfPoolPercent for liquidity provider (inc prev liquidity)
    //2000/1002000 = 0.001996007984031936 so roughtly 0.2%
    expect(shareOfPoolPercent.value).toBe("33.33%");

    // New pool units for liquidity provider (inc prev liquidity)
    expect(totalLiquidityProviderUnits.value).toBe("500000000000000000000000");

    expect(totalPoolUnits.value).toBe("1500000000000000000000000");
  });

  test("Can handle division by zero", () => {
    liquidityProvider.value = LiquidityProvider(ATK, ZERO, "", ZERO, ZERO);
    tokenAAmount.value = "0";
    tokenBAmount.value = "0";
    tokenASymbol.value = "atk";
    tokenBSymbol.value = "rowan";
    expect(state.value).toBe(PoolState.ZERO_AMOUNTS_NEW_POOL);
    expect(aPerBRatioMessage.value).toBe("N/A");
    expect(bPerARatioMessage.value).toBe("N/A");
    expect(aPerBRatioProjectedMessage.value).toBe("N/A");
    expect(bPerARatioProjectedMessage.value).toBe("N/A");
    expect(shareOfPoolPercent.value).toBe("< 0.01%");
  });

  test("Don't allow rowan === 0 when creating new pool", () => {
    tokenAAmount.value = "1000";
    tokenBAmount.value = "0";
    tokenASymbol.value = "atk";
    tokenBSymbol.value = "rowan";
    expect(state.value).toBe(PoolState.ZERO_AMOUNTS_NEW_POOL);
    expect(aPerBRatioMessage.value).toBe("N/A");
    expect(bPerARatioMessage.value).toBe("N/A");
    expect(aPerBRatioProjectedMessage.value).toBe("N/A");
    expect(bPerARatioProjectedMessage.value).toBe("N/A");
  });

  test("Don't allow external token === 0 when creating new pool", () => {
    tokenAAmount.value = "0";
    tokenBAmount.value = "1000";
    tokenASymbol.value = "atk";
    tokenBSymbol.value = "rowan";
    expect(state.value).toBe(PoolState.ZERO_AMOUNTS_NEW_POOL);
    expect(aPerBRatioMessage.value).toBe("N/A");
    expect(bPerARatioMessage.value).toBe("N/A");
    expect(aPerBRatioProjectedMessage.value).toBe("N/A");
    expect(bPerARatioProjectedMessage.value).toBe("N/A");
  });

  test("Allow rowan === 0 when adding to preExistingPool", () => {
    balances.value = [
      AssetAmount(ATK, "1000000000000000000000"),
      AssetAmount(ROWAN, "1000000000000000000000"),
    ];
    poolFinder.mockImplementation(
      () =>
        ref(
          Pool(
            AssetAmount(ATK, "1000000000000000000000000"),
            AssetAmount(ROWAN, "1000000000000000000000000"),
          ),
        ) as Ref<Pool>,
    );
    tokenAAmount.value = "1000";
    tokenBAmount.value = "0";
    tokenASymbol.value = "atk";
    tokenBSymbol.value = "rowan";
    expect(state.value).toBe(PoolState.VALID_INPUT);
    expect(aPerBRatioMessage.value).toBe("1.00000000");
    expect(bPerARatioMessage.value).toBe("1.00000000");
    expect(aPerBRatioProjectedMessage.value).toBe("1.00100000");
    expect(bPerARatioProjectedMessage.value).toBe("0.99900100");
  });

  test("Allow external token === 0 when adding to preExistingPool", () => {
    balances.value = [
      AssetAmount(ATK, "1000000000000000000000"),
      AssetAmount(ROWAN, "1000000000000000000000"),
    ];
    poolFinder.mockImplementation(
      () =>
        ref(
          Pool(
            AssetAmount(ATK, "1000000000000000000000000"),
            AssetAmount(ROWAN, "1000000000000000000000000"),
          ),
        ) as Ref<Pool>,
    );
    tokenAAmount.value = "0";
    tokenBAmount.value = "1000";
    tokenASymbol.value = "atk";
    tokenBSymbol.value = "rowan";
    expect(state.value).toBe(PoolState.VALID_INPUT);
    expect(aPerBRatioMessage.value).toBe("1.00000000");
    expect(bPerARatioMessage.value).toBe("1.00000000");
    expect(aPerBRatioProjectedMessage.value).toBe("0.99900100");
    expect(bPerARatioProjectedMessage.value).toBe("1.00100000");
  });

  test("insufficient funds", () => {
    balances.value = [
      AssetAmount(ATK, "100000000000000000000"),
      AssetAmount(ROWAN, "100000000000000000000"),
    ];
    tokenAAmount.value = "1000";
    tokenBAmount.value = "500";
    tokenASymbol.value = "atk";
    tokenBSymbol.value = "rowan";

    expect(state.value).toBe(PoolState.INSUFFICIENT_FUNDS);
  });

  test("valid funds below limit", () => {
    balances.value = [
      AssetAmount(ATK, "1000000000000000000000"),
      AssetAmount(ROWAN, "500000000000000000000"),
    ];
    tokenAAmount.value = "999";
    tokenBAmount.value = "499";
    tokenASymbol.value = "atk";
    tokenBSymbol.value = "rowan";
    expect(state.value).toBe(PoolState.VALID_INPUT);
  });

  test("valid funds at limit", () => {
    balances.value = [
      AssetAmount(ATK, "1000000000000000000000"),
      AssetAmount(ROWAN, "500000000000000000000"),
    ];
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
