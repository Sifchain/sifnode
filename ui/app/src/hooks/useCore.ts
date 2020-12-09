import {
  createStore,
  createApi,
  createActions,
  getWeb3Provider,
  getFakeTokens,
  Coin,
  Network,
  Asset,
} from "ui-core";
import { toRefs } from "vue";

const api = createApi({
  // TODO: switch on env
  sifAddrPrefix: "sif",
  sifApiUrl: process.env.VUE_APP_SIFNODE_API || "http://127.0.0.1:1317",
  getWeb3Provider,
  nativeAsset: Coin({
    name: "Rowan",
    symbol: "rwn",
    decimals: 18,
    network: Network.SIFCHAIN,
  }),
  loadAssets: getFakeTokens,
});

const store = createStore();
const actions = createActions({ store, api });

const poolFinder = (a: Asset | string, b: Asset | string) => {
  const pools = toRefs(store.pools);
  const key = [a, b]
    .map((x) => (typeof x === "string" ? x : x.symbol))
    .join("_") as keyof typeof pools;
  return pools[key] || null;
};

export function useCore() {
  return {
    store,
    api,
    actions,
    poolFinder,
  };
}
