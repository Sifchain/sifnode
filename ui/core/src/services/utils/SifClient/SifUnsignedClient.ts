import {
  AuthExtension,
  BroadcastMode,
  CosmosClient,
  LcdClient,
  setupAuthExtension,
} from "@cosmjs/launchpad";

import {
  createTendermintSocketPoll,
  TendermintSocketPoll,
} from "./TendermintSocketPoll";

import { ClpExtension, setupClpExtension } from "./x/clp";
import { IDispensationApi, setupDispensationApi } from "./x/dispensation";
import { EthbridgeExtension, setupEthbridgeExtension } from "./x/ethbridge";

type CustomLcdClient = LcdClient &
  AuthExtension &
  ClpExtension &
  EthbridgeExtension &
  IDispensationApi;

function createLcdClient(
  apiUrl: string,
  broadcastMode: BroadcastMode | undefined,
): CustomLcdClient {
  return LcdClient.withExtensions(
    { apiUrl: apiUrl, broadcastMode: broadcastMode },
    setupAuthExtension,
    setupClpExtension,
    setupEthbridgeExtension,
    setupDispensationApi,
  );
}

type IClpApi = ClpExtension["clp"];
type IEthbridgeApi = EthbridgeExtension["ethbridge"];

type HandlerFn<T> = (a: T) => void;
export class SifUnSignedClient
  extends CosmosClient
  implements IClpApi, IEthbridgeApi {
  protected readonly lcdClient: CustomLcdClient;
  private subscriber: TendermintSocketPoll | undefined;
  constructor(
    apiUrl: string,
    wsUrl = "ws://localhost:26657/websocket",
    rpcUrl = "http://localhost:26657",
    broadcastMode?: BroadcastMode,
  ) {
    super(apiUrl, broadcastMode);
    this.lcdClient = createLcdClient(apiUrl, broadcastMode);
    this.swap = this.lcdClient.clp.swap;
    this.getPools = this.lcdClient.clp.getPools;
    this.getAssets = this.lcdClient.clp.getAssets;
    this.addLiquidity = this.lcdClient.clp.addLiquidity;
    this.createPool = this.lcdClient.clp.createPool;
    this.getLiquidityProvider = this.lcdClient.clp.getLiquidityProvider;
    this.removeLiquidity = this.lcdClient.clp.removeLiquidity;
    this.getPool = this.lcdClient.clp.getPool;
    this.burn = this.lcdClient.ethbridge.burn;
    this.lock = this.lcdClient.ethbridge.lock;
    this.claim = this.lcdClient.dispensation.claim;
    this.subscriber = createTendermintSocketPoll(rpcUrl);
  }

  // Clp Extension
  swap: IClpApi["swap"];
  getPools: IClpApi["getPools"];
  getAssets: IClpApi["getAssets"];
  addLiquidity: IClpApi["addLiquidity"];
  createPool: IClpApi["createPool"];
  getLiquidityProvider: IClpApi["getLiquidityProvider"];
  removeLiquidity: IClpApi["removeLiquidity"];
  getPool: IClpApi["getPool"];

  // Ethbridge Extension
  burn: IEthbridgeApi["burn"];
  lock: IEthbridgeApi["lock"];

  // Dispensation
  claim: IDispensationApi["dispensation"]["claim"];

  onNewBlock<T>(handler: HandlerFn<T>) {
    console.log("received onNewBlock handler");
    if (!this.subscriber) console.error("Subscriber not setup");
    this.subscriber?.on("NewBlock", handler);
    return () => {
      this.subscriber?.off("NewBlock", handler);
    };
  }

  onTx<T>(handler: HandlerFn<T>) {
    console.log("received onTx handler");
    if (!this.subscriber) console.error("Subscriber not setup");
    this.subscriber?.on("Tx", handler);
    return () => {
      this.subscriber?.off("Tx", handler);
    };
  }

  onSocketError<T>(handler: HandlerFn<T>) {
    console.log("received onSocketError handler");
    if (!this.subscriber) console.error("Subscriber not setup");
    this.subscriber?.on("error", handler);
    return () => {
      this.subscriber?.off("error", handler);
    };
  }
}
