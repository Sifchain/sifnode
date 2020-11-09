import { RWN } from "../../constants";
import { Asset, AssetAmount, Pair } from "../../entities";
import { SifUnSignedClient } from "../SifService/SifClient";

export type MarketServiceContext = {
  loadAssets: () => Promise<Asset[]>;
  sifApiUrl: string;
  nativeAsset: Asset;
};

function toAssetSymbol(assetOrString: Asset | string) {
  return typeof assetOrString === "string"
    ? assetOrString
    : assetOrString.symbol;
}

function makeQuerablePromise<T>(promise: Promise<T>) {
  let isResolved = false;

  promise.then(() => {
    isResolved = true;
  });

  return {
    isResolved() {
      return isResolved;
    },
  };
}

export default function createMarketService({
  loadAssets,
  sifApiUrl,
}: MarketServiceContext) {
  const sifClient = new SifUnSignedClient(sifApiUrl);
  const pairs = new Map<string, Pair>();

  async function generatePairs() {
    await loadAssets();
    const pools = await sifClient.getPools();
    return pools.map((poolData) => {
      const externalAssetTicker = poolData.external_asset.ticker;

      const pair = Pair(
        AssetAmount(RWN, poolData.native_asset_balance),
        AssetAmount(
          Asset.get(externalAssetTicker),
          poolData.external_asset_balance
        )
      );

      pairs.set(pair.symbol(), pair);
    });
  }

  const pairsGenerated = makeQuerablePromise(generatePairs());

  return {
    find(asset1: Asset | string, asset2: Asset | string) {
      if (!pairsGenerated.isResolved()) return null;
      const key = [asset1, asset2].map(toAssetSymbol).join("_");
      return pairs.get(key) ?? null;
    },
  };
}
