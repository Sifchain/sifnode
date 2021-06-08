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
exports.getEtheriumBalance = exports.transferEther = exports.transferToken = exports.transferAsset = exports.isMetaMaskInpageProvider = exports.isEventEmittingProvider = exports.getTokenBalance = exports.getTokenContract = void 0;
const entities_1 = require("../../../entities");
const erc20TokenAbi_1 = __importDefault(require("./erc20TokenAbi"));
function getTokenContract(web3, asset) {
    return new web3.eth.Contract(erc20TokenAbi_1.default, asset.address);
}
exports.getTokenContract = getTokenContract;
function getTokenBalance(web3, address, asset) {
    return __awaiter(this, void 0, void 0, function* () {
        const contract = getTokenContract(web3, asset);
        let tokenBalance = "0";
        try {
            tokenBalance = yield contract.methods.balanceOf(address).call();
        }
        catch (err) {
            console.log(`Error fetching balance for ${asset.symbol}`);
        }
        return entities_1.AssetAmount(asset, tokenBalance);
    });
}
exports.getTokenBalance = getTokenBalance;
function isEventEmittingProvider(provider) {
    if (!provider || typeof provider === "string")
        return false;
    return typeof provider.on === "function";
}
exports.isEventEmittingProvider = isEventEmittingProvider;
function isMetaMaskInpageProvider(provider) {
    if (!provider || typeof provider === "string")
        return false;
    return typeof provider.request === "function";
}
exports.isMetaMaskInpageProvider = isMetaMaskInpageProvider;
// Transfer token or ether
function transferAsset(web3, fromAddress, toAddress, amount, asset) {
    return __awaiter(this, void 0, void 0, function* () {
        if (asset === null || asset === void 0 ? void 0 : asset.address) {
            return yield transferToken(web3, fromAddress, toAddress, amount, asset);
        }
        return yield transferEther(web3, fromAddress, toAddress, amount);
    });
}
exports.transferAsset = transferAsset;
// Transfer token
function transferToken(web3, fromAddress, toAddress, amount, asset) {
    return __awaiter(this, void 0, void 0, function* () {
        const contract = getTokenContract(web3, asset);
        return new Promise((resolve, reject) => {
            let hash;
            let receipt;
            function resolvePromise() {
                if (receipt && hash)
                    resolve(hash);
            }
            contract.methods
                .transfer(toAddress, amount.toString())
                .send({ from: fromAddress })
                .on("transactionHash", (_hash) => {
                hash = _hash;
                resolvePromise();
            })
                .on("receipt", (_receipt) => {
                receipt = _receipt;
                resolvePromise();
            })
                .on("error", (err) => {
                reject(err);
            });
        });
    });
}
exports.transferToken = transferToken;
// Transfer ether
function transferEther(web3, fromAddress, toAddress, amount) {
    return __awaiter(this, void 0, void 0, function* () {
        return new Promise((resolve, reject) => {
            let hash;
            let receipt;
            function resolvePromise() {
                if (receipt && hash)
                    resolve(hash);
            }
            web3.eth
                .sendTransaction({
                from: fromAddress,
                to: toAddress,
                value: amount.toString(),
            })
                .on("transactionHash", (_hash) => {
                hash = _hash;
                resolvePromise();
            })
                .on("receipt", (_receipt) => {
                receipt = _receipt;
                resolvePromise();
            })
                .on("error", (err) => {
                reject(err);
            });
        });
    });
}
exports.transferEther = transferEther;
function getEtheriumBalance(web3, address) {
    return __awaiter(this, void 0, void 0, function* () {
        const ethBalance = yield web3.eth.getBalance(address);
        // TODO: Pull as search from supported tokens
        return entities_1.AssetAmount({
            symbol: "eth",
            label: "ETH",
            address: "",
            decimals: 18,
            name: "Ethereum",
            network: entities_1.Network.ETHEREUM,
        }, ethBalance);
    });
}
exports.getEtheriumBalance = getEtheriumBalance;
//# sourceMappingURL=ethereumUtils.js.map