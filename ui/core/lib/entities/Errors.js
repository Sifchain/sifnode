"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.getErrorMessage = exports.ErrorCode = void 0;
var ErrorCode;
(function (ErrorCode) {
    ErrorCode[ErrorCode["TX_FAILED_SLIPPAGE"] = 0] = "TX_FAILED_SLIPPAGE";
    ErrorCode[ErrorCode["TX_FAILED"] = 1] = "TX_FAILED";
    ErrorCode[ErrorCode["USER_REJECTED"] = 2] = "USER_REJECTED";
    ErrorCode[ErrorCode["UNKNOWN_FAILURE"] = 3] = "UNKNOWN_FAILURE";
    ErrorCode[ErrorCode["INSUFFICIENT_FUNDS"] = 4] = "INSUFFICIENT_FUNDS";
    ErrorCode[ErrorCode["TX_FAILED_OUT_OF_GAS"] = 5] = "TX_FAILED_OUT_OF_GAS";
    ErrorCode[ErrorCode["TX_FAILED_NOT_ENOUGH_ROWAN_TO_COVER_GAS"] = 6] = "TX_FAILED_NOT_ENOUGH_ROWAN_TO_COVER_GAS";
    ErrorCode[ErrorCode["TX_FAILED_USER_NOT_ENOUGH_BALANCE"] = 7] = "TX_FAILED_USER_NOT_ENOUGH_BALANCE";
})(ErrorCode = exports.ErrorCode || (exports.ErrorCode = {}));
// This may be removed as it is a UX concern
const ErrorMessages = {
    [ErrorCode.TX_FAILED_SLIPPAGE]: "Your transaction has failed - Received amount is below expected",
    [ErrorCode.TX_FAILED]: "Your transaction has failed",
    [ErrorCode.USER_REJECTED]: "You have rejected the transaction",
    [ErrorCode.UNKNOWN_FAILURE]: "There was an unknown failure",
    [ErrorCode.INSUFFICIENT_FUNDS]: "You have insufficient funds",
    [ErrorCode.TX_FAILED_USER_NOT_ENOUGH_BALANCE]: "Not have enough balance",
    [ErrorCode.TX_FAILED_NOT_ENOUGH_ROWAN_TO_COVER_GAS]: "Not enough ROWAN to cover the gas fees",
    [ErrorCode.TX_FAILED_OUT_OF_GAS]: "Your transaction has failed - Out of gas",
};
function getErrorMessage(code) {
    return ErrorMessages[code];
}
exports.getErrorMessage = getErrorMessage;
//# sourceMappingURL=Errors.js.map