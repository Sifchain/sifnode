import {
  createStore,
  createServices,
  createUsecases,
  createPoolFinder,
  getConfig,
} from "ui-core";

const config = getConfig(
  process.env.VUE_APP_DEPLOYMENT_TAG,
  process.env.VUE_APP_SIFCHAIN_ASSET_TAG,
  process.env.VUE_APP_ETHEREUM_ASSET_TAG,
);

const api = createServices(config);
const store = createStore();
const actions = createUsecases({ store, api });
const poolFinder = createPoolFinder(store);

// expose store on window so it is easy to inspect
Object.defineProperty(window, "store", {
  get: function () {
    // Gives us `store` for in console inspection
    // Gives us `store.dump()` for string representation
    // Gives us `store.dumpTab()` for string representation
    const storeString = JSON.stringify(
      store,
      (_, value) => {
        // TODO give all entities a toString so we don't have to do this
        // if AssetAmount
        if (value?.asset && value?.quotient) {
          return value.toString();
        }

        // If Fraction
        if (value?.numerator && value?.denominator) {
          return value.toFixed(18);
        }

        return value;
      },
      2,
    );

    const storeSafe = JSON.parse(storeString);
    storeSafe.dumpTab = () => {
      const x = window.open();
      x?.document.open();
      x?.document.write("<pre>" + storeString + "</pre>");
      x?.document.close();
    };
    storeSafe.dump = () => {
      return storeString;
    };

    return storeSafe;
  },
});

export function useCore() {
  return {
    store,
    api,
    actions,
    poolFinder,
    config,
  };
}
