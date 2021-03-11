import {
  Asset,
  AssetAmount,
  Fraction,
  LiquidityProvider,
  Pool,
} from "../../entities";
import { ActionContext } from "..";
import { PoolStore } from "../../store/pools";
import { effect } from "@vue/reactivity";
import JSBI from "jsbi";

export default ({
  api,
  store,
}: ActionContext<
  "SifService" | "ClpService" | "EventsService",
  "pools" | "wallet" | "accountpools"
>) => {
  const state = api.SifService.getState();

  async function syncPools() {
    const state = api.SifService.getState();

    // UPdate pools
    const pools = await api.ClpService.getPools();
    for (let pool of pools) {
      store.pools[pool.symbol()] = pool;
    }

    // Update lp pools
    if (state.address) {
      const accountPoolSymbols = await api.ClpService.getPoolSymbolsByLiquidityProvider(
        state.address
      );

      // This is a hot method when there are a heap of pools
      // Ideally we would have a better rest endpoint design

      accountPoolSymbols.forEach(async symbol => {
        const lp = await api.ClpService.getLiquidityProvider({
          symbol,
          lpAddress: state.address,
        });
        if (!lp) return;
        const pool = `${symbol}_rowan`;
        store.accountpools[state.address] =
          store.accountpools[state.address] || {};

        store.accountpools[state.address][pool] = { lp, pool };
      });

      // Delete accountpools
      const currentPoolIds = accountPoolSymbols.map(id => `${id}_rowan`);
      if (store.accountpools[state.address]) {
        const existingPoolIds = Object.keys(store.accountpools[state.address]);
        const disjunctiveIds = existingPoolIds.filter(
          id => !currentPoolIds.includes(id)
        );

        disjunctiveIds.forEach(poolToRemove => {
          delete store.accountpools[state.address][poolToRemove];
        });
      }
    }
  }

  // Sync on load
  syncPools().then(() => {
    effect(() => {
      if (Object.keys(store.pools).length === 0) {
        api.EventsService.notify({
          type: "NoLiquidityPoolsFoundEvent",
          payload: {},
        });
      }
    });
  });

  // Then every transaction

  api.SifService.onNewBlock(async () => {
    await syncPools();
  });

  function findPool(pools: PoolStore, a: string, b: string) {
    const key = [a, b].sort().join("_");

    return pools[key] ?? null;
  }

  effect(() => {
    // When sif address changes syncPools
    store.wallet.sif.address;
    syncPools();
  });

  const actions = {
    async swap(
      sentAmount: AssetAmount,
      receivedAsset: Asset,
      minimumReceived: AssetAmount
    ) {
      if (!state.address) throw "No from address provided for swap";

      const tx = await api.ClpService.swap({
        fromAddress: state.address,
        sentAmount,
        receivedAsset,
        minimumReceived,
      });

      const txStatus = await api.SifService.signAndBroadcast(tx.value.msg);

      if (txStatus.state !== "accepted") {
        api.EventsService.notify({
          type: "TransactionErrorEvent",
          payload: {
            txStatus,
            message: txStatus.memo || "There was an error with your swap",
          },
        });
      }

      return txStatus;
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

      const txStatus = await api.SifService.signAndBroadcast(tx.value.msg);
      if (txStatus.state !== "accepted") {
        api.EventsService.notify({
          type: "TransactionErrorEvent",
          payload: {
            txStatus,
            message: txStatus.memo || "There was an error with your swap",
          },
        });
      }
      return txStatus;
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

      const txStatus = await api.SifService.signAndBroadcast(tx.value.msg);

      if (txStatus.state !== "accepted") {
        api.EventsService.notify({
          type: "TransactionErrorEvent",
          payload: {
            txStatus,
            message: txStatus.memo || "There was an error removing liquidity",
          },
        });
      }

      return txStatus;
    },

    async disconnect() {
      api.SifService.purgeClient();
    },
  };

  return actions;
};
