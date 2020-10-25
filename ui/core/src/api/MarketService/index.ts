import { Asset, AssetAmount, Pair } from "../../entities";

export type MarketServiceContext = {
  loadAssets: () => Promise<Asset[]>;
  fetchMarketData: () => Promise<{ name: string; value: number }[][]>;
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
  fetchMarketData,
}: MarketServiceContext) {
  const pairs = new Map<string, Pair>();

  async function generatePairs() {
    await loadAssets();
    const data = await fetchMarketData();

    data.map(([amount1, amount2]) => {
      const asset1 = Asset.get(amount1.name);
      const asset2 = Asset.get(amount2.name);
      const pair = Pair(
        AssetAmount(asset1, amount1.value),
        AssetAmount(asset2, amount2.value)
      );
      pairs.set(pair.symbol(), pair);
    });
  }
  const pairsGenerated = makeQuerablePromise(generatePairs());

  return {
    find(asset1: Asset | string, asset2: Asset | string) {
      if (!pairsGenerated.isResolved()) return null;
      const key = [asset1, asset2].map(toAssetSymbol).join("_");
      return pairs.get(key);
    },
  };
}
