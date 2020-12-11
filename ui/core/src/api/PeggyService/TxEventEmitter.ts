import { TxEventEmitter, TxEventPrepopulated } from "./types";
import { EventEmitter2 } from "eventemitter2";
/**
 * Adds types around EventEmitter2
 * @param txHash transaction hash this emitter responds to
 */
export function createTxEventEmitter(txHash: string) {
  const emitter = new EventEmitter2();
  const instance: TxEventEmitter = {
    emit(e: TxEventPrepopulated) {
      emitter.emit(e.type, { ...e, txHash: e.txHash || txHash });
    },
    onTxEvent(handler) {
      emitter.onAny((e, v) => handler(v));
      return instance;
    },
    onEthConfCountChanged(handler) {
      emitter.on("EthConfCountChanged", handler);
      return instance;
    },
    onEthTxConfirmed(handler) {
      emitter.on("EthTxConfirmed", handler);
      return instance;
    },

    onSifTxConfirmed(handler) {
      emitter.on("SifTxConfirmed", handler);
      return instance;
    },
    onEthTxInitiated(handler) {
      emitter.on("EthTxInitiated", handler);
      return instance;
    },
    onSifTxInitiated(handler) {
      emitter.on("SifTxInitiated", handler);
      return instance;
    },
    onComplete(handler) {
      emitter.on("Complete", handler);
      return instance;
    },
  };
  return instance;
}
