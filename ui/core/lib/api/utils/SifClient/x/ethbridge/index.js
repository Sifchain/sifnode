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
exports.setupEthbridgeExtension = void 0;
function setupEthbridgeExtension(base) {
    return {
        ethbridge: {
            burn: (params) => __awaiter(this, void 0, void 0, function* () {
                console.log(`/ethbridge/burn`, JSON.stringify(params, null, 2));
                return yield base.post(`/ethbridge/burn`, params);
            }),
            lock: (params) => __awaiter(this, void 0, void 0, function* () {
                console.log(`/ethbridge/lock`, JSON.stringify(params, null, 2));
                return yield base.post(`/ethbridge/lock`, params);
            }),
        },
    };
}
exports.setupEthbridgeExtension = setupEthbridgeExtension;
//# sourceMappingURL=index.js.map