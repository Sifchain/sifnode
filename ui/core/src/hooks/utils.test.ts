import { ATK, BTK } from "../constants";
import { AssetAmount, Pair, Coin, Network } from "../entities";
import { assetPriceMessage } from "./utils";

describe("assets with decimals", () => {
  const ASSETS = {
    atk: Coin({
      symbol: "catk",
      name: "Atk",
      network: Network.SIFCHAIN,
      decimals: 18,
    }),
    btk: Coin({
      symbol: "cbtk",
      name: "Btk",
      network: Network.SIFCHAIN,
      decimals: 18,
    }),
  };
  test("assetPriceMessage", () => {
    const msg = assetPriceMessage(
      AssetAmount(ASSETS.atk, "100"),
      Pair(
        AssetAmount(ASSETS.atk, "1000000"),
        AssetAmount(ASSETS.btk, "2000000")
      ),
      4
    );

    expect(msg).toBe("0.5001 CATK per CBTK");
  });

  test("with zero amounts message should be nothing", () => {
    const msg = assetPriceMessage(
      AssetAmount(ASSETS.atk, "0"),
      Pair(AssetAmount(ASSETS.atk, "1"), AssetAmount(ASSETS.btk, "1")),
      4
    );
    expect(msg).toBe("");
  });
});
describe("assets with zero decimals", () => {
  const ASSETS = {
    atk: Coin({
      symbol: "catk",
      name: "Atk",
      network: Network.SIFCHAIN,
      decimals: 0,
    }),
    btk: Coin({
      symbol: "cbtk",
      name: "Btk",
      network: Network.SIFCHAIN,
      decimals: 0,
    }),
  };
  test("with 1 as an amount", () => {
    const msg = assetPriceMessage(
      AssetAmount(ASSETS.atk, "1"),
      Pair(
        AssetAmount(ASSETS.atk, "1000000"),
        AssetAmount(ASSETS.btk, "1000000")
      )
    );
    expect(msg).toBe("1 CATK per CBTK");
  });

  test("with 12 as an amount", () => {
    const msg = assetPriceMessage(
      AssetAmount(ASSETS.atk, "12"),
      Pair(
        AssetAmount(ASSETS.atk, "1000000"),
        AssetAmount(ASSETS.btk, "1000000")
      ),
      4
    );
    expect(msg).toBe("1.0000 CATK per CBTK");
  });
});
