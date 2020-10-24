import { AssetAmount } from "./AssetAmount";
import { Network } from "./Network";
import { Coin } from "./Coin";

const USD = Coin({
  symbol: "USD",
  decimals: 2,
  name: "US Dollar",
  network: Network.ETHEREUM,
});

const ETH = Coin({
  symbol: "ETH",
  decimals: 18,
  name: "Ethereum",
  network: Network.ETHEREUM,
});

test("it should be able to handle whole integars", () => {
  const f = new AssetAmount(USD, "10012");
  expect(f.toFixed(2)).toBe("100.12");
});

test("Shorthand", () => {
  expect(AssetAmount.n(USD, "100.12").toFixed()).toBe("100.12");
  expect(AssetAmount.n(ETH, "10.1234567").toFixed()).toBe(
    "10.123456700000000000"
  );
});
