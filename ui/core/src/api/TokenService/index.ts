import { Asset, Token } from "../../entities";
import imageAssets from "./tokenImages";

export type TokenServiceContext = {
  getSupportedTokens: () => Promise<Token[]>;
  getSupportedAssets: () => Promise<Asset[]>;
};

export default function createTokenService({
  getSupportedTokens,
  getSupportedAssets,
}: TokenServiceContext) {
  // define map to store all assets
  const assetMap = new Map<string, Asset>();

  // update the map from the give assets
  function updateMap(givenAssets: Asset[]) {
    for (const asset of givenAssets) {
      if (!assetMap.has(asset.symbol)) {
        assetMap.set(asset.symbol, asset);
      }
    }
  }

  // async fetch the supported tokens and store in the map
  async function loadAssets() {
    const tokens = await getSupportedTokens();
    const assets = await getSupportedAssets();
    updateMap([...assets, ...tokens]);
  }
  const assetsLoadedPromise = loadAssets();

  return {
    async getTopAssets() {
      await assetsLoadedPromise;
      const topAssets: Asset[] = [];
      for (const [, asset] of assetMap) {
        topAssets.push(asset);
      }

      return topAssets;
    },

    getAsset(symbol: string): Asset | undefined {
      return assetMap.get(symbol);
    },

    getImage(symbol: string): string | undefined {
      return imageAssets[symbol as keyof typeof imageAssets];
    },
  };
}
