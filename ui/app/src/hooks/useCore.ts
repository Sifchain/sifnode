import {
  createStore,
  createApi,
  createActions,
  ApiContext,
  getMetamaskProvider,
  Asset,
  Pool,
  Token,
  Coin,
  Network,
} from "ui-core";
import { Ref, toRefs } from "vue";

import localnetconfig from "../../config.localnet.json";
import testnetconfig from "../../config.testnet.json";

type ConfigMap = { [s: string]: ApiContext };

type TokenConfig = {
  symbol: string;
  decimals: number;
  imageUrl?: string;
  name: string;
  address: string;
  network: Network;
};

type CoinConfig = {
  symbol: string;
  decimals: number;
  imageUrl?: string;
  name: string;
  network: Network;
};

type AssetConfig = CoinConfig | TokenConfig;
function isTokenConfig(a: any): a is TokenConfig {
  return typeof a?.address === "string";
}

function parseAsset(a: unknown) {
  if (isTokenConfig(a)) {
    return Token(a);
  }
  return Coin(a as CoinConfig);
}

function getConfig(tag = "localnet"): ApiContext {
  const configMap: ConfigMap = {
    localnet: {
      sifAddrPrefix: "sif",
      sifApiUrl: "http://127.0.0.1:1317",
      sifWsUrl: "ws://localhost:26657/websocket",
      getWeb3Provider: getMetamaskProvider,
      supportedAssets: (localnetconfig.assets as AssetConfig[]).map(parseAsset),
      nativeAsset: Coin({
        symbol: "rowan",
        decimals: 18,
        name: "Rowan",
        network: Network.SIFCHAIN,
      }),
    },
    testnet: {
      sifAddrPrefix: "sif",
      sifApiUrl: "http://127.0.0.1:1317",
      sifWsUrl: "ws://localhost:26657/websocket",
      getWeb3Provider: getMetamaskProvider,
      supportedAssets: (testnetconfig.assets as AssetConfig[]).map(parseAsset),
      nativeAsset: Coin({
        symbol: "rowan",
        decimals: 18,
        name: "Rowan",
        network: Network.SIFCHAIN,
      }),
    },
    alphanet: {
      sifAddrPrefix: "sif",
      sifApiUrl: "http://127.0.0.1:1317",
      sifWsUrl: "ws://localhost:26657/websocket",
      getWeb3Provider: getMetamaskProvider,
      supportedAssets: (testnetconfig.assets as AssetConfig[]).map(parseAsset),
      nativeAsset: Coin({
        symbol: "rowan",
        decimals: 18,
        name: "Rowan",
        network: Network.SIFCHAIN,
      }),
    },
    betanet: {
      sifAddrPrefix: "sif",
      sifApiUrl: "http://127.0.0.1:1317",
      sifWsUrl: "ws://localhost:26657/websocket",
      getWeb3Provider: getMetamaskProvider,
      supportedAssets: testnetconfig.assets.map(parseAsset),
      nativeAsset: Coin({
        symbol: "rowan",
        decimals: 18,
        name: "Rowan",
        network: Network.SIFCHAIN,
      }),
    },
    mainnet: {
      sifAddrPrefix: "sif",
      sifApiUrl: "http://127.0.0.1:1317",
      sifWsUrl: "ws://localhost:26657/websocket",
      getWeb3Provider: getMetamaskProvider,
      supportedAssets: (testnetconfig.assets as AssetConfig[]).map(parseAsset),
      nativeAsset: Coin({
        symbol: "rowan",
        decimals: 18,
        name: "Rowan",
        network: Network.SIFCHAIN,
      }),
    },
  };

  return configMap[tag.toLowerCase()];
}

const api = createApi(getConfig(process.env.VUE_APP_DEPLOYMENT_TAG));

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
