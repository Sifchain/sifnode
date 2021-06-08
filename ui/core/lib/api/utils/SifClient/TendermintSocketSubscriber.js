"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.createTendermintSocketSubscriber = exports.TendermintSocketSubscriber = void 0;
const reconnecting_websocket_1 = __importDefault(require("reconnecting-websocket"));
const eventemitter2_1 = require("eventemitter2");
const lodash_1 = require("lodash");
const ensureWs_1 = require("./ensureWs");
// Helper to allow us to add listeners to the open websocket
// In kind of a synchronous looking way
function openWebsocket(ws) {
    const wsPromise = new Promise((resolve) => {
        if (ws.readyState === reconnecting_websocket_1.default.OPEN) {
            resolve(ws);
            return;
        }
        ws.addEventListener("open", () => {
            resolve(ws);
        });
    });
    return wsPromise.then.bind(wsPromise);
}
// Simplify subscribing to Tendermintsocket
function TendermintSocketSubscriber({ wsUrl }) {
    const emitter = new eventemitter2_1.EventEmitter2();
    // This let's us wait until the websocket is open before subscribing to messages on it
    const _ws = ensureWs_1.ensureWs(wsUrl);
    const withWebsocket = openWebsocket(_ws);
    withWebsocket((ws) => {
        ws.addEventListener("message", (message) => {
            var _a;
            const data = JSON.parse(message.data);
            const eventData = (_a = data.result) === null || _a === void 0 ? void 0 : _a.data;
            if (!eventData)
                return;
            // Get last part of Tendermint Tx eg. 'tendermint/event/Tx'
            const [eventType] = eventData.type.split("/").slice(-1);
            // console.log("Message received");
            // console.log({ eventType, eventData });
            emitter.emit(eventType, eventData);
        });
    });
    return {
        on(event, handler) {
            // If for error listen immediately
            if (event === "error") {
                _ws.addEventListener("error", handler);
                return;
            }
            if (!emitter.hasListeners(event)) {
                withWebsocket((ws) => {
                    ws.send(JSON.stringify({
                        jsonrpc: "2.0",
                        method: "subscribe",
                        id: lodash_1.uniqueId(),
                        params: {
                            query: `tm.event='${event}'`,
                        },
                    }));
                });
            }
            emitter.on(event, handler);
        },
        off(event, handler) {
            emitter.off(event, handler);
        },
    };
}
exports.TendermintSocketSubscriber = TendermintSocketSubscriber;
function createTendermintSocketSubscriber(wsUrl) {
    return TendermintSocketSubscriber({ wsUrl });
}
exports.createTendermintSocketSubscriber = createTendermintSocketSubscriber;
//# sourceMappingURL=TendermintSocketSubscriber.js.map