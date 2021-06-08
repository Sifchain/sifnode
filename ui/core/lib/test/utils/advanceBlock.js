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
exports.advanceBlock = void 0;
const getWeb3Provider_1 = require("./getWeb3Provider");
// No TS defs yet provided https://github.com/OpenZeppelin/openzeppelin-test-helpers/pull/141
const { time } = require("@openzeppelin/test-helpers");
beforeEach(() => __awaiter(void 0, void 0, void 0, function* () {
    require("@openzeppelin/test-helpers/configure")({
        provider: yield getWeb3Provider_1.getWeb3Provider(),
    });
}));
function advanceBlock(count) {
    return __awaiter(this, void 0, void 0, function* () {
        console.log("Advancing time by " + count + " blocks");
        for (let i = 0; i < count; i++) {
            yield time.advanceBlock();
        }
        console.log("Finished advancing time.");
    });
}
exports.advanceBlock = advanceBlock;
//# sourceMappingURL=advanceBlock.js.map