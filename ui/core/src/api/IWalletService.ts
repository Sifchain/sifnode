import { TxHash, TxParams, Address, Asset, AssetAmount } from "../entities";

type Msg = { type: string; value: any }; // make entity

export type IWalletService = {
  getState: () => {
    address: Address;
    accounts: Address[];
    connected: boolean;
    balances: AssetAmount[];
    log: string;
  };
  onProviderNotFound(handler: () => void): void;
  isConnected(): boolean;
  getSupportedTokens: () => Asset[];
  connect(): Promise<void>;
  disconnect(): Promise<void>;
  transfer(params: TxParams): Promise<TxHash>;
  getBalance(address?: Address, asset?: Asset): Promise<AssetAmount[]>;
  signAndBroadcast(msg: Msg, memo?: string): Promise<any>;
  setPhrase(phrase: string): Promise<Address>;
  purgeClient(): void;

  // FOLLOWING ARE YTI:

  // setNetwork(net: Network): void
  // getNetwork(): Network

  // getExplorerAddressUrl(address: Address): string
  // getExplorerTxUrl(txID: string): string
  // getTransactions(params?: TxHistoryParams): Promise<TxsPage>

  // getFees(): Promise<Fees>

  // transfer(params: TxParams): Promise<TxHash>
  // deposit(params: TxParams): Promise<TxHash>
};
