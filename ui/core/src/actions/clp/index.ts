import { Asset, AssetAmount, LiquidityProvider, Pool } from "../../entities";
import { ActionContext } from "..";
import { PoolStore } from "../../store/pools";
import notify from "../../api/utils/Notifications";
import { toPool } from "../../api/utils/SifClient/toPool";
import { effect } from "@vue/reactivity";

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
          detail: "Create liquidity pool to swap.",
        });
      }
    });
  });

  // Then every transaction

  api.SifService.onNewBlock(async () => {
    await syncPools();
  });

  api.SifService.onSocketError(({ sifWsUrl }) => {
    notify({
      type: "error",
      message: "Websocket Not Connected",
      detail: `${sifWsUrl}`,
    });
  });

  function findPool(pools: PoolStore, a: string, b: string) {
    const key = [a, b].sort().join("_");

    return pools[key] ?? null;
  }

  const actions = {
    async swap(sentAmount: AssetAmount, receivedAsset: Asset, minimumReceived: string) {
      if (!state.address) throw "No from address provided for swap";

      const tx = await api.ClpService.swap({
        fromAddress: state.address,
        sentAmount,
        receivedAsset,
        minimumReceived
      });

      return await api.SifService.signAndBroadcast(tx.value.msg);
    },

    async addLiquidity(
      nativeAssetAmount: AssetAmount,
      externalAssetAmount: AssetAmount
    ) {
      const response = {
        hash: <string>"",
        state: <string>"",
        stateMsg: <string>""
      }
      try {
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

        const txHash = await api.SifService.signAndBroadcast(tx.value.msg);
          
        if (txHash && txHash.rawLog && txHash.rawLog.includes('"type":"added_liquidity"')) {
          response.state = "confirmed";
          response.hash = txHash?.transactionHash ?? "";
        } else {
          response.state = "failed";
          response.stateMsg = "Oops... Something went wrong. Please try again!";
        }

        return response;

      } catch (err) {
        // TODO: coordinate with blockchain to get more standardised errors
        if (err.message) {
          notify({
            type: "error",
            message: err.message,
          });
        }

        // TODO: check the type of error, and notify frontend accordingly
        response.state = "failed";
        response.stateMsg = "Oops... Something went wrong. Please try again!";

        return response;
      }
    },

    async removeLiquidity(
      asset: Asset,
      wBasisPoints: string,
      asymmetry: string
    ) {
      try {
        const tx = await api.ClpService.removeLiquidity({
          fromAddress: state.address,
          asset,
          asymmetry,
          wBasisPoints,
        });

        return await api.SifService.signAndBroadcast(tx.value.msg);
      } catch (err) {
        // TODO: coordinate with blockchain to get more standardised errors
        if (err.message) {
          notify({
            type: "error",
            message: err.message,
          });
        }
      }
    },

    async disconnect() {
      api.SifService.purgeClient();
    },
  };

  return actions;
};
