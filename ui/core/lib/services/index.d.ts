export * from "./EthereumService/utils/getMetamaskProvider";
import { EthereumServiceContext } from "./EthereumService";
import { EthbridgeServiceContext } from "./EthbridgeService";
import { SifServiceContext } from "./SifService";
import { ClpServiceContext } from "./ClpService";
import { EventBusServiceContext } from "./EventBusService";
export declare type Services = ReturnType<typeof createServices>;
export declare type WithService<T extends keyof Services = keyof Services> = {
    services: Pick<Services, T>;
};
export declare type ServiceContext = EthereumServiceContext & SifServiceContext & ClpServiceContext & EthbridgeServiceContext & ClpServiceContext & EventBusServiceContext;
export declare function createServices(context: ServiceContext): {
    clp: {
        getPools: () => Promise<{
            amounts: [import("..").IAssetAmount, import("..").IAssetAmount];
            otherAsset: (asset: import("..").IAsset) => import("..").IAssetAmount;
            symbol: () => string;
            contains: (...assets: import("..").IAsset[]) => boolean;
            toString: () => string;
            getAmount: (asset: string | import("..").IAsset) => import("..").IAssetAmount;
            poolUnits: import("..").IAmount;
            priceAsset(asset: import("..").IAsset): import("..").IAssetAmount;
            calcProviderFee(x: import("..").IAssetAmount): import("..").IAssetAmount;
            calcPriceImpact(x: import("..").IAssetAmount): import("..").IAmount;
            calcSwapResult(x: import("..").IAssetAmount): import("..").IAssetAmount;
            calcReverseSwapResult(Sa: import("..").IAssetAmount): import("..").IAssetAmount;
            calculatePoolUnits(nativeAssetAmount: import("..").IAssetAmount, externalAssetAmount: import("..").IAssetAmount): import("..").IAmount[];
        }[]>;
        getPoolSymbolsByLiquidityProvider: (address: string) => Promise<string[]>;
        swap: (params: {
            fromAddress: string;
            sentAmount: import("..").IAssetAmount;
            receivedAsset: import("..").IAsset;
            minimumReceived: import("..").IAssetAmount;
        }) => any;
        addLiquidity: (params: {
            fromAddress: string;
            nativeAssetAmount: import("..").IAssetAmount;
            externalAssetAmount: import("..").IAssetAmount;
        }) => any;
        createPool: (params: {
            fromAddress: string;
            nativeAssetAmount: import("..").IAssetAmount;
            externalAssetAmount: import("..").IAssetAmount;
        }) => any;
        getLiquidityProvider: (params: {
            symbol: string;
            lpAddress: string;
        }) => Promise<import("..").LiquidityProvider | null>;
        removeLiquidity: (params: {
            wBasisPoints: string;
            asymmetry: string;
            asset: import("..").IAsset;
            fromAddress: string;
        }) => any;
    };
    eth: import("./IWalletService").IWalletService;
    sif: {
        getState(): {
            connected: boolean;
            address: string;
            accounts: string[];
            balances: import("..").IAssetAmount[];
            log: string;
        };
        getSupportedTokens(): import("..").IAsset[];
        setClient(): Promise<void>;
        initProvider(): Promise<void>;
        connect(): Promise<void>;
        isConnected(): boolean;
        onSocketError(handler: (a: any) => void): void;
        onTx(handler: (a: any) => void): void;
        onNewBlock(handler: (a: any) => void): void;
        setPhrase(mnemonic: string): Promise<string>;
        purgeClient(): Promise<void>;
        getBalance(address?: string | undefined, asset?: import("..").IAsset | undefined): Promise<import("..").IAssetAmount[]>;
        transfer(params: import("..").TxParams): Promise<any>;
        signAndBroadcast(msg: import("@cosmjs/launchpad").Msg | import("@cosmjs/launchpad").Msg[], memo?: string | undefined): Promise<import("..").TransactionStatus>;
    };
    ethbridge: {
        approveBridgeBankSpend(account: string, amount: import("..").IAssetAmount): Promise<any>;
        burnToEthereum(params: {
            fromAddress: string;
            ethereumRecipient: string;
            assetAmount: import("..").IAssetAmount;
            feeAmount: import("..").IAssetAmount;
        }): Promise<import("@cosmjs/launchpad").Msg>;
        lockToSifchain(sifRecipient: string, assetAmount: import("..").IAssetAmount, confirmations: number): {
            readonly hash: string | undefined;
            readonly symbol: string | undefined;
            setTxHash(hash: string): void;
            emit(e: import("./EthbridgeService/types").TxEventPrepopulated<import("./EthbridgeService/types").TxEvent>): void;
            onTxEvent(handler: (e: import("./EthbridgeService/types").TxEvent) => void): any;
            onTxHash(handler: (e: import("./EthbridgeService/types").TxEventHashReceived) => void): any;
            onEthConfCountChanged(handler: (e: import("./EthbridgeService/types").TxEventEthConfCountChanged) => void): any;
            onEthTxConfirmed(handler: (e: import("./EthbridgeService/types").TxEventEthTxConfirmed) => void): any;
            onSifTxConfirmed(handler: (e: import("./EthbridgeService/types").TxEventSifTxConfirmed) => void): any;
            onEthTxInitiated(handler: (e: import("./EthbridgeService/types").TxEventEthTxInitiated) => void): any;
            onSifTxInitiated(handler: (e: import("./EthbridgeService/types").TxEventSifTxInitiated) => void): any;
            onComplete(handler: (e: import("./EthbridgeService/types").TxEventComplete) => void): any;
            removeListeners(): any;
            onError(handler: (e: import("./EthbridgeService/types").TxEventError) => void): any;
        };
        lockToEthereum(params: {
            fromAddress: string;
            ethereumRecipient: string;
            assetAmount: import("..").IAssetAmount;
            feeAmount: import("..").IAssetAmount;
        }): Promise<import("@cosmjs/launchpad").Msg>;
        fetchUnconfirmedLockBurnTxs(address: string, confirmations: number): Promise<{
            readonly hash: string | undefined;
            readonly symbol: string | undefined;
            setTxHash(hash: string): void;
            emit(e: import("./EthbridgeService/types").TxEventPrepopulated<import("./EthbridgeService/types").TxEvent>): void;
            onTxEvent(handler: (e: import("./EthbridgeService/types").TxEvent) => void): any;
            onTxHash(handler: (e: import("./EthbridgeService/types").TxEventHashReceived) => void): any;
            onEthConfCountChanged(handler: (e: import("./EthbridgeService/types").TxEventEthConfCountChanged) => void): any;
            onEthTxConfirmed(handler: (e: import("./EthbridgeService/types").TxEventEthTxConfirmed) => void): any;
            onSifTxConfirmed(handler: (e: import("./EthbridgeService/types").TxEventSifTxConfirmed) => void): any;
            onEthTxInitiated(handler: (e: import("./EthbridgeService/types").TxEventEthTxInitiated) => void): any;
            onSifTxInitiated(handler: (e: import("./EthbridgeService/types").TxEventSifTxInitiated) => void): any;
            onComplete(handler: (e: import("./EthbridgeService/types").TxEventComplete) => void): any;
            removeListeners(): any;
            onError(handler: (e: import("./EthbridgeService/types").TxEventError) => void): any;
        }[]>;
        burnToSifchain(sifRecipient: string, assetAmount: import("..").IAssetAmount, confirmations: number, account?: string | undefined): {
            readonly hash: string | undefined;
            readonly symbol: string | undefined;
            setTxHash(hash: string): void;
            emit(e: import("./EthbridgeService/types").TxEventPrepopulated<import("./EthbridgeService/types").TxEvent>): void;
            onTxEvent(handler: (e: import("./EthbridgeService/types").TxEvent) => void): any;
            onTxHash(handler: (e: import("./EthbridgeService/types").TxEventHashReceived) => void): any;
            onEthConfCountChanged(handler: (e: import("./EthbridgeService/types").TxEventEthConfCountChanged) => void): any;
            onEthTxConfirmed(handler: (e: import("./EthbridgeService/types").TxEventEthTxConfirmed) => void): any;
            onSifTxConfirmed(handler: (e: import("./EthbridgeService/types").TxEventSifTxConfirmed) => void): any;
            onEthTxInitiated(handler: (e: import("./EthbridgeService/types").TxEventEthTxInitiated) => void): any;
            onSifTxInitiated(handler: (e: import("./EthbridgeService/types").TxEventSifTxInitiated) => void): any;
            onComplete(handler: (e: import("./EthbridgeService/types").TxEventComplete) => void): any;
            removeListeners(): any;
            onError(handler: (e: import("./EthbridgeService/types").TxEventError) => void): any;
        };
    };
    bus: {
        onAny(handler: import("./EventBusService").EventHandler): void;
        on(eventType: "ErrorEvent" | "TransactionErrorEvent" | "WalletConnectedEvent" | "WalletDisconnectedEvent" | "WalletConnectionErrorEvent" | "PegTransactionPendingEvent" | "PegTransactionCompletedEvent" | "PegTransactionErrorEvent" | "NoLiquidityPoolsFoundEvent" | import("./EventBusService").AppEventTypes, handler: import("./EventBusService").EventHandler): void;
        dispatch(event: import("./EventBusService").AppEvent): void;
    };
};
