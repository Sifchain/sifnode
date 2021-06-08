import { TransactionStatus } from "../../entities";
declare type TxEventBase<T> = {
    txHash: string;
    payload: T;
};
export declare type TxEventEthConfCountChanged = {
    type: "EthConfCountChanged";
} & TxEventBase<number>;
export declare type TxEventSifConfCountChanged = {
    type: "SifConfCountChanged";
} & TxEventBase<number>;
export declare type TxEventEthTxInitiated = {
    type: "EthTxInitiated";
} & TxEventBase<unknown>;
export declare type TxEventHashReceived = {
    type: "HashReceived";
} & TxEventBase<string>;
export declare type TxEventEthTxConfirmed = {
    type: "EthTxConfirmed";
} & TxEventBase<unknown>;
export declare type TxEventSifTxInitiated = {
    type: "SifTxInitiated";
} & TxEventBase<unknown>;
export declare type TxEventSifTxConfirmed = {
    type: "SifTxConfirmed";
} & TxEventBase<unknown>;
export declare type TxEventComplete = {
    type: "Complete";
} & TxEventBase<unknown>;
export declare type TxEventError = {
    type: "Error";
} & TxEventBase<TransactionStatus>;
export declare type TxEvent = TxEventEthConfCountChanged | TxEventSifConfCountChanged | TxEventEthTxInitiated | TxEventEthTxConfirmed | TxEventSifTxInitiated | TxEventSifTxConfirmed | TxEventHashReceived | TxEventError | TxEventComplete;
export declare type TxEventPrepopulated<T extends TxEvent = TxEvent> = Omit<T, "txHash"> & {
    txHash?: string;
};
export {};
