import { Services, WithService } from "../services";
import { Store, WithStore } from "../store";
export declare type UsecaseContext<T extends keyof Services = keyof Services, S extends keyof Store = keyof Store> = WithService<T> & WithStore<S>;
export declare function createUsecases(context: UsecaseContext): {
    clp: {
        swap(sentAmount: import("..").IAssetAmount, receivedAsset: import("..").IAsset, minimumReceived: import("..").IAssetAmount): Promise<import("..").TransactionStatus>;
        addLiquidity(nativeAssetAmount: import("..").IAssetAmount, externalAssetAmount: import("..").IAssetAmount): Promise<import("..").TransactionStatus>;
        removeLiquidity(asset: import("..").IAsset, wBasisPoints: string, asymmetry: string): Promise<import("..").TransactionStatus>;
        disconnect(): Promise<void>;
    };
    wallet: {
        sif: {
            getCosmosBalances(address: string): Promise<import("..").IAssetAmount[]>;
            connect(mnemonic: string): Promise<string>;
            sendCosmosTransaction(params: import("..").TxParams): Promise<any>;
            disconnect(): Promise<void>;
            connectToWallet(): Promise<void>;
        };
        eth: {
            isSupportedNetwork(): boolean;
            disconnectWallet(): Promise<void>;
            connectToWallet(): Promise<void>;
            transferEthWallet(amount: number, recipient: string, asset: import("..").IAsset): Promise<string>;
        };
    };
    peg: {
        subscribeToUnconfirmedPegTxs: () => () => void;
        getSifTokens(): import("..").IAsset[];
        getEthTokens(): import("..").IAsset[];
        calculateUnpegFee(asset: import("..").IAsset): import("..").IAssetAmount;
        unpeg(assetAmount: import("..").IAssetAmount): Promise<import("..").TransactionStatus>;
        approve(address: string, assetAmount: import("..").IAssetAmount): Promise<any>;
        peg(assetAmount: import("..").IAssetAmount): Promise<import("..").TransactionStatus>;
    };
};
