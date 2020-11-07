import { Asset, AssetAmount, CompositePair, Pair } from "../../entities";

export type MarketServiceContext = {
  loadAssets: () => Promise<Asset[]>;
  getPools: () => Promise<Pair[]>;
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
  getPools,
}: MarketServiceContext) {
  const pairs = new Map<string, Pair>();

  async function generatePairs() {
    await loadAssets();
    const pools = await getPools();

    pools.map((pair) => {
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
