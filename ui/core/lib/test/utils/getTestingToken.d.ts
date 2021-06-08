import { IAssetAmount } from "../../entities";
export declare function getTestingToken(tokenSymbol: string): import("../../entities").IAsset;
export declare function getTestingTokens(tokens: string[]): import("../../entities").IAsset[];
export declare function getBalance(balances: IAssetAmount[], symbol: string): IAssetAmount;
