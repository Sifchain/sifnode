"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.createPegTxEventEmitter = void 0;
const eventemitter2_1 = require("eventemitter2");
/**
 * Adds types around EventEmitter2
 * @param txHash transaction hash this emitter responds to
 */
function createPegTxEventEmitter(txHash, symbol) {
    let _txHash = txHash;
    let _symbol = symbol;
    const emitter = new eventemitter2_1.EventEmitter2();
    const instance = {
        get hash() {
            return _txHash;
        },
        get symbol() {
            return _symbol;
        },
        setTxHash(hash) {
            _txHash = hash;
            this.emit({ type: "HashReceived", payload: hash });
        },
        emit(e) {
            emitter.emit(e.type, Object.assign(Object.assign({}, e), { txHash: e.txHash || _txHash }));
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
        removeListeners() {
            emitter.removeAllListeners();
            return instance;
        },
        onError(handler) {
            emitter.on("Error", (e) => {
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
exports.createPegTxEventEmitter = createPegTxEventEmitter;
//# sourceMappingURL=PegTxEventEmitter.js.map