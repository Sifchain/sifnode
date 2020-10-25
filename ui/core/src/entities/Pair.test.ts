import { Asset } from "./Asset";
import { AssetAmount } from "./AssetAmount";
import { Network } from "./Network";
import { Pair } from "./Pair";
import { Token } from "./Token";

describe("Pair", () => {
  const ATK = Token({
    decimals: 6,
    symbol: "atk",
    name: "AppleToken",
    address: "123",
    network: Network.ETHEREUM,
  });
  const BTK = Token({
    decimals: 18,
    symbol: "btk",
    name: "BananaToken",
    address: "1234",
    network: Network.ETHEREUM,
  });
  const ETH = Token({
    decimals: 18,
    symbol: "eth",
    name: "Ethereum",
    address: "1234",
    network: Network.ETHEREUM,
  });

  // TODO: Confirm with Blockscience

  test("when equal it should be 1.0", () => {
    const pair = Pair(AssetAmount(ATK, "10"), AssetAmount(BTK, "10"));
    expect(pair.priceA().toFixed()).toEqual("1.000000000000000000");
  });

  describe("when half", () => {
    const pair = Pair(AssetAmount(ATK, "5"), AssetAmount(BTK, "10"));

    test("pair has symbol", () => {
      expect(pair.symbol()).toEqual("atk_btk");
    });

    test("priceA should be 2", () => {
      expect(pair.priceA().toFixed()).toEqual("2.000000000000000000"); // In terms of B
    });

    test("priceB should be 0.5", () => {
      expect(pair.priceB().toFixed()).toEqual("0.500000");
    });

    test("can match assets", () => {
      expect(pair.contains(ATK, BTK)).toBe(true);
      expect(pair.contains(BTK, ATK)).toBe(true);
      expect(pair.contains(ATK, ATK)).toBe(false);
      expect(pair.contains(ATK, ETH)).toBe(false);
    });
  });
});
