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
const entities_1 = require("../../entities");
const reactivity_1 = require("@vue/reactivity");
const utils_1 = require("../utils");
exports.default = ({ services, store, }) => {
    const reportTransactionError = utils_1.ReportTransactionError(services.bus);
    const state = services.sif.getState();
    function syncPools() {
        return __awaiter(this, void 0, void 0, function* () {
            const state = services.sif.getState();
            // UPdate pools
            const pools = yield services.clp.getPools();
            for (let pool of pools) {
                store.pools[pool.symbol()] = pool;
            }
            // Update lp pools
            if (state.address) {
                const accountPoolSymbols = yield services.clp.getPoolSymbolsByLiquidityProvider(state.address);
                // This is a hot method when there are a heap of pools
                // Ideally we would have a better rest endpoint design
                accountPoolSymbols.forEach((symbol) => __awaiter(this, void 0, void 0, function* () {
                    const lp = yield services.clp.getLiquidityProvider({
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
                services.bus.dispatch({
                    type: "NoLiquidityPoolsFoundEvent",
                    payload: {},
                });
            }
        });
    });
    // Then every transaction
    services.sif.onNewBlock(() => __awaiter(void 0, void 0, void 0, function* () {
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
                const tx = yield services.clp.swap({
                    fromAddress: state.address,
                    sentAmount,
                    receivedAsset,
                    minimumReceived,
                });
                const txStatus = yield services.sif.signAndBroadcast(tx.value.msg);
                if (txStatus.state !== "accepted") {
                    // Edge case where we have run out of native balance and need to represent that
                    if (txStatus.code === entities_1.ErrorCode.TX_FAILED_USER_NOT_ENOUGH_BALANCE &&
                        sentAmount.symbol === "rowan") {
                        return reportTransactionError(Object.assign(Object.assign({}, txStatus), { code: entities_1.ErrorCode.TX_FAILED_NOT_ENOUGH_ROWAN_TO_COVER_GAS, memo: entities_1.getErrorMessage(entities_1.ErrorCode.TX_FAILED_NOT_ENOUGH_ROWAN_TO_COVER_GAS) }));
                    }
                    return reportTransactionError(txStatus);
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
                    ? services.clp.addLiquidity
                    : services.clp.createPool;
                const tx = yield provideLiquidity({
                    fromAddress: state.address,
                    nativeAssetAmount,
                    externalAssetAmount,
                });
                const txStatus = yield services.sif.signAndBroadcast(tx.value.msg);
                if (txStatus.state !== "accepted") {
                    // Edge case where we have run out of native balance and need to represent that
                    if (txStatus.code === entities_1.ErrorCode.TX_FAILED_USER_NOT_ENOUGH_BALANCE) {
                        return reportTransactionError(Object.assign(Object.assign({}, txStatus), { code: entities_1.ErrorCode.TX_FAILED_NOT_ENOUGH_ROWAN_TO_COVER_GAS, memo: entities_1.getErrorMessage(entities_1.ErrorCode.TX_FAILED_NOT_ENOUGH_ROWAN_TO_COVER_GAS) }));
                    }
                }
                return txStatus;
            });
        },
        removeLiquidity(asset, wBasisPoints, asymmetry) {
            return __awaiter(this, void 0, void 0, function* () {
                const tx = yield services.clp.removeLiquidity({
                    fromAddress: state.address,
                    asset,
                    asymmetry,
                    wBasisPoints,
                });
                const txStatus = yield services.sif.signAndBroadcast(tx.value.msg);
                if (txStatus.state !== "accepted") {
                    services.bus.dispatch({
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
                services.sif.purgeClient();
            });
        },
    };
    return actions;
};
//# sourceMappingURL=index.js.map