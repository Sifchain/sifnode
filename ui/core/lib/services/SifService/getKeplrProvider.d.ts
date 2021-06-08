import { OfflineSigner } from "@cosmjs/launchpad";
declare type Keplr = {
    experimentalSuggestChain?(chainInfo: any): Promise<any>;
    enable(chainId: string): Promise<void>;
    getKey(chainId: string): Promise<any>;
    getTxConfig(chainId: string, config: any): Promise<any>;
    sign(chainId: string, signer: string, message: Uint8Array): Promise<any>;
    sendTx(chainId: string, stdTx: any, mode: any): Promise<any>;
    suggestToken(chainId: string, contractAddress: string): Promise<void>;
    requestTx(chainId: string, txBytes: Uint8Array, mode: "sync" | "async" | "commit", isRestAPI: boolean): Promise<void>;
    requestTxWithResult(chainId: string, txBytes: Uint8Array, mode: "sync" | "async" | "commit", isRestAPI: boolean): Promise<any>;
    getSecret20ViewingKey(chainId: string, contractAddress: string): Promise<string>;
    getOfflineSigner: (chainId?: string) => OfflineSigner;
};
declare type provider = Keplr;
export default function getKeplrProvider(): Promise<provider | null>;
export {};
