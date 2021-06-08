import { TxHash, TxParams, Address, Asset, IAssetAmount } from "../entities";
declare type Msg = {
    type: string;
    value: any;
};
export declare type IWalletService = {
    getState: () => {
        address: Address;
        accounts: Address[];
        connected: boolean;
        balances: IAssetAmount[];
        log: string;
    };
    onProviderNotFound(handler: () => void): void;
    onChainIdDetected(handler: (chainId: string) => void): void;
    isConnected(): boolean;
    getSupportedTokens: () => Asset[];
    connect(): Promise<void>;
    disconnect(): Promise<void>;
    transfer(params: TxParams): Promise<TxHash>;
    getBalance(address?: Address, asset?: Asset): Promise<IAssetAmount[]>;
    signAndBroadcast(msg: Msg, memo?: string): Promise<any>;
    setPhrase(phrase: string): Promise<Address>;
    purgeClient(): void;
};
export {};
