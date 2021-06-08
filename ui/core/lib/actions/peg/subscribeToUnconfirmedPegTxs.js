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
exports.SubscribeToUnconfirmedPegTxs = void 0;
const subscribeToTx_1 = require("./utils/subscribeToTx");
exports.SubscribeToUnconfirmedPegTxs = ({ api, store, config, }) => () => {
    // Update a tx state in the store
    const subscribeToTx = subscribeToTx_1.SubscribeToTx({ store, api });
    function getSubscriptions() {
        return __awaiter(this, void 0, void 0, function* () {
            const pendingTxs = yield api.EthbridgeService.fetchUnconfirmedLockBurnTxs(store.wallet.eth.address, config.ethConfirmations);
            return pendingTxs.map(subscribeToTx);
        });
    }
    // Need to keep subscriptions syncronous so using promise
    const subscriptionsPromise = getSubscriptions();
    // Return unsubscribe synchronously
    return () => {
        subscriptionsPromise.then((subscriptions) => subscriptions.forEach((unsubscribe) => unsubscribe()));
    };
};
//# sourceMappingURL=subscribeToUnconfirmedPegTxs.js.map