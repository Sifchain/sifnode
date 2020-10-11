export declare type Asset = {
    decimals: number;
    symbol: string;
    name: string;
};
export declare function createAsset(decimals: number, symbol: string, name: string): Asset;
