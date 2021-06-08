import { AuthExtension, BroadcastMode, CosmosClient, LcdClient } from "@cosmjs/launchpad";
import { ClpExtension } from "./x/clp";
import { EthbridgeExtension } from "./x/ethbridge";
declare type CustomLcdClient = LcdClient & AuthExtension & ClpExtension & EthbridgeExtension;
declare type IClpApi = ClpExtension["clp"];
declare type IEthbridgeApi = EthbridgeExtension["ethbridge"];
declare type HandlerFn<T> = (a: T) => void;
export declare class SifUnSignedClient extends CosmosClient implements IClpApi, IEthbridgeApi {
    protected readonly lcdClient: CustomLcdClient;
    private subscriber;
    constructor(apiUrl: string, wsUrl?: string, rpcUrl?: string, broadcastMode?: BroadcastMode);
    swap: IClpApi["swap"];
    getPools: IClpApi["getPools"];
    getAssets: IClpApi["getAssets"];
    addLiquidity: IClpApi["addLiquidity"];
    createPool: IClpApi["createPool"];
    getLiquidityProvider: IClpApi["getLiquidityProvider"];
    removeLiquidity: IClpApi["removeLiquidity"];
    getPool: IClpApi["getPool"];
    burn: IEthbridgeApi["burn"];
    lock: IEthbridgeApi["lock"];
    onNewBlock<T>(handler: HandlerFn<T>): () => void;
    onTx<T>(handler: HandlerFn<T>): () => void;
    onSocketError<T>(handler: HandlerFn<T>): () => void;
}
export {};
