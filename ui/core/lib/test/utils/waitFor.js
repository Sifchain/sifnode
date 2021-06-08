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
exports.waitFor = void 0;
const sleep_1 = require("./sleep");
function waitFor(getter, expected, name) {
    return __awaiter(this, void 0, void 0, function* () {
        console.log(`Starting wait: "${name}" for value to be ${expected.toString()}`);
        let value;
        for (let i = 0; i < 100; i++) {
            yield sleep_1.sleep(1000);
            value = yield getter();
            console.log(`${value.toString()} ==? ${expected.toString()}`);
            if (value.toString() === expected.toString()) {
                return;
            }
        }
        throw new Error(`${value.toString()} never was ${expected.toString()} in wait: ${name}`);
    });
}
exports.waitFor = waitFor;
//# sourceMappingURL=waitFor.js.map