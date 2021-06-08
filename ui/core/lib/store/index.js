"use strict";
var __createBinding = (this && this.__createBinding) || (Object.create ? (function(o, m, k, k2) {
    if (k2 === undefined) k2 = k;
    Object.defineProperty(o, k2, { enumerable: true, get: function() { return m[k]; } });
}) : (function(o, m, k, k2) {
    if (k2 === undefined) k2 = k;
    o[k2] = m[k];
}));
var __exportStar = (this && this.__exportStar) || function(m, exports) {
    for (var p in m) if (p !== "default" && !exports.hasOwnProperty(p)) __createBinding(exports, m, p);
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.createStore = void 0;
const reactivity_1 = require("@vue/reactivity");
const wallet_1 = require("./wallet");
const asset_1 = require("./asset");
const pools_1 = require("./pools");
const tx_1 = require("./tx");
__exportStar(require("./poolFinder"), exports);
function createStore() {
    return reactivity_1.reactive({
        wallet: wallet_1.wallet,
        asset: asset_1.asset,
        pools: pools_1.pools,
        tx: tx_1.tx,
        accountpools: pools_1.accountpools,
    });
}
exports.createStore = createStore;
//# sourceMappingURL=index.js.map