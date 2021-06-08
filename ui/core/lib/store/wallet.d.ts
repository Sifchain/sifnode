import { Address, IAssetAmount } from "../entities";
export declare type WalletStore = {
    eth: {
        chainId?: string;
        balances: IAssetAmount[];
        isConnected: boolean;
        address: Address;
    };
    sif: {
        balances: IAssetAmount[];
        isConnected: boolean;
        address: Address;
    };
};
export declare const wallet: WalletStore;
