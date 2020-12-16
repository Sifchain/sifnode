import { Asset } from "../../entities";

export type TokenServiceContext = {
  assets: Asset[];
};

export default function createTokenService({ assets }: TokenServiceContext) {
  // // define map to store all assets
  let _assets: Asset[] = [];

  // async fetch the supported tokens and store in the map
  const cacheLoadedAssets = async () => {
    _assets = assets;
  };

  const cacheIsReady = cacheLoadedAssets();

  return {
    async getTopAssets() {
      await cacheIsReady;
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
