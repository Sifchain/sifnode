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
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
const reactivity_1 = require("@vue/reactivity");
const B_1 = __importDefault(require("../entities/utils/B"));
const utils_1 = require("./utils");
exports.default = ({ api, store, }) => {
    api.EthereumService.onProviderNotFound(() => {
        api.EventBusService.dispatch({
            type: "WalletConnectionErrorEvent",
            payload: {
                walletType: "eth",
                message: "Metamask not found.",
            },
        });
    });
    api.EthereumService.onChainIdDetected((chainId) => {
        store.wallet.eth.chainId = chainId;
    });
    const etheriumState = api.EthereumService.getState();
    const actions = {
        isSupportedNetwork() {
            return utils_1.isSupportedEVMChain(store.wallet.eth.chainId);
        },
        disconnectWallet() {
            return __awaiter(this, void 0, void 0, function* () {
                yield api.EthereumService.disconnect();
            });
        },
        connectToWallet() {
            return __awaiter(this, void 0, void 0, function* () {
                try {
                    yield api.EthereumService.connect();
                }
                catch (err) {
                    api.EventBusService.dispatch({
                        type: "WalletConnectionErrorEvent",
                        payload: {
                            walletType: "eth",
                            message: "Failed to connect to Metamask.",
                        },
                    });
                }
            });
        },
        transferEthWallet(amount, recipient, asset) {
            return __awaiter(this, void 0, void 0, function* () {
                const hash = yield api.EthereumService.transfer({
                    amount: B_1.default(amount, asset.decimals),
                    recipient,
                    asset,
                });
                return hash;
            });
        },
    };
    reactivity_1.effect(() => {
        // Only show connected when we have an address
        if (store.wallet.eth.isConnected !== etheriumState.connected) {
            store.wallet.eth.isConnected =
                etheriumState.connected && !!etheriumState.address;
            if (store.wallet.eth.isConnected) {
                api.EventBusService.dispatch({
                    type: "WalletConnectedEvent",
                    payload: {
                        walletType: "eth",
                        address: store.wallet.eth.address,
                    },
                });
            }
        }
    });
    reactivity_1.effect(() => {
        store.wallet.eth.address = etheriumState.address;
    });
    reactivity_1.effect(() => {
        store.wallet.eth.balances = etheriumState.balances;
    });
    reactivity_1.effect(() => __awaiter(void 0, void 0, void 0, function* () {
        etheriumState.log; // trigger on log change
        yield api.EthereumService.getBalance();
    }));
    return actions;
};
//# sourceMappingURL=ethWallet.js.map