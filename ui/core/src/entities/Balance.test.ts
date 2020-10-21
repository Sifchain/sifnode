import { Balance } from "./Balance";
import { createAsset } from "./Asset";
import { ChainId } from "./ChainId";

const USD = createAsset("USD", 2, "US Dollar", ChainId.ETHEREUM);
const ETH = createAsset("ETH", 18, "Ethereum", ChainId.ETHEREUM);

test("it should be able to handle whole integars", () => {
  const f = new Balance(USD, "10012");
  expect(f.toFixed(2)).toBe("100.12");
});

test("Shorthand", () => {
  expect(Balance.n(USD, "100.12").toFixed()).toBe("100.12");
  expect(Balance.n(ETH, "10.1234567").toFixed()).toBe("10.123456700000000000");
});
