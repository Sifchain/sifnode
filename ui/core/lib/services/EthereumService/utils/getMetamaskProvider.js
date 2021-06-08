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
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.getMetamaskProvider = void 0;
const detect_provider_1 = __importDefault(require("@metamask/detect-provider"));
// Detect mossible metamask provider from browser
exports.getMetamaskProvider = () => __awaiter(void 0, void 0, void 0, function* () {
    const mmp = yield detect_provider_1.default();
    const win = window;
    if (!mmp || !win)
        return null;
    if (mmp) {
        return mmp;
    }
    // if a wallet has left web3 on the page we can use the current provider
    if (win.web3) {
        return win.web3.currentProvider;
    }
    return null;
});
//# sourceMappingURL=getMetamaskProvider.js.map