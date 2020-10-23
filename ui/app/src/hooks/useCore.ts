import {
  createStore,
  createApi,
  createActions,
  getWeb3Provider,
  getFakeTokens as getSupportedTokens,
} from "../../../core";

const api = createApi({
  // TODO: switch on env
  sifAddrPrefix: "sif",
  sifApiUrl: "http://127.0.0.1:1317",
  getWeb3Provider,
  getSupportedTokens,
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
