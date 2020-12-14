// Everything here represents services that are effectively remote data storage
export * from "./EthereumService/utils/getMetamaskProvider";

import ethereumService, { EthereumServiceContext } from "./EthereumService";
import sifService, { SifServiceContext } from "./SifService";
import clpService, { ClpServiceContext } from "./ClpService";

export type Api = ReturnType<typeof createApi>;

export type WithApi<T extends keyof Api = keyof Api> = {
  api: Pick<Api, T>;
};

export type ApiContext = EthereumServiceContext &
  SifServiceContext &
  ClpServiceContext &
  Omit<ClpServiceContext, "getPools">; // add contexts from other APIs

import localnetconfig from "../config.localnet.json";
import testnetconfig from "../config.testnet.json";
import { parseConfig, ChainConfig } from "./utils/parseConfig";

type ConfigMap = { [s: string]: ApiContext };

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

  const SifService = sifService(context);
  const ClpService = clpService(context);
  return {
    ClpService,
    EthereumService,

    SifService,
  };
}
