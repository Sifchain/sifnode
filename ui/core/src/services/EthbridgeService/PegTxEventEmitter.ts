import {
  TxEvent,
  TxEventComplete,
  TxEventError,
  TxEventEthConfCountChanged,
  TxEventEthTxConfirmed,
  TxEventEthTxInitiated,
  TxEventHashReceived,
  TxEventPrepopulated,
  TxEventSifTxConfirmed,
  TxEventSifTxInitiated,
} from "./types";
import { EventEmitter2 } from "eventemitter2";
export type PegTxEventEmitter = ReturnType<typeof createPegTxEventEmitter>;
/**
 * Adds types around EventEmitter2
 * @param txHash transaction hash this emitter responds to
 */
export function createPegTxEventEmitter(txHash?: string, symbol?: string) {
  let _txHash = txHash;
  let _symbol = symbol;
  const emitter = new EventEmitter2();

  const instance = {
    get hash() {
      return _txHash;
    },
    get symbol() {
      return _symbol;
    },
    setTxHash(hash: string) {
      _txHash = hash;
      this.emit({ type: "HashReceived", payload: hash });
    },
    emit(e: TxEventPrepopulated) {
      emitter.emit(e.type, { ...e, txHash: e.txHash || _txHash });
    },
    onTxEvent(handler: (e: TxEvent) => void) {
      emitter.onAny((e, v) => handler(v));
      return instance;
    },
    onTxHash(handler: (e: TxEventHashReceived) => void) {
      emitter.on("HashReceived", handler);
      return instance;
    },
    onEthConfCountChanged(handler: (e: TxEventEthConfCountChanged) => void) {
      emitter.on("EthConfCountChanged", handler);
      return instance;
    },
    onEthTxConfirmed(handler: (e: TxEventEthTxConfirmed) => void) {
      emitter.on("EthTxConfirmed", handler);
      return instance;
    },

    onSifTxConfirmed(handler: (e: TxEventSifTxConfirmed) => void) {
      emitter.on("SifTxConfirmed", handler);
      return instance;
    },
    onEthTxInitiated(handler: (e: TxEventEthTxInitiated) => void) {
      emitter.on("EthTxInitiated", handler);
      return instance;
    },
    onSifTxInitiated(handler: (e: TxEventSifTxInitiated) => void) {
      emitter.on("SifTxInitiated", handler);
      return instance;
    },
    onComplete(handler: (e: TxEventComplete) => void) {
      emitter.on("Complete", handler);
      return instance;
    },
    removeListeners() {
      emitter.removeAllListeners();
      return instance;
    },
    onError(handler: (e: TxEventError) => void) {
      emitter.on("Error", (e: TxEventError) => {
        handler(e);

        // We assume the yx is in an error state
        // so dont want the listener to transmit
        // events after an error
        emitter.removeAllListeners();
      });

      return instance;
    },
  };

  return instance;
}
