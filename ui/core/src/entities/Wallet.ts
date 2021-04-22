import { IAssetAmount } from "./AssetAmount";

export type Wallet = {
  addresses: WalletAddress[];
};

export type WalletAddress = {
  address: string;
  balance: IAssetAmount;
};

export type Mnemonic = string;
