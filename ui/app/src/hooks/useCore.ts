import {
  createStore,
  createApi,
  createActions,
  getWeb3Provider,
  getFakeTokens as getSupportedTokens,
} from "../../../core";

const api = createApi({
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
