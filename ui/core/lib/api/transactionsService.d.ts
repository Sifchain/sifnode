import { BroadcastingResult, EncodedTransaction, Transaction } from "../entities/Transaction";
export declare const transactionService: {
    getByhash(hash: string): Promise<void>;
    search(actions: string, sender: string, page: number, limit: number, txheight: number): Promise<Transaction[]>;
    broadcast(tx: Transaction): Promise<BroadcastingResult>;
    encode(tx: Transaction): Promise<EncodedTransaction>;
    decode(tx: EncodedTransaction): Promise<Transaction>;
};
