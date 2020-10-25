import {
  createStore,
  createApi,
  createActions,
  getWeb3Provider,
  getFakeTokens,
  // loadAssets,
} from "../../../core";

// import tokens from "../../../core/data/topErc20Tokens.json";

const api = createApi({
  // TODO: switch on env
  sifAddrPrefix: "sif",
  sifApiUrl: "http://127.0.0.1:1317",
  getWeb3Provider,
  loadAssets: getFakeTokens,
  fetchMarketData: async () => {
    // Setup a new fake pools
    return [
      [
        { name: "atk", value: 200 },
        { name: "btk", value: 100 },
      ],
      [
        { name: "atk", value: 100 },
        { name: "eth", value: 5 },
      ],
      [
        { name: "btk", value: 150 },
        { name: "eth", value: 5 },
      ],
    ];
  },
});

const store = createStore();
const actions = createActions({ store, api });
export function useCore() {
  return {
    store,
    api,
    actions,
  };
}
