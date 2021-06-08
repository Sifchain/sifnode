"use strict";
// Consolodated place where we can setup testing services
var __awaiter = (this && this.__awaiter) || function (thisArg, _arguments, P, generator) {
    function adopt(value) { return value instanceof P ? value : new P(function (resolve) { resolve(value); }); }
    return new (P || (P = Promise))(function (resolve, reject) {
        function fulfilled(value) { try { step(generator.next(value)); } catch (e) { reject(e); } }
        function rejected(value) { try { step(generator["throw"](value)); } catch (e) { reject(e); } }
        function step(result) { result.done ? resolve(result.value) : adopt(result.value).then(fulfilled, rejected); }
        step((generator = generator.apply(thisArg, _arguments || [])).next());
    });
};
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.createTestEthService = exports.createTestSifService = void 0;
const SifService_1 = __importDefault(require("../../services/SifService"));
const EthereumService_1 = __importDefault(require("../../services/EthereumService"));
const getTestingToken_1 = require("./getTestingToken");
const getWeb3Provider_1 = require("./getWeb3Provider");
function createTestSifService(account) {
    return __awaiter(this, void 0, void 0, function* () {
        const sif = SifService_1.default({
            sifApiUrl: "http://localhost:1317",
            sifAddrPrefix: "sif",
            sifWsUrl: "ws://localhost:26657/websocket",
            sifRpcUrl: "http://localhost:26657",
            assets: getTestingToken_1.getTestingTokens(["CATK", "CBTK", "CETH", "ROWAN"]),
            keplrChainConfig: {},
        });
        if (account) {
            console.log("logging in to account with: " + account.mnemonic);
            yield sif.setPhrase(account.mnemonic);
        }
        return sif;
    });
}
exports.createTestSifService = createTestSifService;
function createTestEthService() {
    return __awaiter(this, void 0, void 0, function* () {
        const eth = EthereumService_1.default({
            assets: getTestingToken_1.getTestingTokens(["ATK", "BTK", "ETH", "EROWAN"]),
            getWeb3Provider: getWeb3Provider_1.getWeb3Provider,
        });
        console.log("Connecting to eth service");
        yield eth.connect();
        console.log("Finished connecting to eth service");
        return eth;
    });
}
exports.createTestEthService = createTestEthService;
//# sourceMappingURL=services.js.map