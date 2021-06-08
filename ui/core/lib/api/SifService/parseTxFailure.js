"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.parseTxFailure = void 0;
const Errors_1 = require("../../entities/Errors");
function parseTxFailure(txFailure) {
    console.log({ "txFailure.rawLog": txFailure.rawLog });
    // TODO: synchronise with backend and use error codes at the service level
    // and provide a localized error lookup on frontend for messages
    if (txFailure.rawLog.toLowerCase().includes("below expected")) {
        return {
            code: Errors_1.ErrorCode.TX_FAILED_SLIPPAGE,
            hash: txFailure.transactionHash,
            memo: Errors_1.getErrorMessage(Errors_1.ErrorCode.TX_FAILED_SLIPPAGE),
            state: "failed",
        };
    }
    if (txFailure.rawLog.toLowerCase().includes("swap_failed")) {
        return {
            code: Errors_1.ErrorCode.TX_FAILED,
            hash: txFailure.transactionHash,
            memo: Errors_1.getErrorMessage(Errors_1.ErrorCode.TX_FAILED),
            state: "failed",
        };
    }
    if (txFailure.rawLog.toLowerCase().includes("request rejected")) {
        return {
            code: Errors_1.ErrorCode.USER_REJECTED,
            hash: txFailure.transactionHash,
            memo: Errors_1.getErrorMessage(Errors_1.ErrorCode.USER_REJECTED),
            state: "rejected",
        };
    }
    if (txFailure.rawLog.toLowerCase().includes("out of gas")) {
        return {
            code: Errors_1.ErrorCode.TX_FAILED_OUT_OF_GAS,
            hash: txFailure.transactionHash,
            memo: Errors_1.getErrorMessage(Errors_1.ErrorCode.TX_FAILED_OUT_OF_GAS),
            state: "out_of_gas",
        };
    }
    if (txFailure.rawLog.toLowerCase().includes("insufficient funds")) {
        return {
            code: Errors_1.ErrorCode.INSUFFICIENT_FUNDS,
            hash: txFailure.transactionHash,
            memo: Errors_1.getErrorMessage(Errors_1.ErrorCode.INSUFFICIENT_FUNDS),
            state: "failed",
        };
    }
    return {
        code: Errors_1.ErrorCode.UNKNOWN_FAILURE,
        hash: txFailure.transactionHash,
        memo: Errors_1.getErrorMessage(Errors_1.ErrorCode.UNKNOWN_FAILURE),
        state: "failed",
    };
}
exports.parseTxFailure = parseTxFailure;
//# sourceMappingURL=parseTxFailure.js.map