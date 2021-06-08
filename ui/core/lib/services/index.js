"use strict";
var __createBinding = (this && this.__createBinding) || (Object.create ? (function(o, m, k, k2) {
    if (k2 === undefined) k2 = k;
    Object.defineProperty(o, k2, { enumerable: true, get: function() { return m[k]; } });
}) : (function(o, m, k, k2) {
    if (k2 === undefined) k2 = k;
    o[k2] = m[k];
}));
var __exportStar = (this && this.__exportStar) || function(m, exports) {
    for (var p in m) if (p !== "default" && !exports.hasOwnProperty(p)) __createBinding(exports, m, p);
};
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.createServices = void 0;
// Everything here represents services that are effectively remote data storage
__exportStar(require("./EthereumService/utils/getMetamaskProvider"), exports);
const EthereumService_1 = __importDefault(require("./EthereumService"));
const EthbridgeService_1 = __importDefault(require("./EthbridgeService"));
const SifService_1 = __importDefault(require("./SifService"));
const ClpService_1 = __importDefault(require("./ClpService"));
const EventBusService_1 = __importDefault(require("./EventBusService"));
function createServices(context) {
    const EthereumService = EthereumService_1.default(context);
    const EthbridgeService = EthbridgeService_1.default(context);
    const SifService = SifService_1.default(context);
    const ClpService = ClpService_1.default(context);
    const EventBusService = EventBusService_1.default(context);
    return {
        clp: ClpService,
        eth: EthereumService,
        sif: SifService,
        ethbridge: EthbridgeService,
        bus: EventBusService,
    };
}
exports.createServices = createServices;
//# sourceMappingURL=index.js.map