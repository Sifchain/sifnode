import { AssetAmount, IAssetAmount } from "../entities";
import { getTestingTokens } from "../test/utils/getTestingToken";
import { fromBaseUnits } from "../utils";
const [ETH] = getTestingTokens(["ETH"]);

export function toDerived(assetAmount: IAssetAmount) {
  return assetAmount.amount.multiply(fromBaseUnits("1", assetAmount.asset));
}

test("derived", () => {
  const oneeth = AssetAmount("eth", "1000000000000000000");
  expect(toDerived(oneeth).toString()).toEqual("1.000000000000000000");
});
