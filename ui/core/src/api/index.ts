// Everything here represents services that are effectively remote data storage
export * from "./EthereumService/utils/getFakeTokens";
export * from "./EthereumService/utils/getWeb3Provider";
export * from "./EthereumService/utils/loadAssets";

import EthereumService, { EthereumServiceContext } from "./EthereumService";
import tokenService, { TokenServiceContext } from "./TokenService";
import sifService, { SifServiceContext } from "./SifService";
import marketService, { MarketServiceContext } from "./MarketService";

export type Api = ReturnType<typeof createApi>;

export type WithApi<T extends keyof Api = keyof Api> = {
  api: Pick<Api, T>;
};

type ApiContext = EthereumServiceContext &
  TokenServiceContext &
  SifServiceContext &
  MarketServiceContext; // add contexts from other APIs

export function createApi(context: ApiContext) {
  return {
    EthereumService: EthereumService(context),
    TokenService: tokenService(context),
    SifService: sifService(context),
    MarketService: marketService(context),
  };
}
