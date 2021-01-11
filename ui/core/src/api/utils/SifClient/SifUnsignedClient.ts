import {
  AuthExtension,
  BroadcastMode,
  CosmosClient,
  LcdClient,
  setupAuthExtension,
} from "@cosmjs/launchpad";

import {
  createTendermintSocketSubscriber,
  TendermintSocketSubscriber,
} from "./TendermintSocketSubscriber";

import { ClpExtension, setupClpExtension } from "./x/clp";
import { EthbridgeExtension, setupEthbridgeExtension } from "./x/ethbridge";

type CustomLcdClient = LcdClient &
  AuthExtension &
  ClpExtension &
  EthbridgeExtension;

function createLcdClient(
  apiUrl: string,
  broadcastMode: BroadcastMode | undefined
): CustomLcdClient {
  return LcdClient.withExtensions(
    { apiUrl: apiUrl, broadcastMode: broadcastMode },
    setupAuthExtension,
    setupClpExtension,
    setupEthbridgeExtension
  );
}

type IClpApi = ClpExtension["clp"];
type IEthbridgeApi = EthbridgeExtension["ethbridge"];

type HandlerFn<T> = (a: T) => void;
export class SifUnSignedClient extends CosmosClient
  implements IClpApi, IEthbridgeApi {
  protected readonly lcdClient: CustomLcdClient;
  private subscriber: TendermintSocketSubscriber | undefined;
  constructor(
    apiUrl: string,
    wsUrl: string = "ws://localhost:26657/websocket",
    broadcastMode?: BroadcastMode
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
    if (wsUrl) this.subscriber = createTendermintSocketSubscriber(wsUrl);
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

  onNewBlock<T>(handler: HandlerFn<T>) {
    this.subscriber?.on("NewBlock", handler);
  }

  onTx<T>(handler: HandlerFn<T>) {
    this.subscriber?.on("Tx", handler);
  }

  onSocketError<T>(handler: HandlerFn<T>) {
    this.subscriber?.on("error", handler);
  }
}
