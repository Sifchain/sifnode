import localethereumassets from "../../assets.ethereum.localnet.json";
import localsifassets from "../../assets.sifchain.localnet.json";

import { parseAssets } from "../../api/utils/parseConfig";

const assets = [...localethereumassets.assets, ...localsifassets.assets];

export function getTestingToken(tokenSymbol: string) {
  const supportedTokens = parseAssets(assets as any[]);

  const asset = supportedTokens.find(
    ({ symbol }) => symbol.toUpperCase() === tokenSymbol.toUpperCase()
  );

  if (!asset) throw new Error(`${tokenSymbol} not returned`);

  return asset;
}

export function getTestingTokens(tokens: string[]) {
  return tokens.map(getTestingToken);
}
