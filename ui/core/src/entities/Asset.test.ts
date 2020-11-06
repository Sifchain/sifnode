import { createAsset } from "./Asset";
import { ChainId } from "./ChainId";

test("it should represent an asset", () => {
  const USD = createAsset("USDC", 2, "USD Coin", ChainId.ETHEREUM);
  expect(USD).toEqual({
    decimals: 2,
    symbol: "USDC",
    name: "USD Coin",
    chainId: 1,
  });
});
