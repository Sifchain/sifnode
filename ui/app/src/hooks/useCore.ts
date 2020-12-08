import {
  createStore,
  createApi,
  createActions,
  getWeb3Provider,
  getFakeTokens,
  Asset,
  Pool,
} from "ui-core";
import { Ref, toRefs } from "vue";

const api = createApi({
  // TODO: switch on env
  sifAddrPrefix: "sif",
  sifApiUrl: process.env.VUE_APP_SIFNODE_API || "http://127.0.0.1:1317",
  sifWsUrl:
    process.env.VUE_APP_SIFNODE_WS_API || "ws://localhost:26657/websocket",
  getWeb3Provider,
  loadAssets: getFakeTokens,
});

const store = createStore();
const actions = createActions({ store, api });

type PoolFinderFn = (a: Asset | string, b: Asset | string) => Ref<Pool> | null;
const poolFinder: PoolFinderFn = (a: Asset | string, b: Asset | string) => {
  const pools = toRefs(store.pools);
  const key = [a, b]
    .map((x) => (typeof x === "string" ? x : x.symbol))
    .join("_") as keyof typeof pools;

  const poolRef = pools[key] as Ref<Pool> | undefined;
  return poolRef ?? null;
};

export function useCore() {
  return {
    store,
    api,
    actions,
    poolFinder,
  };
}
