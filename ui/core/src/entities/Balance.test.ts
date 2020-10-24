import { Balance } from "./Balance";
import { ChainId } from "./ChainId";
import { Coin } from "./Coin";

const USD = Coin({
  symbol: "USD",
  decimals: 2,
  name: "US Dollar",
  chainId: ChainId.ETHEREUM,
});

const ETH = Coin({
  symbol: "ETH",
  decimals: 18,
  name: "Ethereum",
  chainId: ChainId.ETHEREUM,
});

test("it should be able to handle whole integars", () => {
  const f = new Balance(USD, "10012");
  expect(f.toFixed(2)).toBe("100.12");
});

test("Shorthand", () => {
  expect(Balance.n(USD, "100.12").toFixed()).toBe("100.12");
  expect(Balance.n(ETH, "10.1234567").toFixed()).toBe("10.123456700000000000");
});
