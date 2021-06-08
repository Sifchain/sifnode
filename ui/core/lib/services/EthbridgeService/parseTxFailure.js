"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.parseTxFailure = void 0;
const Errors_1 = require("../../entities/Errors");
// TODO: Should this go in a shared ethereum client mimicking sifchain?
function parseTxFailure({ hash = "", log = "", }) {
    // Ethereum events
    if (log.toString().toLowerCase().includes("request rejected") ||
        log.toString().toLowerCase().includes("user denied transaction")) {
        return {
            code: Errors_1.ErrorCode.USER_REJECTED,
            memo: Errors_1.getErrorMessage(Errors_1.ErrorCode.USER_REJECTED),
            hash,
            state: "rejected",
        };
    }
    return {
        code: Errors_1.ErrorCode.UNKNOWN_FAILURE,
        memo: Errors_1.getErrorMessage(Errors_1.ErrorCode.UNKNOWN_FAILURE),
        hash,
        state: "failed",
    };
}
exports.parseTxFailure = parseTxFailure;
//# sourceMappingURL=parseTxFailure.js.map