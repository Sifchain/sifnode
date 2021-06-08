import { TransactionStatus } from "../../entities";
import { AppEvent } from "../../services/EventBusService";
export declare function isSupportedEVMChain(chainId?: string): boolean;
export declare const ReportTransactionError: (bus: {
    dispatch: (event: AppEvent) => void;
}) => (txStatus: TransactionStatus) => TransactionStatus;
