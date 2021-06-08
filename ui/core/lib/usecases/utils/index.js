"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.ReportTransactionError = exports.isSupportedEVMChain = void 0;
function isSupportedEVMChain(chainId) {
    if (!chainId)
        return false;
    // List of supported EVM chainIds
    const supportedEVMChainIds = [
        "0x1",
        "0x3",
        "0x539",
    ];
    return supportedEVMChainIds.includes(chainId);
}
exports.isSupportedEVMChain = isSupportedEVMChain;
exports.ReportTransactionError = (bus) => (txStatus) => {
    bus.dispatch({
        type: "TransactionErrorEvent",
        payload: {
            txStatus,
            message: txStatus.memo || "There was an error with your swap",
        },
    });
    return txStatus;
};
//# sourceMappingURL=index.js.map