import { ServiceContext } from "./services";
export declare type AppConfig = ServiceContext;
export declare function getConfig(config?: string, sifchainAssetTag?: string, ethereumAssetTag?: string): AppConfig;
