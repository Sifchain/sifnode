import { TransactionStatus } from "../../entities";
declare type WalletType = "sif" | "eth";
declare type ErrorEvent = {
    type: "ErrorEvent";
    payload: {
        message: string;
        detail?: {
            type: "etherscan" | "info";
            message: string;
        };
    };
};
declare type TransactionErrorEvent = {
    type: "TransactionErrorEvent";
    payload: {
        txStatus: TransactionStatus;
        message: string;
    };
};
declare type WalletConnectedEvent = {
    type: "WalletConnectedEvent";
    payload: {
        walletType: WalletType;
        address: string;
    };
};
declare type WalletDisconnectedEvent = {
    type: "WalletDisconnectedEvent";
    payload: {
        walletType: WalletType;
        address: string;
    };
};
declare type WalletConnectionErrorEvent = {
    type: "WalletConnectionErrorEvent";
    payload: {
        walletType: WalletType;
        message: string;
    };
};
declare type PegTransactionPendingEvent = {
    type: "PegTransactionPendingEvent";
    payload: {
        hash: string;
    };
};
declare type PegTransactionCompletedEvent = {
    type: "PegTransactionCompletedEvent";
    payload: {
        hash: string;
    };
};
declare type PegTransactionErrorEvent = {
    type: "PegTransactionErrorEvent";
    payload: {
        txStatus: TransactionStatus;
        message: string;
    };
};
declare type NoLiquidityPoolsFoundEvent = {
    type: "NoLiquidityPoolsFoundEvent";
    payload: {};
};
export declare type AppEvent = ErrorEvent | WalletConnectedEvent | WalletDisconnectedEvent | WalletConnectionErrorEvent | PegTransactionPendingEvent | PegTransactionCompletedEvent | NoLiquidityPoolsFoundEvent | TransactionErrorEvent | PegTransactionErrorEvent;
export {};
