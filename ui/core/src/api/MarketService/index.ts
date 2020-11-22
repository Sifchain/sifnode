import { RWN } from "../../constants";
import { Asset, AssetAmount, Pool } from "../../entities";
import { Fraction } from "../../entities/fraction/Fraction";
import { SifUnSignedClient } from "../utils/SifClient";
import { RawPool } from "../utils/x/clp";

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

function processPool(poolData: RawPool) {
  const externalAssetTicker = poolData.external_asset.ticker;

  return Pool(
    AssetAmount(RWN, poolData.native_asset_balance),
    AssetAmount(
      Asset.get(externalAssetTicker),
      poolData.external_asset_balance
    ),
    new Fraction(poolData.pool_units)
  );
}

export default function createMarketService({
  loadAssets,
  sifApiUrl,
}: MarketServiceContext) {
  const sifClient = new SifUnSignedClient(sifApiUrl);
  const poolMap = new Map<string, Pool>();

  async function initialize() {
    await loadAssets();
    instance.getPools();
  }

  const pairsGenerated = makeQuerablePromise(initialize());

  const instance = {
    async getPools() {
      const rawPools = await sifClient.getPools();
      const pools = rawPools.map(processPool);

      pools.forEach((pool) => {
        poolMap.set(pool.symbol(), pool);
      }, poolMap);

      return pools;
    },
    find(asset1: Asset | string, asset2: Asset | string) {
      if (!pairsGenerated.isResolved()) return null;
      const key = [asset1, asset2]
        .map(toAssetSymbol)
        .sort()
        .join("_");
      return poolMap.get(key) ?? null;
    },
  };

  return instance;
}
