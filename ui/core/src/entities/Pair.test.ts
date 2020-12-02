import JSBI from "jsbi";
import { Asset } from "./Asset";
import { AssetAmount } from "./AssetAmount";
import { Coin } from "./Coin";
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
  const ETH = Coin({
    decimals: 18,
    symbol: "eth",
    name: "Ethereum",
    network: Network.ETHEREUM,
  });
  const RWN = Coin({
    decimals: 18,
    symbol: "rwn",
    name: "Rowan",
    network: Network.SIFCHAIN,
  });

  test("contains()", () => {
    const pair = Pair(AssetAmount(ATK, "10"), AssetAmount(BTK, "10"));

    expect(pair.contains(ATK)).toBe(true);
    expect(pair.contains(BTK)).toBe(true);
    expect(pair.contains(RWN)).toBe(false);
  });
  test("other asset", () => {
    const pair = Pair(AssetAmount(ATK, "10"), AssetAmount(BTK, "10"));
    expect(pair.otherAsset(ATK).symbol).toBe("btk");
  });
  describe("when half", () => {
    const pair = Pair(AssetAmount(ATK, "5"), AssetAmount(BTK, "10"));

    test("pair has symbol", () => {
      expect(pair.symbol()).toEqual("atk_btk");
    });
  });
});
