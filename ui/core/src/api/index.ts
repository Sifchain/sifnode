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

export type FullApi = ReturnType<typeof createApi>;

export type Api<
  T extends keyof FullApi = keyof FullApi,
  U extends object = {}
> = {
  api: Pick<FullApi, T>;
} & U;
