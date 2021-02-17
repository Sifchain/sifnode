import {Asset, AssetAmount, Fraction, LiquidityProvider, Pool} from "../../entities";
import { ActionContext } from "..";
import { PoolStore } from "../../store/pools";
import notify from "../../api/utils/Notifications";
import { effect } from "@vue/reactivity";
import JSBI from "jsbi";

export default ({
  api,
  store,
}: ActionContext<
  "SifService" | "ClpService",
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

      const accountPools: { lp: LiquidityProvider; pool: Pool }[] = [];
      for (const symbol of accountPoolSymbols) {
        const lp = await api.ClpService.getLiquidityProvider({
          symbol,
          lpAddress: state.address,
        });
        if (!lp) continue;
        const pool = store.pools[`${symbol}_rowan`];
        accountPools.push({ lp, pool });
      }
      store.accountpools = accountPools;
    }
  }

  // Sync on load
  syncPools().then(() => {
    effect(() => {
      if (Object.keys(store.pools).length === 0) {
        notify({
          type: "error",
          message: "No Liquidity Pools Found",
          detail: {
            type: "info",
            message: "Create liquidity pool to swap.",
          }
        });
      }
    });
  });

  // Then every transaction

  api.SifService.onNewBlock(async () => {
    await syncPools();
  });

  api.SifService.onSocketError((instance) => {
    notify({
      type: "error",
      message: "Websocket Not Connected",
      detail: {
        type: 'websocket',
        message: instance.target.url
      }
    });
  });

  function findPool(pools: PoolStore, a: string, b: string) {
    const key = [a, b].sort().join("_");

    return pools[key] ?? null;
  }

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
        notify({
          type: "error",
          message: txStatus.memo || "There was an error with your swap",
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
        notify({
          type: "error",
          message: txStatus.memo || "There was an error adding liquidity",
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
        notify({
          type: "error",
          message: txStatus.memo || "There was an error removing liquidity",
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
