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
exports.SifClient = void 0;
const launchpad_1 = require("@cosmjs/launchpad");
const SifUnsignedClient_1 = require("./SifUnsignedClient");
class SifClient extends launchpad_1.SigningCosmosClient {
    constructor(apiUrl, senderAddress, signer, wsUrl, rpcUrl, gasPrice, gasLimits, broadcastMode) {
        super(apiUrl, senderAddress, signer, gasPrice, gasLimits, broadcastMode);
        this.wallet = signer;
        this.unsignedClient = new SifUnsignedClient_1.SifUnSignedClient(apiUrl, wsUrl, rpcUrl, broadcastMode);
    }
    getAccounts() {
        return __awaiter(this, void 0, void 0, function* () {
            const accounts = yield this.wallet.getAccounts();
            return accounts.map(({ address }) => address);
        });
    }
    getUnsignedClient() {
        return this.unsignedClient;
    }
}
exports.SifClient = SifClient;
//# sourceMappingURL=SifClient.js.map