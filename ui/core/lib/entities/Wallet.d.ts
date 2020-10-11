import { AssetAmount } from "./AssetAmount";
export declare type Wallet = {
    addresses: WalletAddress[];
};
export declare type WalletAddress = {
    address: string;
    assetAmount: AssetAmount;
};
