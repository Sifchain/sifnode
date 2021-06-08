"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.isSupportedEVMChain = void 0;
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
//# sourceMappingURL=index.js.map