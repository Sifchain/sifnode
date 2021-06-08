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
const reactivity_1 = require("@vue/reactivity");
exports.default = ({ api, store, }) => {
    const state = api.SifService.getState();
    function syncPools() {
        return __awaiter(this, void 0, void 0, function* () {
            const state = api.SifService.getState();
            // UPdate pools
            const pools = yield api.ClpService.getPools();
            for (let pool of pools) {
                store.pools[pool.symbol()] = pool;
            }
            // Update lp pools
            if (state.address) {
                const accountPoolSymbols = yield api.ClpService.getPoolSymbolsByLiquidityProvider(state.address);
                // This is a hot method when there are a heap of pools
                // Ideally we would have a better rest endpoint design
                accountPoolSymbols.forEach((symbol) => __awaiter(this, void 0, void 0, function* () {
                    const lp = yield api.ClpService.getLiquidityProvider({
                        symbol,
                        lpAddress: state.address,
                    });
                    if (!lp)
                        return;
                    const pool = `${symbol}_rowan`;
                    store.accountpools[state.address] =
                        store.accountpools[state.address] || {};
                    store.accountpools[state.address][pool] = { lp, pool };
                }));
                // Delete accountpools
                const currentPoolIds = accountPoolSymbols.map((id) => `${id}_rowan`);
                if (store.accountpools[state.address]) {
                    const existingPoolIds = Object.keys(store.accountpools[state.address]);
                    const disjunctiveIds = existingPoolIds.filter((id) => !currentPoolIds.includes(id));
                    disjunctiveIds.forEach((poolToRemove) => {
                        delete store.accountpools[state.address][poolToRemove];
                    });
                }
            }
        });
    }
    // Sync on load
    syncPools().then(() => {
        reactivity_1.effect(() => {
            if (Object.keys(store.pools).length === 0) {
                api.EventBusService.dispatch({
                    type: "NoLiquidityPoolsFoundEvent",
                    payload: {},
                });
            }
        });
    });
    // Then every transaction
    api.SifService.onNewBlock(() => __awaiter(void 0, void 0, void 0, function* () {
        yield syncPools();
    }));
    function findPool(pools, a, b) {
        var _a;
        const key = [a, b].sort().join("_");
        return (_a = pools[key]) !== null && _a !== void 0 ? _a : null;
    }
    reactivity_1.effect(() => {
        // When sif address changes syncPools
        store.wallet.sif.address;
        syncPools();
    });
    const actions = {
        swap(sentAmount, receivedAsset, minimumReceived) {
            return __awaiter(this, void 0, void 0, function* () {
                if (!state.address)
                    throw "No from address provided for swap";
                const tx = yield api.ClpService.swap({
                    fromAddress: state.address,
                    sentAmount,
                    receivedAsset,
                    minimumReceived,
                });
                const txStatus = yield api.SifService.signAndBroadcast(tx.value.msg);
                if (txStatus.state !== "accepted") {
                    api.EventBusService.dispatch({
                        type: "TransactionErrorEvent",
                        payload: {
                            txStatus,
                            message: txStatus.memo || "There was an error with your swap",
                        },
                    });
                }
                return txStatus;
            });
        },
        addLiquidity(nativeAssetAmount, externalAssetAmount) {
            return __awaiter(this, void 0, void 0, function* () {
                if (!state.address)
                    throw "No from address provided for swap";
                const hasPool = !!findPool(store.pools, nativeAssetAmount.asset.symbol, externalAssetAmount.asset.symbol);
                const provideLiquidity = hasPool
                    ? api.ClpService.addLiquidity
                    : api.ClpService.createPool;
                const tx = yield provideLiquidity({
                    fromAddress: state.address,
                    nativeAssetAmount,
                    externalAssetAmount,
                });
                const txStatus = yield api.SifService.signAndBroadcast(tx.value.msg);
                if (txStatus.state !== "accepted") {
                    api.EventBusService.dispatch({
                        type: "TransactionErrorEvent",
                        payload: {
                            txStatus,
                            message: txStatus.memo || "There was an error with your swap",
                        },
                    });
                }
                return txStatus;
            });
        },
        removeLiquidity(asset, wBasisPoints, asymmetry) {
            return __awaiter(this, void 0, void 0, function* () {
                const tx = yield api.ClpService.removeLiquidity({
                    fromAddress: state.address,
                    asset,
                    asymmetry,
                    wBasisPoints,
                });
                const txStatus = yield api.SifService.signAndBroadcast(tx.value.msg);
                if (txStatus.state !== "accepted") {
                    api.EventBusService.dispatch({
                        type: "TransactionErrorEvent",
                        payload: {
                            txStatus,
                            message: txStatus.memo || "There was an error removing liquidity",
                        },
                    });
                }
                return txStatus;
            });
        },
        disconnect() {
            return __awaiter(this, void 0, void 0, function* () {
                api.SifService.purgeClient();
            });
        },
    };
    return actions;
};
//# sourceMappingURL=index.js.map