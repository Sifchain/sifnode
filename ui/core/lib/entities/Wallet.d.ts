import { IAssetAmount } from "./AssetAmount";
export declare type Wallet = {
    addresses: WalletAddress[];
};
export declare type WalletAddress = {
    address: string;
    balance: IAssetAmount;
};
export declare type Mnemonic = string;
