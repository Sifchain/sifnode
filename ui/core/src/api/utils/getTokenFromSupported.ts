import { parseAssets, AssetConfig } from "./parseConfig";

export function getTokenFromSupported(assets: any[], tokenSymbol: string) {
  const supportedTokens = parseAssets(assets);

  const asset = supportedTokens.find(
    ({ symbol }) => symbol.toUpperCase() === tokenSymbol.toUpperCase()
  );

  if (!asset) throw new Error(`${tokenSymbol} not returned`);

  return asset;
}
