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
exports.setupClpExtension = void 0;
function setupClpExtension(base) {
    return {
        clp: {
            getPools: () => __awaiter(this, void 0, void 0, function* () {
                var _a;
                return (_a = (yield base.get(`/clp/getPools`)).result) === null || _a === void 0 ? void 0 : _a.Pools;
            }),
            getAssets: (address) => __awaiter(this, void 0, void 0, function* () {
                return (yield base.get(`/clp/getAssets?lpAddress=${address}`)).result;
            }),
            swap: (params) => __awaiter(this, void 0, void 0, function* () {
                return yield base.post(`/clp/swap`, params);
            }),
            addLiquidity: (params) => __awaiter(this, void 0, void 0, function* () {
                return yield base.post(`/clp/addLiquidity`, params);
            }),
            createPool: (params) => __awaiter(this, void 0, void 0, function* () {
                return yield base.post(`/clp/createPool`, params);
            }),
            getLiquidityProvider: ({ symbol, lpAddress }) => __awaiter(this, void 0, void 0, function* () {
                return yield base.get(`/clp/getLiquidityProvider?symbol=${symbol}&lpAddress=${lpAddress}`);
            }),
            removeLiquidity: (params) => __awaiter(this, void 0, void 0, function* () {
                return yield base.post(`/clp/removeLiquidity`, params);
            }),
            getPool: ({ ticker }) => __awaiter(this, void 0, void 0, function* () {
                return (yield base.get(`/clp/getPool?ticker=${ticker}`)).result;
            }),
        },
    };
}
exports.setupClpExtension = setupClpExtension;
//# sourceMappingURL=index.js.map