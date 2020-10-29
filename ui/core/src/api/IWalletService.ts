import { TxHash, TxParams, Address, Asset, Balance, Token } from "../entities";

export type IWalletService = {
  getState: () => {
    address: Address;
    accounts: Address[];
    connected: boolean;
    balances: Balance[];
    log: string;
  };
  isConnected(): boolean;
  connect(): Promise<void>;
  disconnect(): Promise<void>;
  transfer(params: TxParams): Promise<TxHash>;
  getBalance(address?: Address, asset?: Asset | Token): Promise<Balance[]>;

  // FOLLOWING ARE YTI:

  setPhrase(phrase: string): Promise<Address>;
  // setNetwork(net: Network): void
  // getNetwork(): Network

  // getExplorerAddressUrl(address: Address): string
  // getExplorerTxUrl(txID: string): string
  // getTransactions(params?: TxHistoryParams): Promise<TxsPage>

  // getFees(): Promise<Fees>

  // transfer(params: TxParams): Promise<TxHash>
  // deposit(params: TxParams): Promise<TxHash>

  purgeClient(): void;
};
