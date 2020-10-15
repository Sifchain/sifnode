// Everything here represents services that are effectively remote data storage
export * from "./utils/getFakeTokens";
export * from "./utils/getSupportedTokens";
export * from "./utils/getWeb3";

import walletService, { WalletServiceContext } from "./walletService";
import tokenService, { TokenServiceContext } from "./tokenService";

type ApiContext = WalletServiceContext & TokenServiceContext; // add contexts from other APIs

export function createApi(context: ApiContext) {
  return {
    walletService: walletService(context),
    tokenService: tokenService(context),
  };
}

export type Api = ReturnType<typeof createApi>;

export type WithApi<T extends keyof Api = keyof Api> = {
  api: Pick<Api, T>;
};
