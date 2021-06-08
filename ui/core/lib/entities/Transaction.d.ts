import JSBI from "jsbi";
import { Address } from "./Address";
import { Asset } from "./Asset";
export declare type TransactionStatus = {
    code?: number;
    hash: string;
    state: "requested" | "accepted" | "failed" | "rejected" | "out_of_gas" | "completed";
    memo?: string;
    symbol?: string;
};
export declare type TxParams = {
    asset?: Asset;
    amount: JSBI;
    recipient: Address;
    feeRate?: number;
    memo?: string;
};
export declare type TxHash = string;
