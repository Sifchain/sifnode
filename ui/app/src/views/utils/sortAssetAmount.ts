import { IAssetAmount, Asset } from "ui-core";
import { format } from "ui-core/src/utils/format";

export function sortAssetAmount<
  T extends {
    amount: IAssetAmount | null | undefined;
    asset: Asset;
  }[]
>(assetAmounts: T): T {
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

      if (!b.amount?.amount) {
        return -1;
      }
      if (!a.amount?.amount) {
        return 1;
      }

      // #TODO - TD - There is likely a much better way to do this an entire sort
      //              But I couldn't figure it out #refactor
      //              Asset balances needed to be sorted once precision is applied

      const aValue = format(a.amount.amount, a.asset, {
        mantissa: 18,
      });

      const bValue = format(b.amount.amount, b.asset, {
        mantissa: 18,
      });

      return Number(bValue) >= Number(aValue) ? 1 : -1;
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
