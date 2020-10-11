import { BankBalances, BankSenderAndTxInfo, StdTx } from "../entities/Bank";
export declare const bankService: {
    getBalances(address: string): Promise<BankBalances>;
    sendCoins(address: string, account: BankSenderAndTxInfo): Promise<StdTx>;
};
