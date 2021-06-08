import { Msg } from "@cosmjs/launchpad";
import { Address, Asset, IAssetAmount, TransactionStatus, TxParams } from "../../entities";
import { Mnemonic } from "../../entities";
import { KeplrChainConfig } from "../../utils/parseConfig";
export declare type SifServiceContext = {
    sifAddrPrefix: string;
    sifApiUrl: string;
    sifWsUrl: string;
    sifRpcUrl: string;
    keplrChainConfig: KeplrChainConfig;
    assets: Asset[];
};
declare type HandlerFn<T> = (a: T) => void;
/**
 * Constructor for SifService
 *
 * SifService handles communication between our ui core Domain and the SifNode blockchain
 */
export default function createSifService({ sifAddrPrefix, sifApiUrl, sifWsUrl, sifRpcUrl, keplrChainConfig, assets, }: SifServiceContext): {
    /**
     * getState returns the service's reactive state to be listened to by consuming clients.
     */
    getState(): {
        connected: boolean;
        address: Address;
        accounts: Address[];
        balances: IAssetAmount[];
        log: string;
    };
    getSupportedTokens(): import("../../entities").IAsset[];
    setClient(): Promise<void>;
    initProvider(): Promise<void>;
    connect(): Promise<void>;
    isConnected(): boolean;
    onSocketError(handler: HandlerFn<any>): void;
    onTx(handler: HandlerFn<any>): void;
    onNewBlock(handler: HandlerFn<any>): void;
    setPhrase(mnemonic: Mnemonic): Promise<Address>;
    purgeClient(): Promise<void>;
    getBalance(address?: string | undefined, asset?: import("../../entities").IAsset | undefined): Promise<IAssetAmount[]>;
    transfer(params: TxParams): Promise<any>;
    signAndBroadcast(msg: Msg | Msg[], memo?: string | undefined): Promise<TransactionStatus>;
};
export {};
