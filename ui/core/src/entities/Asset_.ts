import { Network } from "./Network";

export type IAsset = {
  address: string;
  decimals: number;
  imageUrl?: string;
  name: string;
  network: Network;
  symbol: string;
  label: string;
};
type ReadonlyAsset = Readonly<IAsset>;
const assetMap = new Map<string, ReadonlyAsset>();

function isAsset(value: any): value is IAsset {
  return (
    typeof value?.symbol === "string" && typeof value?.decimals === "number"
  );
}

export function Asset(assetOrSymbol: IAsset | string): ReadonlyAsset {
  // If it is an asset then cache it and return it
  if (isAsset(assetOrSymbol)) {
    assetMap.set(
      assetOrSymbol.symbol.toLowerCase(),
      assetOrSymbol as ReadonlyAsset
    );
    return assetOrSymbol;
  }

  // Return it from cache
  const found = assetOrSymbol
    ? assetMap.get(assetOrSymbol.toLowerCase())
    : false;
  if (!found) {
    throw new Error(
      `Attempt to retrieve the asset with key "${assetOrSymbol}" before it had been cached.`
    );
  }

  return found;
}
