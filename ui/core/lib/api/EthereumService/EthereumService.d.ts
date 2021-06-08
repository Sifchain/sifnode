import { provider } from "web3-core";
import { IWalletService } from "../IWalletService";
import { TxHash, TxParams, Asset, IAssetAmount } from "../../entities";
import { EIPProvider } from "./utils/ethereumUtils";
import { Msg } from "@cosmjs/launchpad";
declare type Address = string;
declare type Balances = IAssetAmount[];
declare type PossibleProvider = provider | EIPProvider;
export declare type EthereumServiceContext = {
    getWeb3Provider: () => Promise<provider>;
    assets: Asset[];
};
export declare class EthereumService implements IWalletService {
    private web3;
    private supportedTokens;
    private blockSubscription;
    private provider;
    private providerPromise;
    private reportProviderNotFound;
    private chainIdDetectedHandler;
    private state;
    constructor(getWeb3Provider: () => Promise<PossibleProvider>, assets: Asset[]);
    onChainIdDetected(handler: (chainId: string) => void): void;
    onProviderNotFound(handler: () => void): void;
    getState(): {
        connected: boolean;
        address: string;
        accounts: string[];
        balances: IAssetAmount[];
        log: string;
    };
    private updateData;
    getAddress(): Address;
    isConnected(): boolean;
    getSupportedTokens(): import("../../entities").IAsset[];
    connect(): Promise<void>;
    addWeb3Subscription(): void;
    removeWeb3Subscription(): void;
    disconnect(): Promise<void>;
    getBalance(address?: Address, asset?: Asset): Promise<Balances>;
    transfer(params: TxParams): Promise<TxHash>;
    signAndBroadcast(msg: Msg, mmo?: string): Promise<void>;
    setPhrase(args: string): Promise<string>;
    purgeClient(): void;
    static create({ getWeb3Provider, assets, }: EthereumServiceContext): IWalletService;
}
declare const _default: typeof EthereumService.create;
export default _default;
