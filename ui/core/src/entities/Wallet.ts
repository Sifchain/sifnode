import { Balance } from "./Balance";

export type Wallet = {
  addresses: WalletAddress[];
};

export type WalletAddress = {
  address: string;
  balance: Balance;
};
