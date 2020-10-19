// Everything here represents services that are effectively remote data storage
export * from "./EthereumService/utils/getFakeTokens";
export * from "./EthereumService/utils/getWeb3Provider";

import JSBI from "jsbi";
import { Address, Asset, Balances, Token } from "../entities";
import EthereumService, { EthereumServiceContext } from "./EthereumService";
import tokenService, { TokenServiceContext } from "./TokenService";
import sifService, { SifServiceContext } from "./SifService";

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
  getState: () => {
    address: Address;
    accounts: Address[];
    connected: boolean;
    log: string;
  };
  isConnected(): boolean;
  connect(): Promise<void>;
  disconnect(): Promise<void>;
  transfer(params: TxParams): Promise<TxHash>;
  getBalance(address?: Address, asset?: Asset | Token): Promise<Balances>;
  // FOLLOWING ARE YTI:

  // setPhrase(phrase: string): Address
  // setNetwork(net: Network): void
  // getNetwork(): Network

  // getExplorerAddressUrl(address: Address): string
  // getExplorerTxUrl(txID: string): string
  // getTransactions(params?: TxHistoryParams): Promise<TxsPage>

  // getFees(): Promise<Fees>

  // transfer(params: TxParams): Promise<TxHash>
  // deposit(params: TxParams): Promise<TxHash>

  // purgeClient(): void
};

type ApiContext = EthereumServiceContext &
  TokenServiceContext &
  SifServiceContext; // add contexts from other APIs

export function createApi(context: ApiContext) {
  return {
    EthereumService: EthereumService(context),
    TokenService: tokenService(context),
    SifService: sifService(context),
  };
}
