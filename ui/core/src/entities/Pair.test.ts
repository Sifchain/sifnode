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

  // TODO: Confirm with Blockscience

  test("when equal it should be 1.0", () => {
    const pair = Pair({ a: AssetAmount(ATK, "10"), b: AssetAmount(BTK, "10") });
    expect(pair.priceA().toFixed()).toEqual("1.000000000000000000");
  });

  describe("when half", () => {
    const pair = Pair({ a: AssetAmount(ATK, "5"), b: AssetAmount(BTK, "10") });
    test("priceA should be 2", () => {
      expect(pair.priceA().toFixed()).toEqual("2.000000000000000000"); // In terms of B
    });
    test("priceB should be 0.5", () => {
      expect(pair.priceB().toFixed()).toEqual("0.500000");
    });
  });
});
