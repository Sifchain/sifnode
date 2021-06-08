import { Network } from "./Network";
export declare function Token(p: {
    address: string;
    decimals: number;
    imageUrl?: string;
    name: string;
    network: Network;
    symbol: string;
}): {
    address: string;
    decimals: number;
    imageUrl?: string | undefined;
    name: string;
    network: Network;
    symbol: string;
};
export declare type Token = ReturnType<typeof Token>;
