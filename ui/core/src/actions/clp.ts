import { Asset, AssetAmount } from "../entities";
import { ActionContext } from ".";
import { PoolStore } from "../store/pools";

export default ({
  api,
  store,
}: ActionContext<
  "SifService" | "MarketService" | "ClpService",
  "pools" | "wallet"
>) => {
  const state = api.SifService.getState();

  // Sync MarketService with pool store
  api.MarketService.onPoolsUpdated((pools) => {
    for (let pool of pools) {
      store.pools[pool.symbol()] = pool;
    }
  });

  function findPool(pools: PoolStore, a: string, b: string) {
    const key = [a, b].sort().join("_");

    return pools[key] ?? null;
  }

  const actions = {
    async swap(sentAmount: AssetAmount, receivedAsset: Asset) {
      if (!state.address) throw "No from address provided for swap";

      const tx = await api.ClpService.swap({
        fromAddress: state.address,
        sentAmount,
        receivedAsset,
      });

      return await api.SifService.signAndBroadcast(tx.value.msg);
    },

    async addLiquidity(
      nativeAssetAmount: AssetAmount,
      externalAssetAmount: AssetAmount
    ) {
      if (!state.address) throw "No from address provided for swap";
      const hasPool = !!findPool(
        store.pools,
        nativeAssetAmount.asset.symbol,
        externalAssetAmount.asset.symbol
      );

      const provideLiquidity = hasPool
        ? api.ClpService.addLiquidity
        : api.ClpService.createPool;

      const tx = await provideLiquidity({
        fromAddress: state.address,
        nativeAssetAmount,
        externalAssetAmount,
      });

      return await api.SifService.signAndBroadcast(tx.value.msg);
    },

    async removeLiquidity(
      asset: Asset,
      wBasisPoints: string,
      asymmetry: string
    ) {
      const tx = await api.ClpService.removeLiquidity({
        fromAddress: state.address,
        asset,
        asymmetry,
        wBasisPoints,
      });

      return await api.SifService.signAndBroadcast(tx.value.msg);
    },

    async disconnect() {
      api.SifService.purgeClient();
    },
  };

  return actions;
};
