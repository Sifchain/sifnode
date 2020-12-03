// Everything here represents services that are effectively remote data storage
export * from "./EthereumService/utils/getFakeTokens";
export * from "./EthereumService/utils/getWeb3Provider";
export * from "./EthereumService/utils/loadAssets";

import ethereumService, { EthereumServiceContext } from "./EthereumService";
import tokenService, { TokenServiceContext } from "./TokenService";
import sifService, { SifServiceContext } from "./SifService";
import marketService, { MarketServiceContext } from "./MarketService";
import clpService, { ClpServiceContext } from "./ClpService";

export type Api = ReturnType<typeof createApi>;

export type WithApi<T extends keyof Api = keyof Api> = {
  api: Pick<Api, T>;
};

type ApiContext = EthereumServiceContext &
  TokenServiceContext &
  SifServiceContext &
  ClpServiceContext &
  Omit<MarketServiceContext, "getPools">; // add contexts from other APIs

export function createApi(context: ApiContext) {
  const EthereumService = ethereumService(context);
  const TokenService = tokenService(context);
  const SifService = sifService(context);
  const MarketService = marketService(context);
  const ClpService = clpService(context);
  return {
    MarketService,
    EthereumService,
    TokenService,
    SifService,
    ClpService,
  };
}
