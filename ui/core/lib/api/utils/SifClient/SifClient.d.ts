import { BroadcastMode, CosmosFeeTable, GasLimits, GasPrice, OfflineSigner, SigningCosmosClient } from "@cosmjs/launchpad";
import { SifUnSignedClient } from "./SifUnsignedClient";
export declare class SifClient extends SigningCosmosClient {
    private wallet;
    private unsignedClient;
    constructor(apiUrl: string, senderAddress: string, signer: OfflineSigner, wsUrl: string, rpcUrl: string, gasPrice?: GasPrice, gasLimits?: Partial<GasLimits<CosmosFeeTable>>, broadcastMode?: BroadcastMode);
    getAccounts(): Promise<string[]>;
    getUnsignedClient(): SifUnSignedClient;
}
