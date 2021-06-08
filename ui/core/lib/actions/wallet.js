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
const bip39_1 = require("bip39");
const reactivity_1 = require("@vue/reactivity");
exports.default = ({ api, store, }) => {
    const state = api.SifService.getState();
    const actions = {
        getCosmosBalances(address) {
            return __awaiter(this, void 0, void 0, function* () {
                // TODO: validate sif prefix
                return yield api.SifService.getBalance(address);
            });
        },
        connect(mnemonic) {
            return __awaiter(this, void 0, void 0, function* () {
                if (!mnemonic)
                    throw "Mnemonic must be defined";
                if (!bip39_1.validateMnemonic(mnemonic))
                    throw "Invalid Mnemonic. Not sent.";
                return yield api.SifService.setPhrase(mnemonic);
            });
        },
        sendCosmosTransaction(params) {
            return __awaiter(this, void 0, void 0, function* () {
                return yield api.SifService.transfer(params);
            });
        },
        disconnect() {
            return __awaiter(this, void 0, void 0, function* () {
                api.SifService.purgeClient();
            });
        },
        connectToWallet() {
            return __awaiter(this, void 0, void 0, function* () {
                try {
                    // TODO type
                    yield api.SifService.connect();
                    store.wallet.sif.isConnected = true;
                }
                catch (error) {
                    api.EventBusService.dispatch({
                        type: "WalletConnectionErrorEvent",
                        payload: {
                            walletType: "sif",
                            message: "Failed to connect to Keplr.",
                        },
                    });
                }
            });
        },
    };
    reactivity_1.effect(() => {
        if (store.wallet.sif.isConnected !== state.connected) {
            store.wallet.sif.isConnected = state.connected;
            if (store.wallet.sif.isConnected) {
                api.EventBusService.dispatch({
                    type: "WalletConnectedEvent",
                    payload: {
                        walletType: "sif",
                        address: store.wallet.sif.address,
                    },
                });
            }
        }
    });
    reactivity_1.effect(() => {
        store.wallet.sif.address = state.address;
    });
    reactivity_1.effect(() => {
        store.wallet.sif.balances = state.balances;
    });
    return actions;
};
//# sourceMappingURL=wallet.js.map