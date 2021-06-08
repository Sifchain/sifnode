"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.createPoolFinder = void 0;
const reactivity_1 = require("@vue/reactivity");
exports.createPoolFinder = (s) => (a, // externalAsset
b) => {
    const pools = reactivity_1.toRefs(s.pools);
    const key = [a, b]
        .map((x) => (typeof x === "string" ? x : x.symbol))
        .join("_");
    const poolRef = pools[key];
    return poolRef !== null && poolRef !== void 0 ? poolRef : null;
};
//# sourceMappingURL=poolFinder.js.map