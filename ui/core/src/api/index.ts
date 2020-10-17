// Everything here represents services that are effectively remote data storage
export * from "./EtheriumService/utils/getFakeTokens";
export * from "./EtheriumService/utils/getWeb3Provider";

import { Address, Asset, Balances, Token } from "../entities";
import etheriumService, { EtheriumServiceContext } from "./EtheriumService";
import tokenService, { TokenServiceContext } from "./TokenService";

export type Api = ReturnType<typeof createApi>;

export type WithApi<T extends keyof Api = keyof Api> = {
  api: Pick<Api, T>;
};

export type IWalletService = {
  onDisconnected(handler: (...a: any[]) => void): void;
  onConnected(handler: (...a: any[]) => void): void;
  getAddress(): Promise<Address | null>;
  isConnected(): boolean;
  connect(): Promise<void>;
  disconnect(): Promise<void>;
  getBalance(address?: Address, asset?: Asset | Token): Promise<Balances>;
};

type ApiContext = EtheriumServiceContext & TokenServiceContext; // add contexts from other APIs

export function createApi(context: ApiContext) {
  return {
    EtheriumService: etheriumService(context),
    TokenService: tokenService(context),
  };
}
