import {
  createStore,
  createApi,
  getWeb3,
  createActions,
  getFakeTokens as getSupportedTokens,
} from "../../../core";

const api = createApi({
  getWeb3,
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
