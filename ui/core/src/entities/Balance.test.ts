import { Balance } from "./Balance";
import { createAsset } from "./Asset";

const USD = createAsset(2, "USD", "US Dollar");

test("it should be able to handle whole integars", () => {
  const f = new Balance(USD, "10012");
  expect(f.toFixed(2)).toBe("100.12");
});
