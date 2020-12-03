import { AssetAmount } from "./AssetAmount";

export type Wallet = {
  addresses: WalletAddress[];
};

export type WalletAddress = {
  address: string;
  balance: AssetAmount;
};

export type Mnemonic = string;
