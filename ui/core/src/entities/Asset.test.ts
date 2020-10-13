import { createAsset } from "./Asset";

test("it should represent and asset", () => {
  const USD = createAsset(2, "USD", "US Dollar");
  expect(USD).toEqual({ decimals: 2, symbol: "USD", name: "US Dollar" });
});
