import { PegTxEventEmitter, TxEventError, TxEventPrepopulated } from "./types";
import { EventEmitter2 } from "eventemitter2";
/**
 * Adds types around EventEmitter2
 * @param txHash transaction hash this emitter responds to
 */
export function createPegTxEventEmitter(txHash?: string) {
  let _txHash = txHash;
  const emitter = new EventEmitter2();

  const instance: PegTxEventEmitter = {
    setTxHash(hash: string) {
      _txHash = hash;
      this.emit({ type: "HashReceived", payload: hash });
    },
    emit(e: TxEventPrepopulated) {
      emitter.emit(e.type, { ...e, txHash: e.txHash || _txHash });
    },
    onTxEvent(handler) {
      emitter.onAny((e, v) => handler(v));
      return instance;
    },
    onTxHash(handler) {
      emitter.on("HashReceived", handler);
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
    onError(handler) {
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
