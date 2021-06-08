"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.wallet = void 0;
const reactivity_1 = require("@vue/reactivity");
exports.wallet = reactivity_1.reactive({
    eth: {
        isConnected: false,
        address: "",
        balances: [],
    },
    sif: {
        isConnected: false,
        address: "",
        balances: [],
    },
});
//# sourceMappingURL=wallet.js.map