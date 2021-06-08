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
const launchpad_1 = require("@cosmjs/launchpad");
const reactivity_1 = require("@vue/reactivity");
const lodash_1 = require("lodash");
const entities_1 = require("../../entities");
const SifClient_1 = require("../utils/SifClient");
const utils_1 = require("./utils");
const getKeplrProvider_1 = __importDefault(require("./getKeplrProvider"));
const parseTxFailure_1 = require("./parseTxFailure");
/**
 * Constructor for SifService
 *
 * SifService handles communication between our ui core Domain and the SifNode blockchain
 */
function createSifService({ sifAddrPrefix, sifApiUrl, sifWsUrl, sifRpcUrl, keplrChainConfig, assets, }) {
    const {} = sifAddrPrefix;
    const state = reactivity_1.reactive({
        connected: false,
        accounts: [],
        address: "",
        balances: [],
        log: "unset",
    });
    const keplrProviderPromise = getKeplrProvider_1.default();
    let keplrProvider;
    let client = null;
    let polling;
    let connecting = false;
    const unSignedClient = new SifClient_1.SifUnSignedClient(sifApiUrl, sifWsUrl, sifRpcUrl);
    const supportedTokens = assets.filter((asset) => asset.network === entities_1.Network.SIFCHAIN);
    function createSifClientFromMnemonic(mnemonic) {
        return __awaiter(this, void 0, void 0, function* () {
            const wallet = yield launchpad_1.Secp256k1HdWallet.fromMnemonic(mnemonic, launchpad_1.makeCosmoshubPath(0), sifAddrPrefix);
            const accounts = yield wallet.getAccounts();
            const address = accounts.length > 0 ? accounts[0].address : "";
            if (!address) {
                throw new Error("No address on sif account");
            }
            return new SifClient_1.SifClient(sifApiUrl, address, wallet, sifWsUrl, sifRpcUrl);
        });
    }
    const triggerUpdate = lodash_1.debounce(() => __awaiter(this, void 0, void 0, function* () {
        try {
            if (!polling) {
                polling = setInterval(() => {
                    triggerUpdate();
                }, 2000);
            }
            yield instance.setClient();
            if (!client) {
                state.connected = false;
                state.address = "";
                state.balances = [];
                state.accounts = [];
                state.log = "";
                return;
            }
            state.connected = !!client;
            state.address = client.senderAddress;
            state.accounts = yield client.getAccounts();
            state.balances = yield instance.getBalance(client.senderAddress);
        }
        catch (e) {
            if (!e.toString().toLowerCase().includes("no address found on chain")) {
                state.connected = false;
                state.address = "";
                state.balances = [];
                state.accounts = [];
                state.log = "";
                if (polling) {
                    clearInterval(polling);
                    polling = null;
                }
            }
        }
    }), 100, { leading: true });
    const instance = {
        /**
         * getState returns the service's reactive state to be listened to by consuming clients.
         */
        getState() {
            return state;
        },
        getSupportedTokens() {
            return supportedTokens;
        },
        setClient() {
            return __awaiter(this, void 0, void 0, function* () {
                if (!keplrProvider) {
                    return;
                }
                if (connecting || state.connected) {
                    return;
                }
                connecting = true;
                /*
                  Only load dev env keplr configs.
                  Will need to change chain id in devnet, testnet so keplr asks to add experimental chain.
                  Otherwise, if sifchain, auto maps to production chain per keplr code.
                */
                if (!state.connected && keplrChainConfig.chainId !== "sifchain") {
                    yield this.connect();
                }
                const offlineSigner = keplrProvider.getOfflineSigner(keplrChainConfig.chainId);
                const accounts = yield offlineSigner.getAccounts();
                console.log("account", accounts);
                const address = accounts.length > 0 ? accounts[0].address : "";
                if (!address) {
                    throw "No address on sif account";
                }
                client = new SifClient_1.SifClient(sifApiUrl, address, offlineSigner, sifWsUrl, sifRpcUrl);
                connecting = false;
            });
        },
        initProvider() {
            return __awaiter(this, void 0, void 0, function* () {
                try {
                    keplrProvider = yield keplrProviderPromise;
                    if (!keplrProvider) {
                        return;
                    }
                    triggerUpdate();
                }
                catch (e) {
                    console.log("initProvider", e);
                }
            });
        },
        connect() {
            return __awaiter(this, void 0, void 0, function* () {
                if (!keplrProvider) {
                    keplrProvider = yield keplrProviderPromise;
                }
                // open extension
                if (keplrProvider.experimentalSuggestChain) {
                    try {
                        yield keplrProvider.experimentalSuggestChain(keplrChainConfig);
                        yield keplrProvider.enable(keplrChainConfig.chainId);
                        triggerUpdate();
                    }
                    catch (error) {
                        console.log(error);
                        throw { message: "Failed to Suggest Chain" };
                    }
                }
                else {
                    throw {
                        message: "Keplr Outdated",
                        detail: { type: "info", message: "Need at least 0.6.4" },
                    };
                }
            });
        },
        isConnected() {
            return state.connected;
        },
        onSocketError(handler) {
            unSignedClient.onSocketError(handler);
        },
        onTx(handler) {
            unSignedClient.onTx(handler);
        },
        onNewBlock(handler) {
            unSignedClient.onNewBlock(handler);
        },
        // Required solely for testing purposes
        setPhrase(mnemonic) {
            return __awaiter(this, void 0, void 0, function* () {
                try {
                    if (!mnemonic) {
                        throw "No mnemonic. Can't generate wallet.";
                    }
                    client = yield createSifClientFromMnemonic(mnemonic);
                    return client.senderAddress;
                }
                catch (error) {
                    throw error;
                }
            });
        },
        purgeClient() {
            return __awaiter(this, void 0, void 0, function* () {
                // We currently delegate auth to Keplr so this is irrelevant
            });
        },
        getBalance(address, asset) {
            return __awaiter(this, void 0, void 0, function* () {
                if (!client) {
                    throw "No client. Please sign in.";
                }
                if (!address) {
                    throw "Address undefined. Fail";
                }
                utils_1.ensureSifAddress(address);
                try {
                    const account = yield client.getAccount(address);
                    if (!account) {
                        throw "No Address found on chain";
                    } // todo handle this better
                    const supportedTokenSymbols = supportedTokens.map((s) => s.symbol);
                    return account.balance
                        .filter((balance) => supportedTokenSymbols.includes(balance.denom))
                        .map(({ amount, denom }) => {
                        const asset = supportedTokens.find((token) => token.symbol === denom); // will be found because of filter above
                        return entities_1.AssetAmount(asset, amount);
                    })
                        .filter((balance) => {
                        // If an aseet is supplied filter for it
                        if (!asset) {
                            return true;
                        }
                        return balance.asset.symbol === asset.symbol;
                    });
                }
                catch (error) {
                    throw error;
                }
            });
        },
        transfer(params) {
            return __awaiter(this, void 0, void 0, function* () {
                if (!client) {
                    throw "No client. Please sign in.";
                }
                if (!params.asset) {
                    throw "No asset.";
                }
                try {
                    const msg = {
                        type: "cosmos-sdk/MsgSend",
                        value: {
                            amount: [
                                {
                                    amount: params.amount.toString(),
                                    denom: params.asset.symbol,
                                },
                            ],
                            from_address: client.senderAddress,
                            to_address: params.recipient,
                        },
                    };
                    const fee = {
                        amount: launchpad_1.coins(250000, params.asset.symbol),
                        gas: "500000",
                    };
                    return yield client.signAndBroadcast([msg], fee, params.memo);
                }
                catch (err) {
                    console.error(err);
                }
            });
        },
        signAndBroadcast(msg, memo) {
            return __awaiter(this, void 0, void 0, function* () {
                if (!client) {
                    throw "No client. Please sign in.";
                }
                try {
                    const fee = {
                        // Keplr overwrites this in app but for unit/integration tests where we
                        // dont connect to keplr we need to specify an amount of rowan to pay for the fee.
                        amount: launchpad_1.coins(250000, "rowan"),
                        gas: "500000",
                    };
                    const msgArr = Array.isArray(msg) ? msg : [msg];
                    const result = yield client.signAndBroadcast(msgArr, fee, memo);
                    if (launchpad_1.isBroadcastTxFailure(result)) {
                        /* istanbul ignore next */ // TODO: fix coverage
                        return parseTxFailure_1.parseTxFailure(result);
                    }
                    return {
                        hash: result.transactionHash,
                        memo,
                        state: "accepted",
                    };
                }
                catch (err) {
                    console.log("signAndBroadcast ERROR", err);
                    return parseTxFailure_1.parseTxFailure({ transactionHash: "", rawLog: err.message });
                }
            });
        },
    };
    instance.initProvider();
    return instance;
}
exports.default = createSifService;
//# sourceMappingURL=SifService.js.map