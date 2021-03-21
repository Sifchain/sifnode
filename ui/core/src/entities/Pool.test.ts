import { getTestingTokens } from "../test/utils/getTestingToken";
import { Asset } from "./Asset";
import { AssetAmount } from "./AssetAmount";
import { Network } from "./Network";
import { Pool, CompositePool } from "./Pool";

describe("Pool", () => {
  const [ATK, BTK, ETH, ROWAN] = getTestingTokens([
    "ATK",
    "BTK",
    "ETH",
    "ROWAN",
  ]);

  // TODO: Confirm with Blockscience
  // x = Sent Asset Amount, X = Sent Asset Pool Balance, Y = Received Asset Pool Balance

  // Liquidity Fee = ( x^2 * Y ) / ( x + X )^2
  // Trade Slip = x * (2*X + x) / (X * X)
  // Swap Result = ( x * X * Y ) / ( x + X )^2

  test("It calculates the correct swap amount", () => {
    const pair = Pool(
      AssetAmount(ATK, "10000000000000000000"),
      AssetAmount(BTK, "10000000000000000000"),
    );

    expect(
      pair.calcSwapResult(AssetAmount(ATK, "1000000000000000000")).toString(),
    ).toEqual(
      "826446280991735537 BTK",
      // 0.826446280991735537 = (1 * 10 * 10) / ((1 + 10) * (1 + 10));
    );

    expect(
      pair.calcSwapResult(AssetAmount(BTK, "1000000000000000000")).toString(),
    ).toEqual(
      "826446280991735537 ATK",
      // 0.826446280991735537 = (1 * 10 * 10) / ((1 + 10) * (1 + 10));
    );
  });

  test("calculate swap amount 1:1", () => {
    const pair = Pool(
      AssetAmount(ATK, "1000000000000000000000000000000"),
      AssetAmount(BTK, "1000000000000000000000000000000"),
    );

    expect(
      pair.calcSwapResult(AssetAmount(ATK, "1000000000000000000")).toString(),
    ).toEqual(
      "999999999998000000 BTK",
      // 0.826446280991735537 = (1 * 1000000000000 * 1000000000000) / ((1 + 1000000000000) * (1 + 1000000000000));
    );

    expect(
      pair.calcSwapResult(AssetAmount(BTK, "1000000000000000000")).toString(),
    ).toEqual(
      "999999999998000000 ATK",
      // 0.826446280991735537 = (1 * 1000000000000 * 1000000000000) / ((1 + 1000000000000) * (1 + 1000000000000));
    );
  });
  test("calculate swap amount 2:1", () => {
    const pair = Pool(
      AssetAmount(ATK, "2000000000000000000000000000000"),
      AssetAmount(BTK, "1000000000000000000000000000000"),
    );

    expect(
      pair.calcSwapResult(AssetAmount(ATK, "1000000000000000000")).toString(),
    ).toEqual(
      "499999999999500000 BTK",
      // 0.826446280991735537 = (1 * 1000000000000 * 1000000000000) / ((1 + 1000000000000) * (1 + 1000000000000));
    );

    expect(
      pair.calcSwapResult(AssetAmount(BTK, "1000000000000000000")).toString(),
    ).toEqual(
      "1999999999996000000 ATK",
      // 0.826446280991735537 = (1 * 1000000000000 * 1000000000000) / ((1 + 1000000000000) * (1 + 1000000000000));
    );
  });

  test("swap of 0", () => {
    const pair = Pool(
      AssetAmount(ATK, "10000000000000000000"),
      AssetAmount(BTK, "10000000000000000000"),
    );

    expect(pair.calcSwapResult(AssetAmount(ATK, "0")).toString()).toEqual(
      "0 BTK",
    );
  });

  test("Reverse swap", () => {
    const pair = Pool(
      AssetAmount(ATK, "1000000000000000000000000"),
      AssetAmount(BTK, "1000000000000000000000000"),
    );

    expect(
      pair
        .calcReverseSwapResult(AssetAmount(BTK, "100000000000000000000"))
        .toString(),
    ).toEqual("100020005001400420132 ATK");
  });

  test("Reverse swap of 0", () => {
    const pair = Pool(
      AssetAmount(ATK, "10000000000000000000"),
      AssetAmount(BTK, "10000000000000000000"),
    );

    expect(
      pair.calcReverseSwapResult(AssetAmount(BTK, "0")).toString(),
    ).toEqual("0 ATK");
  });

  test("Cannot calulate swap result for an asset that does not exist within the pair", () => {
    const pair = Pool(
      AssetAmount(ATK, "10000000000000000000"),
      AssetAmount(BTK, "10000000000000000000"),
    );

    expect(() => {
      pair.calcSwapResult(AssetAmount(ROWAN, "10000000000000000000"));
    }).toThrow();
  });

  test("contains()", () => {
    const pair = Pool(AssetAmount(ATK, "10"), AssetAmount(BTK, "10"));
    expect(pair.contains(ATK)).toBe(true);
    expect(pair.contains(BTK)).toBe(true);
    expect(pair.contains(ROWAN)).toBe(false);
  });

  describe("when half", () => {
    const pair = Pool(
      AssetAmount(ATK, "5000000000000000000"),
      AssetAmount(BTK, "10000000000000000000"),
    );

    test("pair has symbol", () => {
      expect(pair.symbol()).toEqual("atk_btk");
    });

    test("calcSwapResult should be 1388..", () => {
      expect(
        pair.calcSwapResult(AssetAmount(ATK, "1000000000000000000")).toString(),
      ).toEqual(
        "1388888888888888889 BTK",
        // 1.388888888888888889 = (1 * 5 * 10) / ((1 + 5) * (1 + 5));
      );
      expect(
        pair.calcSwapResult(AssetAmount(BTK, "1000000000000000000")).toString(),
      ).toEqual(
        "413223140495867769 ATK",
        // 0.413223140495867769 = (1 * 10 * 5) / ((1 + 10) * (1 + 10));
      );
    });
  });

  describe("poolUnits", () => {
    test("poolUnits", () => {
      const pool = Pool(
        AssetAmount(ATK, "1000000"),
        AssetAmount(BTK, "1000000"),
      );

      expect(pool.poolUnits.toString()).toBe("1000000");
    });

    test("addLiquidity:calculate pool units", () => {
      const pool = Pool(
        AssetAmount(ATK, "1000000"),
        AssetAmount(BTK, "1000000"),
      );
      const [units, lpunits] = pool.calculatePoolUnits(
        AssetAmount(ATK, "10000"),
        AssetAmount(BTK, "14000"),
      );
      expect(units.toString()).toBe("1011953");
      expect(lpunits.divide(units).multiply("10000").toString()).toBe("118");
    });
  });

  describe("CompositePool", () => {
    test("Cannot create composite pair with pairs that have no shared asset", () => {
      const pair1 = Pool(
        AssetAmount(ATK, "1000000"),
        AssetAmount(BTK, "1000000"),
      );

      const pair2 = Pool(
        AssetAmount(ROWAN, "1000000"),
        AssetAmount(ETH, "1000000"),
      );
      expect(() => {
        CompositePool(pair1, pair2);
      }).toThrow();
    });

    test("CompositePool contains", () => {
      const pair1 = Pool(
        AssetAmount(ATK, "1000000000000"),
        AssetAmount(ROWAN, "1000000000000"),
      );

      const pair2 = Pool(
        AssetAmount(ROWAN, "1000000000000"),
        AssetAmount(BTK, "1000000000000"),
      );

      const compositePool = CompositePool(pair1, pair2);

      expect(compositePool.contains(BTK)).toBe(true);
      expect(compositePool.contains(ATK)).toBe(true);
      expect(compositePool.contains(ETH)).toBe(false);
    });

    test("CompositePool getAmount()", () => {
      const pair1 = Pool(
        AssetAmount(ATK, "2000000000000"),
        AssetAmount(ROWAN, "1000000000000"),
      );

      const pair2 = Pool(
        AssetAmount(ROWAN, "1000000000000"),
        AssetAmount(BTK, "1000000000000"),
      );

      const compositePool = CompositePool(pair1, pair2);

      expect(compositePool.getAmount("atk").toString()).toBe(
        "2000000000000 ATK",
      );
      expect(compositePool.getAmount(ATK).toString()).toBe("2000000000000 ATK");
      expect(compositePool.getAmount("btk").toString()).toBe(
        "1000000000000 BTK",
      );

      expect(() => {
        compositePool.getAmount("rowan");
      }).toThrow();
    });

    test("CompositePool does two swaps", () => {
      const pair1 = Pool(
        AssetAmount(ATK, "1000000000000000000000000000000"),
        AssetAmount(ROWAN, "1000000000000000000000000000000"),
      );

      const pair2 = Pool(
        AssetAmount(ROWAN, "1000000000000000000000000000000"),
        AssetAmount(BTK, "1000000000000000000000000000000"),
      );

      const compositePool = CompositePool(pair1, pair2);

      const inputAmount = AssetAmount(ATK, "10000000000000000000");

      const compositeSwapAmount = compositePool.calcSwapResult(inputAmount);

      expect(compositeSwapAmount.toString()).toEqual("9999999999600000000 BTK"); // Adjustment for fees

      const output = pair2
        .calcSwapResult(pair1.calcSwapResult(inputAmount))
        .toString();

      expect(compositeSwapAmount.toString()).toEqual(output.toString()); // Adjustment for fees
    });

    test("CompositePool 2:1", () => {
      const pair1 = Pool(
        AssetAmount(ATK, "2000000000000000000000000000000"),
        AssetAmount(ROWAN, "1000000000000000000000000000000"),
      );

      const pair2 = Pool(
        AssetAmount(ROWAN, "1000000000000000000000000000000"),
        AssetAmount(BTK, "1000000000000000000000000000000"),
      );

      const compositePool = CompositePool(pair1, pair2);
      const inputAmount = AssetAmount(ATK, "1000000000000000000");
      const compositeSwapAmount = compositePool.calcSwapResult(inputAmount);

      expect(compositeSwapAmount.toString()).toEqual("499999999999000000 BTK"); // Adjustment for fees
    });

    test("copmosite pair reverseswap", () => {
      const TOKENS = {
        atk: Asset({
          decimals: 18,
          symbol: "atk",
          label: "ATK",
          name: "AppleToken",
          address: "123",
          network: Network.ETHEREUM,
        }),
        btk: Asset({
          decimals: 18,
          symbol: "btk",
          label: "BTK",
          name: "BananaToken",
          address: "1234",
          network: Network.ETHEREUM,
        }),
        rowan: Asset({
          decimals: 18,
          symbol: "rowan",
          label: "ROWAN",
          name: "Rowan",
          address: "1234",
          network: Network.ETHEREUM,
        }),
        eth: Asset({
          decimals: 18,
          symbol: "eth",
          label: "ETH",
          name: "Ethereum",
          address: "1234",
          network: Network.ETHEREUM,
        }),
      };
      const pair1 = Pool(
        AssetAmount(TOKENS.atk, "2000000000000000000000000000000"),
        AssetAmount(TOKENS.rowan, "1000000000000000000000000000000"),
      );

      const pair2 = Pool(
        AssetAmount(TOKENS.btk, "1000000000000000000000000000000"),
        AssetAmount(TOKENS.rowan, "1000000000000000000000000000000"),
      );
      const compositePool = CompositePool(pair1, pair2);
      expect(
        compositePool
          .calcSwapResult(AssetAmount(TOKENS.atk, "100000000000000000000"))
          .toString(),
      ).toEqual("49999999990000000002 BTK");

      expect(
        compositePool
          .calcReverseSwapResult(
            AssetAmount(TOKENS.btk, "50000000000000000000"),
          )
          .toString(),
      ).toEqual("100000000019999999998 ATK");
    });
  });
});
