import {
  AuthExtension,
  BroadcastMode,
  CosmosFeeTable,
  GasLimits,
  GasPrice,
  LcdClient,
  OfflineSigner,
  setupAuthExtension,
  SigningCosmosClient,
} from "@cosmjs/launchpad";

export class SifClient extends SigningCosmosClient {
  protected readonly lcdClient: LcdClient & AuthExtension;

  constructor(
    apiUrl: string,
    senderAddress: string,
    signer: OfflineSigner,
    gasPrice?: GasPrice,
    gasLimits?: Partial<GasLimits<CosmosFeeTable>>,
    broadcastMode?: BroadcastMode
  ) {
    super(apiUrl, senderAddress, signer, gasPrice, gasLimits, broadcastMode);
    this.lcdClient = LcdClient.withExtensions(
      { apiUrl: apiUrl, broadcastMode: broadcastMode },
      setupAuthExtension
    );
  }
}
