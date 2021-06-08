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
exports.EthereumService = void 0;
const reactivity_1 = require("@vue/reactivity");
const web3_1 = __importDefault(require("web3"));
const lodash_1 = require("lodash");
const entities_1 = require("../../entities");
const ethereumUtils_1 = require("./utils/ethereumUtils");
const initState = {
    connected: false,
    accounts: [],
    balances: [],
    address: "",
    log: "unset",
};
// TODO: Refactor to be Module pattern with constructor function ie. `EthereumService()`
class EthereumService {
    constructor(getWeb3Provider, assets) {
        this.web3 = null;
        this.supportedTokens = [];
        this.reportProviderNotFound = () => { };
        this.chainIdDetectedHandler = (_chainId) => { };
        this.updateData = lodash_1.debounce(() => __awaiter(this, void 0, void 0, function* () {
            var _a;
            if (!this.web3) {
                this.state.connected = false;
                this.state.accounts = [];
                this.state.address = "";
                this.state.balances = [];
                return;
            }
            this.state.connected = true;
            this.state.accounts = (_a = (yield this.web3.eth.getAccounts())) !== null && _a !== void 0 ? _a : [];
            this.state.address = this.state.accounts[0];
            this.state.balances = yield this.getBalance();
        }), 100, { leading: true });
        this.state = reactivity_1.reactive(Object.assign({}, initState));
        this.supportedTokens = assets.filter((t) => t.network === entities_1.Network.ETHEREUM);
        this.providerPromise = getWeb3Provider();
        this.providerPromise
            .then((provider) => {
            // Provider not found
            if (!provider) {
                this.provider = null;
                this.reportProviderNotFound();
                return;
            }
            if (ethereumUtils_1.isEventEmittingProvider(provider)) {
                provider.on("chainChanged", () => window.location.reload());
                provider.on("accountsChanged", () => this.updateData());
            }
            // What network is the provider on
            if (ethereumUtils_1.isMetaMaskInpageProvider(provider)) {
                provider.request({ method: "eth_chainId" }).then((chainId) => {
                    this.chainIdDetectedHandler(chainId);
                });
            }
            this.web3 = new web3_1.default(provider);
            this.provider = provider;
            this.addWeb3Subscription();
            this.updateData();
        })
            .catch((error) => {
            console.log("error", error);
        });
    }
    onChainIdDetected(handler) {
        this.chainIdDetectedHandler = handler;
    }
    onProviderNotFound(handler) {
        this.reportProviderNotFound = handler;
    }
    getState() {
        return this.state;
    }
    getAddress() {
        return this.state.address;
    }
    isConnected() {
        return this.state.connected;
    }
    getSupportedTokens() {
        return this.supportedTokens;
    }
    connect() {
        return __awaiter(this, void 0, void 0, function* () {
            const provider = yield this.providerPromise;
            try {
                if (!provider) {
                    throw new Error("Cannot connect because provider is not yet loaded!");
                }
                this.web3 = new web3_1.default(provider);
                if (ethereumUtils_1.isMetaMaskInpageProvider(provider)) {
                    if (provider.request) {
                        yield provider.request({ method: "eth_requestAccounts" });
                    }
                }
                this.addWeb3Subscription();
                yield this.updateData();
            }
            catch (err) {
                this.web3 = null;
                this.removeWeb3Subscription();
                throw err;
            }
        });
    }
    addWeb3Subscription() {
        var _a;
        this.blockSubscription = (_a = this.web3) === null || _a === void 0 ? void 0 : _a.eth.subscribe("newBlockHeaders", (error, blockHeader) => {
            if (blockHeader) {
                this.updateData();
                this.state.log = blockHeader.hash;
            }
            else {
                this.state.log = error.message;
            }
        });
    }
    removeWeb3Subscription() {
        var _a, _b;
        const success = (_a = this.blockSubscription) === null || _a === void 0 ? void 0 : _a.unsubscribe();
        if (success) {
            this.blockSubscription = null;
        }
        else {
            // try again if not success
            (_b = this.blockSubscription) === null || _b === void 0 ? void 0 : _b.unsubscribe();
        }
    }
    disconnect() {
        return __awaiter(this, void 0, void 0, function* () {
            if (ethereumUtils_1.isMetaMaskInpageProvider(this.provider)) {
                this.provider.disconnect &&
                    this.provider.disconnect(0, "Website disconnected wallet");
            }
            this.removeWeb3Subscription();
            this.web3 = null;
            yield this.updateData();
        });
    }
    getBalance(address, asset) {
        return __awaiter(this, void 0, void 0, function* () {
            const supportedTokens = this.getSupportedTokens();
            const addr = address || this.state.address;
            if (!this.web3 || !addr) {
                return [];
            }
            const web3 = this.web3;
            let balances = [];
            if (asset) {
                if (!asset.address) {
                    // Asset must be eth
                    const ethBalance = yield ethereumUtils_1.getEtheriumBalance(web3, addr);
                    balances = [ethBalance];
                }
                else {
                    // Asset must be ERC-20
                    const tokenBalance = yield ethereumUtils_1.getTokenBalance(web3, addr, asset);
                    balances = [tokenBalance];
                }
            }
            else {
                // No address no asset get everything
                balances = yield Promise.all([
                    ethereumUtils_1.getEtheriumBalance(web3, addr),
                    ...supportedTokens
                        .filter((t) => t.symbol !== "eth")
                        .map((token) => {
                        if (token.address)
                            return ethereumUtils_1.getTokenBalance(web3, addr, token);
                        return entities_1.AssetAmount(token, "0");
                    }),
                ]);
            }
            return balances;
        });
    }
    transfer(params) {
        return __awaiter(this, void 0, void 0, function* () {
            // TODO: validate params!!
            if (!this.web3) {
                throw new Error("Cannot do transfer because there is not yet a connection to Ethereum.");
            }
            const { amount, recipient, asset } = params;
            const from = this.getAddress();
            if (!from) {
                throw new Error("Transaction attempted but 'from' address cannot be determined!");
            }
            return yield ethereumUtils_1.transferAsset(this.web3, from, recipient, amount, asset);
        });
    }
    signAndBroadcast(msg, mmo) {
        return __awaiter(this, void 0, void 0, function* () { });
    }
    setPhrase(args) {
        return __awaiter(this, void 0, void 0, function* () {
            // We currently delegate auth to metamask so this is irrelavent
            return "";
        });
    }
    purgeClient() {
        // We currently delegate auth to metamask so this is irrelavent
    }
    static create({ getWeb3Provider, assets, }) {
        return new EthereumService(getWeb3Provider, assets);
    }
}
exports.EthereumService = EthereumService;
exports.default = EthereumService.create;
//# sourceMappingURL=EthereumService.js.map