// Everything here represents services that are effectively remote data storage
export * from "./EthereumService/utils/getFakeTokens";
export * from "./EthereumService/utils/getMetamaskProvider";

import ethereumService, { EthereumServiceContext } from "./EthereumService";
import tokenService, { TokenServiceContext } from "./TokenService";
import sifService, { SifServiceContext } from "./SifService";
import clpService, { ClpServiceContext } from "./ClpService";

export type Api = ReturnType<typeof createApi>;

export type WithApi<T extends keyof Api = keyof Api> = {
  api: Pick<Api, T>;
};

export type ApiContext = EthereumServiceContext &
  TokenServiceContext &
  SifServiceContext &
  ClpServiceContext &
  Omit<ClpServiceContext, "getPools">; // add contexts from other APIs

import localnetconfig from "../../config.localnet.json";
import testnetconfig from "../../config.testnet.json";
import { Coin, Network, Token } from "../entities";
import { getMetamaskProvider } from "./EthereumService/utils/getMetamaskProvider";

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

type ChainConfig = {
  sifAddrPrefix: string;
  sifApiUrl: string;
  sifWsUrl: string;
  web3Provider: "metamask" | string;
  assets: AssetConfig[];
  nativeAsset: string; // symbol
};

function parseConfig(config: ChainConfig): ApiContext {
  const nativeAsset = config.assets.find(
    (a) => a.symbol === config.nativeAsset
  );

  if (!nativeAsset)
    throw new Error(
      "No nativeAsset defined for chain config:" + JSON.stringify(config)
    );

  return {
    sifAddrPrefix: config.sifAddrPrefix,
    sifApiUrl: config.sifApiUrl,
    sifWsUrl: config.sifWsUrl,
    getWeb3Provider:
      config.web3Provider === "metamask"
        ? getMetamaskProvider
        : async () => config.web3Provider,
    assets: (config.assets as AssetConfig[]).map(parseAsset),
    nativeAsset: parseAsset(nativeAsset),
  };
}

function getConfig(tag = "localnet"): ApiContext {
  const configMap: ConfigMap = {
    localnet: parseConfig(localnetconfig as ChainConfig),
    testnet: parseConfig(testnetconfig as ChainConfig),
  };

  return configMap[tag.toLowerCase()];
}

export function createApi(tag?: string) {
  const context = getConfig(tag);
  const EthereumService = ethereumService(context);
  const TokenService = tokenService(context);
  const SifService = sifService(context);
  const ClpService = clpService(context);
  return {
    ClpService,
    EthereumService,
    TokenService,
    SifService,
  };
}
