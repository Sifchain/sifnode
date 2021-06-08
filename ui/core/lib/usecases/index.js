"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.createUsecases = void 0;
const eth_1 = __importDefault(require("./wallet/eth"));
const clp_1 = __importDefault(require("./clp"));
const sif_1 = __importDefault(require("./wallet/sif"));
const peg_1 = __importDefault(require("./peg"));
function createUsecases(context) {
    return {
        clp: clp_1.default(context),
        wallet: {
            sif: sif_1.default(context),
            eth: eth_1.default(context),
        },
        peg: peg_1.default(context),
    };
}
exports.createUsecases = createUsecases;
//# sourceMappingURL=index.js.map