import { TransactionStatus } from "../../entities";
export declare function parseTxFailure(txFailure: {
    transactionHash: string;
    rawLog: string;
}): TransactionStatus;
