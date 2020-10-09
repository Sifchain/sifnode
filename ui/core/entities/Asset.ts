export type Asset = {
  decimals: number;
  symbol: string;
  name: string;
};

export function createAsset(
  decimals: number,
  symbol: string,
  name: string
): Asset {
  return {
    decimals,
    symbol,
    name,
  };
}
