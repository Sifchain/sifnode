import localethereumassets from "../../assets.ethereum.localnet.json";
import localsifassets from "../../assets.sifchain.localnet.json";

import { parseAssets } from "../../utils/parseConfig";
import { Asset, IAssetAmount } from "../../entities";

const assets = [...localethereumassets.assets, ...localsifassets.assets];

export function getTestingToken(tokenSymbol: string) {
  const supportedTokens = parseAssets(assets as any[]).map(asset => {
    Asset.set(asset.symbol, asset);
    return asset;
  });

  const asset = supportedTokens.find(
    ({ symbol }) => symbol.toUpperCase() === tokenSymbol.toUpperCase()
  );

  if (!asset) throw new Error(`${tokenSymbol} not returned`);

  return asset;
}

export function getTestingTokens(tokens: string[]) {
  return tokens.map(getTestingToken);
}

export function getBalance(balances: IAssetAmount[], symbol: string) {
  const bal = balances.find(
    ({ asset }) => asset.symbol.toUpperCase() === symbol.toUpperCase()
  );
  if (!bal) throw new Error("Symbol not found in balances");
  return bal;
}
