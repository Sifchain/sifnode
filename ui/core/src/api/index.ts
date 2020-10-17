// Everything here represents services that are effectively remote data storage
export * from "./EtheriumService/utils/getFakeTokens";
export * from "./EtheriumService/utils/getWeb3Provider";

import JSBI from "jsbi";
import { Address, Asset, Balance, Balances, Token } from "../entities";
import etheriumService, { EtheriumServiceContext } from "./EtheriumService";
import tokenService, { TokenServiceContext } from "./TokenService";

export type Api = ReturnType<typeof createApi>;

export type WithApi<T extends keyof Api = keyof Api> = {
  api: Pick<Api, T>;
};

export type TxParams = {
  asset?: Asset;
  amount: JSBI;
  recipient: Address;
  feeRate?: number; // optional feeRate
  memo?: string; // optional memo to pass
};

export type TxHash = string;

export type IWalletService = {
  onDisconnected(handler: (...a: any[]) => void): void;
  onConnected(handler: (...a: any[]) => void): void;
  onChange(handler: (...a: any[]) => void): void;
  getAddress(): Promise<Address | null>;
  isConnected(): boolean;
  connect(): Promise<void>;
  disconnect(): Promise<void>;
  transfer(params: TxParams): Promise<TxHash>;
  getBalance(address?: Address, asset?: Asset | Token): Promise<Balances>;
};

type ApiContext = EtheriumServiceContext & TokenServiceContext; // add contexts from other APIs

export function createApi(context: ApiContext) {
  return {
    EtheriumService: etheriumService(context),
    TokenService: tokenService(context),
  };
}
