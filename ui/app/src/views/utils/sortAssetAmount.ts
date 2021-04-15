import { IAssetAmount, Asset } from "ui-core";

export function sortAssetAmount(
  assetAmounts: {
    amount: IAssetAmount | null | undefined;
    asset: Asset;
  }[],
): {
  amount: IAssetAmount | null | undefined;
  asset: Asset;
}[] {
  return assetAmounts
    .sort((a, b) => {
      // Sort alphabetically
      if (a.asset.symbol < b.asset.symbol) {
        return -1;
      }
      if (a.asset.symbol > b.asset.symbol) {
        return 1;
      }
      return 0;
    })
    .sort((a, b) => {
      // Next sort by balance
      if (!b.amount || !a.amount) {
        return 0;
      }
      return b.amount.greaterThanOrEqual(a.amount) ? 1 : -1;
      // return b.amount.subtract(a.amount).toBigInt();
    })
    .sort((a) => {
      // Finally, sort and move rowan, erowan to the top
      if (["rowan", "erowan"].includes(a.asset.symbol.toLowerCase())) {
        return -1;
      } else {
        return 1;
      }
    });
}
