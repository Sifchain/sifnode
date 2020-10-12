import { ChainId } from "./ChainId";
import { Asset } from "./Asset";
export declare type Token = Asset & {
    chainId: ChainId;
    address: string;
};
export declare function createToken(chainId: ChainId, address: string, decimals: number, symbol: string, name: string): Token;
