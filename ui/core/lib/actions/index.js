"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.createActions = void 0;
const ethWallet_1 = __importDefault(require("./ethWallet"));
const clp_1 = __importDefault(require("./clp"));
const wallet_1 = __importDefault(require("./wallet"));
const peg_1 = __importDefault(require("./peg"));
function createActions(context) {
    return {
        ethWallet: ethWallet_1.default(context),
        clp: clp_1.default(context),
        wallet: wallet_1.default(context),
        peg: peg_1.default(context),
    };
}
exports.createActions = createActions;
//# sourceMappingURL=index.js.map