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
exports.getBridgeBankContract = void 0;
const json = require("../../../../../smart-contracts/build/contracts/BridgeBank.json");
function getBridgeBankContract(web3, address) {
    return __awaiter(this, void 0, void 0, function* () {
        return new web3.eth.Contract(json.abi, address);
    });
}
exports.getBridgeBankContract = getBridgeBankContract;
//# sourceMappingURL=bridgebankContract.js.map