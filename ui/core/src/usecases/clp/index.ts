import { IAsset, IAssetAmount } from "../../entities";
import { UsecaseContext } from "..";
import { PoolStore } from "../../store/pools";
import { effect } from "@vue/reactivity";

export default ({
  services,
  store,
}: UsecaseContext<
  "SifService" | "ClpService" | "EventBusService",
  "pools" | "wallet" | "accountpools"
>) => {
  const state = services.SifService.getState();

  async function syncPools() {
    const state = services.SifService.getState();

    // UPdate pools
    const pools = await services.ClpService.getPools();
    for (let pool of pools) {
      store.pools[pool.symbol()] = pool;
    }

    // Update lp pools
    if (state.address) {
      const accountPoolSymbols = await services.ClpService.getPoolSymbolsByLiquidityProvider(
        state.address,
      );

      // This is a hot method when there are a heap of pools
      // Ideally we would have a better rest endpoint design

      accountPoolSymbols.forEach(async (symbol) => {
        const lp = await services.ClpService.getLiquidityProvider({
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
      const currentPoolIds = accountPoolSymbols.map((id) => `${id}_rowan`);
      if (store.accountpools[state.address]) {
        const existingPoolIds = Object.keys(store.accountpools[state.address]);
        const disjunctiveIds = existingPoolIds.filter(
          (id) => !currentPoolIds.includes(id),
        );

        disjunctiveIds.forEach((poolToRemove) => {
          delete store.accountpools[state.address][poolToRemove];
        });
      }
    }
  }

  // Sync on load
  syncPools().then(() => {
    effect(() => {
      if (Object.keys(store.pools).length === 0) {
        services.EventBusService.dispatch({
          type: "NoLiquidityPoolsFoundEvent",
          payload: {},
        });
      }
    });
  });

  // Then every transaction

  services.SifService.onNewBlock(async () => {
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
      sentAmount: IAssetAmount,
      receivedAsset: IAsset,
      minimumReceived: IAssetAmount,
    ) {
      if (!state.address) throw "No from address provided for swap";

      const tx = await services.ClpService.swap({
        fromAddress: state.address,
        sentAmount,
        receivedAsset,
        minimumReceived,
      });

      const txStatus = await services.SifService.signAndBroadcast(tx.value.msg);

      if (txStatus.state !== "accepted") {
        services.EventBusService.dispatch({
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
      nativeAssetAmount: IAssetAmount,
      externalAssetAmount: IAssetAmount,
    ) {
      if (!state.address) throw "No from address provided for swap";
      const hasPool = !!findPool(
        store.pools,
        nativeAssetAmount.asset.symbol,
        externalAssetAmount.asset.symbol,
      );

      const provideLiquidity = hasPool
        ? services.ClpService.addLiquidity
        : services.ClpService.createPool;

      const tx = await provideLiquidity({
        fromAddress: state.address,
        nativeAssetAmount,
        externalAssetAmount,
      });

      const txStatus = await services.SifService.signAndBroadcast(tx.value.msg);
      if (txStatus.state !== "accepted") {
        services.EventBusService.dispatch({
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
      asset: IAsset,
      wBasisPoints: string,
      asymmetry: string,
    ) {
      const tx = await services.ClpService.removeLiquidity({
        fromAddress: state.address,
        asset,
        asymmetry,
        wBasisPoints,
      });

      const txStatus = await services.SifService.signAndBroadcast(tx.value.msg);

      if (txStatus.state !== "accepted") {
        services.EventBusService.dispatch({
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
      services.SifService.purgeClient();
    },
  };

  return actions;
};
