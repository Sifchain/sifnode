import { Asset, Token } from "../../entities";

export type TokenServiceContext = {
  getSupportedTokens: () => Promise<Token[]>;
  getSupportedAssets: () => Promise<Asset[]>;
};

export default function createTokenService({
  getSupportedTokens,
  getSupportedAssets,
}: TokenServiceContext) {
  // // define map to store all assets
  const _assets: Asset[] = [];

  // async fetch the supported tokens and store in the map
  async function loadAssets() {
    const tokens = await getSupportedTokens();
    const assets = await getSupportedAssets();
    _assets.push(...assets, ...tokens);
  }
  const assetsLoadedPromise = loadAssets();

  return {
    async getTopAssets() {
      await assetsLoadedPromise;
      return _assets;
    },

    getAsset(symbol: string): Asset | undefined {
      return Asset.get(symbol);
    },

    getImage(symbol: string): string | undefined {
      return Asset.get(symbol).imageUrl;
    },
  };
}
