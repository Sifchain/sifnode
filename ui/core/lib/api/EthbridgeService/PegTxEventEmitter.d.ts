import { TxEvent, TxEventComplete, TxEventError, TxEventEthConfCountChanged, TxEventEthTxConfirmed, TxEventEthTxInitiated, TxEventHashReceived, TxEventPrepopulated, TxEventSifTxConfirmed, TxEventSifTxInitiated } from "./types";
export declare type PegTxEventEmitter = ReturnType<typeof createPegTxEventEmitter>;
/**
 * Adds types around EventEmitter2
 * @param txHash transaction hash this emitter responds to
 */
export declare function createPegTxEventEmitter(txHash?: string, symbol?: string): {
    readonly hash: string | undefined;
    readonly symbol: string | undefined;
    setTxHash(hash: string): void;
    emit(e: TxEventPrepopulated): void;
    onTxEvent(handler: (e: TxEvent) => void): any;
    onTxHash(handler: (e: TxEventHashReceived) => void): any;
    onEthConfCountChanged(handler: (e: TxEventEthConfCountChanged) => void): any;
    onEthTxConfirmed(handler: (e: TxEventEthTxConfirmed) => void): any;
    onSifTxConfirmed(handler: (e: TxEventSifTxConfirmed) => void): any;
    onEthTxInitiated(handler: (e: TxEventEthTxInitiated) => void): any;
    onSifTxInitiated(handler: (e: TxEventSifTxInitiated) => void): any;
    onComplete(handler: (e: TxEventComplete) => void): any;
    removeListeners(): any;
    onError(handler: (e: TxEventError) => void): any;
};
