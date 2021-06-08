import { TestSifAccount } from "./accounts";
export declare function createTestSifService(account?: TestSifAccount): Promise<{
    getState(): {
        connected: boolean;
        address: string;
        accounts: string[];
        balances: import("../..").IAssetAmount[];
        log: string;
    };
    getSupportedTokens(): import("../..").IAsset[];
    setClient(): Promise<void>;
    initProvider(): Promise<void>;
    connect(): Promise<void>;
    isConnected(): boolean;
    onSocketError(handler: (a: any) => void): void;
    onTx(handler: (a: any) => void): void;
    onNewBlock(handler: (a: any) => void): void;
    setPhrase(mnemonic: string): Promise<string>;
    purgeClient(): Promise<void>;
    getBalance(address?: string | undefined, asset?: import("../..").IAsset | undefined): Promise<import("../..").IAssetAmount[]>;
    transfer(params: import("../..").TxParams): Promise<any>;
    signAndBroadcast(msg: import("@cosmjs/launchpad").Msg | import("@cosmjs/launchpad").Msg[], memo?: string | undefined): Promise<import("../..").TransactionStatus>;
}>;
export declare function createTestEthService(): Promise<import("../../services/IWalletService").IWalletService>;
