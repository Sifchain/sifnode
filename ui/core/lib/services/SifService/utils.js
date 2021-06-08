"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.ensureSifAddress = void 0;
function ensureSifAddress(address) {
    if (address.length !== 42)
        throw "Address not valid (length). Fail"; // this is simple check, limited to default address type (check bech32);
    if (!address.match(/^sif/))
        throw "Address not valid (format). Fail"; // this is simple check, limited to default address type (check bech32);
    // TODO: add invariant address starts with "sif" (double check this is correct)
    return address;
}
exports.ensureSifAddress = ensureSifAddress;
//# sourceMappingURL=utils.js.map