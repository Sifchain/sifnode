"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.ensureWs = void 0;
const reconnecting_websocket_1 = __importDefault(require("reconnecting-websocket"));
// Pool socket connections
const wsPool = {};
function ensureWs(wsUrl) {
    if (!wsPool[wsUrl]) {
        wsPool[wsUrl] = new reconnecting_websocket_1.default(wsUrl);
    }
    return wsPool[wsUrl];
}
exports.ensureWs = ensureWs;
//# sourceMappingURL=ensureWs.js.map