import { AssetAmount, Pool, Network, Asset } from "../entities";
import { getTestingTokens } from "../test/utils/getTestingToken";
import { assetPriceMessage } from "./utils";

describe("assets with decimals", () => {
  const [CATK, CBTK] = getTestingTokens(["CATK", "CBTK"]);
  const ASSETS = {
    atk: CATK,
    btk: CBTK,
  };
  test("assetPriceMessage", () => {
    const msg = assetPriceMessage(
      AssetAmount(ASSETS.atk, "100"),
      Pool(
        AssetAmount(ASSETS.atk, "1000000"),
        AssetAmount(ASSETS.btk, "2000000"),
      ),
      4,
    );

    expect(msg).toBe("1.9996 cBTK per cATK");
  });

  test("with zero amounts message should be nothing", () => {
    const msg = assetPriceMessage(
      AssetAmount(ASSETS.atk, "0"),
      Pool(AssetAmount(ASSETS.atk, "1"), AssetAmount(ASSETS.btk, "1")),
      4,
    );
    expect(msg).toBe("");
  });
});
describe("assets with zero decimals", () => {
  const ASSETS = {
    atk: Asset({
      symbol: "catk",
      name: "Atk",
      network: Network.SIFCHAIN,
      decimals: 0,
      label: "cATK",
    }),
    btk: Asset({
      symbol: "cbtk",
      name: "Btk",
      network: Network.SIFCHAIN,
      decimals: 0,
      label: "cBTK",
    }),
  };
  test("with 1 as an amount", () => {
    const msg = assetPriceMessage(
      AssetAmount(ASSETS.atk, "1"),
      Pool(
        AssetAmount(ASSETS.atk, "1000000"),
        AssetAmount(ASSETS.btk, "1000000"),
      ),
      0,
    );
    expect(msg).toBe("1 cBTK per cATK");
  });

  test("with 12 as an amount", () => {
    const msg = assetPriceMessage(
      AssetAmount(ASSETS.atk, "12"),
      Pool(
        AssetAmount(ASSETS.atk, "1000000"),
        AssetAmount(ASSETS.btk, "1000000"),
      ),
      4,
    );
    expect(msg).toBe("1.0000 cBTK per cATK");
  });
});
