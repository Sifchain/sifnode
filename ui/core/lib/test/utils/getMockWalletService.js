"use strict";
var __awaiter = (this && this.__awaiter) || function (thisArg, _arguments, P, generator) {
    function adopt(value) { return value instanceof P ? value : new P(function (resolve) { resolve(value); }); }
    return new (P || (P = Promise))(function (resolve, reject) {
        function fulfilled(value) { try { step(generator.next(value)); } catch (e) { reject(e); } }
        function rejected(value) { try { step(generator["throw"](value)); } catch (e) { reject(e); } }
        function step(result) { result.done ? resolve(result.value) : adopt(result.value).then(fulfilled, rejected); }
        step((generator = generator.apply(thisArg, _arguments || [])).next());
    });
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.getMockWalletService = void 0;
function getMockWalletService(state, walletBalances, service = {}) {
    return Object.assign(Object.assign({ setPhrase: () => __awaiter(this, void 0, void 0, function* () { return ""; }), purgeClient: () => { }, getState: () => state, transfer: () => __awaiter(this, void 0, void 0, function* () { return ""; }), getBalance: jest.fn(() => __awaiter(this, void 0, void 0, function* () { return walletBalances; })), getSupportedTokens: () => [], connect: jest.fn(() => __awaiter(this, void 0, void 0, function* () {
            state.connected = true;
            state.balances = walletBalances;
        })), disconnect: jest.fn(() => __awaiter(this, void 0, void 0, function* () {
            state.connected = false;
        })), isConnected: () => true }, service), { signAndBroadcast: (msg, memo) => __awaiter(this, void 0, void 0, function* () { }), onProviderNotFound: () => { }, onChainIdDetected: () => { } });
}
exports.getMockWalletService = getMockWalletService;
//# sourceMappingURL=getMockWalletService.js.map