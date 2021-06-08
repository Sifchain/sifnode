import { TransactionStatus } from "../../entities";
export declare function parseTxFailure({ hash, log, }: {
    hash: string;
    log: string;
}): TransactionStatus;
