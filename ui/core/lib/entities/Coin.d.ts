import { Network } from "./Network";
export declare function Coin(p: {
    decimals: number;
    imageUrl?: string;
    name: string;
    network: Network;
    symbol: string;
}): {
    decimals: number;
    imageUrl?: string | undefined;
    name: string;
    network: Network;
    symbol: string;
};
export declare type Coin = ReturnType<typeof Coin>;
