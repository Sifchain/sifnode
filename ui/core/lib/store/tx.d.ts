import { TransactionStatus } from "../entities";
export declare type TxStore = {
    eth: {
        [address: string]: {
            [hash: string]: TransactionStatus;
        };
    };
};
export declare const tx: TxStore;
