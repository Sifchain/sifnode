import {
  BroadcastMode,
  CosmosFeeTable,
  GasLimits,
  GasPrice,
  OfflineSigner,
  SigningCosmosClient,
} from "@cosmjs/launchpad";
import { SifUnSignedClient } from "./SifUnsignedClient";

export class SifClient extends SigningCosmosClient {
  private wallet: OfflineSigner;
  private unsignedClient: SifUnSignedClient;

  constructor(
    apiUrl: string,
    senderAddress: string,
    signer: OfflineSigner,
    wsUrl: string,
    rpcUrl: string,
    gasPrice?: GasPrice,
    gasLimits?: Partial<GasLimits<CosmosFeeTable>>,
    broadcastMode?: BroadcastMode,
  ) {
    super(apiUrl, senderAddress, signer, gasPrice, gasLimits, broadcastMode);
    this.wallet = signer;
    this.unsignedClient = new SifUnSignedClient(
      apiUrl,
      wsUrl,
      rpcUrl,
      broadcastMode,
    );
  }

  async getAccounts(): Promise<string[]> {
    const accounts = await this.wallet.getAccounts();
    return accounts.map(({ address }) => address);
  }

  getUnsignedClient() {
    return this.unsignedClient;
  }
}
