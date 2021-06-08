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
const web3_1 = __importDefault(require("web3"));
const bridgebankContract_1 = require("./bridgebankContract");
const tokenContract_1 = require("./tokenContract");
const PegTxEventEmitter_1 = require("./PegTxEventEmitter");
const confirmTx_1 = require("./utils/confirmTx");
const SifClient_1 = require("../utils/SifClient");
const parseTxFailure_1 = require("./parseTxFailure");
const jsbi_1 = __importDefault(require("jsbi"));
const ETH_ADDRESS = "0x0000000000000000000000000000000000000000";
function createEthbridgeService({ sifApiUrl, sifWsUrl, sifRpcUrl, sifChainId, bridgebankContractAddress, getWeb3Provider, sifUnsignedClient = new SifClient_1.SifUnSignedClient(sifApiUrl, sifWsUrl, sifRpcUrl), }) {
    // Pull this out to a util?
    // How to handle context/dependency injection?
    let _web3 = null;
    function ensureWeb3() {
        return __awaiter(this, void 0, void 0, function* () {
            if (!_web3) {
                _web3 = new web3_1.default(yield getWeb3Provider());
            }
            return _web3;
        });
    }
    /**
     * Create an event listener to report status of a peg transaction.
     * Usage:
     * const tx = createPegTx(50)
     * tx.setTxHash('0x52ds.....'); // set the hash to lookup and confirm on the blockchain
     * @param confirmations number of confirmations before pegtx is considered confirmed
     */
    function createPegTx(confirmations, symbol, txHash) {
        const emitter = PegTxEventEmitter_1.createPegTxEventEmitter(txHash, symbol);
        // decorate pegtx to invert dependency to web3 and confirmations
        emitter.onTxHash(({ payload: txHash }) => __awaiter(this, void 0, void 0, function* () {
            const web3 = yield ensureWeb3();
            confirmTx_1.confirmTx({
                web3,
                txHash,
                confirmations,
                onSuccess() {
                    emitter.emit({ type: "Complete", payload: null });
                },
                onCheckConfirmation(count) {
                    emitter.emit({ type: "EthConfCountChanged", payload: count });
                },
            });
        }));
        return emitter;
    }
    /**
     * Gets a list of transactionHashes found as _from keys within the given events within a given blockRange from the current block
     * @param {*} address eth address to correlate transactions with
     * @param {*} contract web3 contract
     * @param {*} eventList event name list of events (must have an addresskey)
     * @param {*} blockRange number of blocks from the current block header to search
     */
    function getEventTxsInBlockrangeFromAddress(address, contract, eventList, blockRange) {
        var _a, _b, _c;
        return __awaiter(this, void 0, void 0, function* () {
            const web3 = yield ensureWeb3();
            const latest = yield web3.eth.getBlockNumber();
            const fromBlock = Math.max(latest - blockRange, 0);
            const allEvents = yield contract.getPastEvents("allEvents", {
                // filter:{_from:address}, // if _from was indexed we could do this
                fromBlock,
                toBlock: "latest",
            });
            // unfortunately because _from is not an indexed key we have to manually filter
            // TODO: ask peggy team to index the _from field which would make this more efficient
            const txs = [];
            for (let event of allEvents) {
                const isEventWeCareAbout = eventList.includes(event.event);
                const matchesInputAddress = address &&
                    ((_b = (_a = event === null || event === void 0 ? void 0 : event.returnValues) === null || _a === void 0 ? void 0 : _a._from) === null || _b === void 0 ? void 0 : _b.toLowerCase()) === address.toLowerCase();
                if (isEventWeCareAbout && matchesInputAddress && event.transactionHash) {
                    txs.push({
                        hash: event.transactionHash,
                        symbol: (_c = event.returnValues) === null || _c === void 0 ? void 0 : _c._symbol,
                    });
                }
            }
            return txs;
        });
    }
    return {
        approveBridgeBankSpend(account, amount) {
            return __awaiter(this, void 0, void 0, function* () {
                // This will popup an approval request in metamask
                const web3 = yield ensureWeb3();
                const tokenContract = yield tokenContract_1.getTokenContract(web3, amount.asset.address);
                const sendArgs = {
                    from: account,
                    value: 0,
                    gas: 100000,
                };
                // TODO - give interface option to approve unlimited spend via web3.utils.toTwosComplement(-1);
                // NOTE - We may want to move this out into its own separate function.
                // Although I couldn't think of a situation we'd call allowance separately from approve
                const hasAlreadyApprovedSpend = yield tokenContract.methods
                    .allowance(account, bridgebankContractAddress)
                    .call();
                if (jsbi_1.default.lessThanOrEqual(amount.toBigInt(), jsbi_1.default.BigInt(hasAlreadyApprovedSpend))) {
                    // dont request approve again
                    console.log("approveBridgeBankSpend: spend already approved", hasAlreadyApprovedSpend);
                    return;
                }
                const res = yield tokenContract.methods
                    .approve(bridgebankContractAddress, amount.toBigInt().toString())
                    .send(sendArgs);
                console.log("approveBridgeBankSpend:", res);
                return res;
            });
        },
        burnToEthereum(params) {
            var _a;
            return __awaiter(this, void 0, void 0, function* () {
                const web3 = yield ensureWeb3();
                const ethereumChainId = yield web3.eth.net.getId();
                const tokenAddress = (_a = params.assetAmount.asset.address) !== null && _a !== void 0 ? _a : ETH_ADDRESS;
                console.log("burnToEthereum: start: ", tokenAddress);
                const txReceipt = yield sifUnsignedClient.burn({
                    ethereum_receiver: params.ethereumRecipient,
                    base_req: {
                        chain_id: sifChainId,
                        from: params.fromAddress,
                    },
                    amount: params.assetAmount.toBigInt().toString(),
                    symbol: params.assetAmount.asset.symbol,
                    cosmos_sender: params.fromAddress,
                    ethereum_chain_id: `${ethereumChainId}`,
                    token_contract_address: tokenAddress,
                    ceth_amount: params.feeAmount.toBigInt().toString(),
                });
                console.log("burnToEthereum: txReceipt: ", txReceipt, tokenAddress);
                return txReceipt;
            });
        },
        lockToSifchain(sifRecipient, assetAmount, confirmations) {
            const pegTx = createPegTx(confirmations, assetAmount.asset.symbol);
            function handleError(err) {
                console.log("lockToSifchain: handleError: ", err);
                pegTx.emit({
                    type: "Error",
                    payload: parseTxFailure_1.parseTxFailure({ hash: "", log: err.message.toString() }),
                });
            }
            (function () {
                return __awaiter(this, void 0, void 0, function* () {
                    const web3 = yield ensureWeb3();
                    const cosmosRecipient = web3_1.default.utils.utf8ToHex(sifRecipient);
                    const bridgeBankContract = yield bridgebankContract_1.getBridgeBankContract(web3, bridgebankContractAddress);
                    const accounts = yield web3.eth.getAccounts();
                    const coinDenom = assetAmount.asset.address || ETH_ADDRESS; // eth address is ""
                    const amount = assetAmount.toBigInt().toString();
                    const fromAddress = accounts[0];
                    const sendArgs = {
                        from: fromAddress,
                        value: coinDenom === ETH_ADDRESS ? amount : 0,
                        gas: 150000,
                    };
                    console.log("lockToSifchain: bridgeBankContract.lock", JSON.stringify({ cosmosRecipient, coinDenom, amount, sendArgs }));
                    bridgeBankContract.methods
                        .lock(cosmosRecipient, coinDenom, amount)
                        .send(sendArgs)
                        .on("transactionHash", (hash) => {
                        console.log("lockToSifchain: bridgeBankContract.lock TX", hash);
                        pegTx.setTxHash(hash);
                    })
                        .on("error", (err) => {
                        console.log("lockToSifchain: bridgeBankContract.lock ERROR", err);
                        handleError(err);
                    });
                });
            })().catch((err) => {
                handleError(err);
            });
            return pegTx;
        },
        lockToEthereum(params) {
            var _a;
            return __awaiter(this, void 0, void 0, function* () {
                const web3 = yield ensureWeb3();
                const ethereumChainId = yield web3.eth.net.getId();
                const tokenAddress = (_a = params.assetAmount.asset.address) !== null && _a !== void 0 ? _a : ETH_ADDRESS;
                const lockParams = {
                    ethereum_receiver: params.ethereumRecipient,
                    base_req: {
                        chain_id: sifChainId,
                        from: params.fromAddress,
                    },
                    amount: params.assetAmount.toBigInt().toString(),
                    symbol: params.assetAmount.asset.symbol,
                    cosmos_sender: params.fromAddress,
                    ethereum_chain_id: `${ethereumChainId}`,
                    token_contract_address: tokenAddress,
                    ceth_amount: params.feeAmount.toBigInt().toString(),
                };
                console.log("lockToEthereum: TRY LOCK", tokenAddress);
                const lockReceipt = yield sifUnsignedClient.lock(lockParams);
                console.log("lockToEthereum: LOCKED", lockReceipt);
                return lockReceipt;
            });
        },
        /**
         * Get a list of unconfirmed transaction hashes associated with
         * a particular address and return pegTxs associated with that hash
         * @param address contract address
         * @param confirmations number of confirmations required
         */
        fetchUnconfirmedLockBurnTxs(address, confirmations) {
            return __awaiter(this, void 0, void 0, function* () {
                const web3 = yield ensureWeb3();
                const bridgeBankContract = yield bridgebankContract_1.getBridgeBankContract(web3, bridgebankContractAddress);
                const txs = yield getEventTxsInBlockrangeFromAddress(address, bridgeBankContract, ["LogBurn", "LogLock"], confirmations);
                return txs.map(({ hash, symbol }) => createPegTx(confirmations, symbol, hash));
            });
        },
        burnToSifchain(sifRecipient, assetAmount, confirmations, account) {
            const pegTx = createPegTx(confirmations, assetAmount.asset.symbol);
            function handleError(err) {
                console.log("burnToSifchain: handleError ERROR", err);
                pegTx.emit({
                    type: "Error",
                    payload: parseTxFailure_1.parseTxFailure({ hash: "", log: err }),
                });
            }
            (function () {
                return __awaiter(this, void 0, void 0, function* () {
                    const web3 = yield ensureWeb3();
                    const cosmosRecipient = web3_1.default.utils.utf8ToHex(sifRecipient);
                    const bridgeBankContract = yield bridgebankContract_1.getBridgeBankContract(web3, bridgebankContractAddress);
                    const accounts = yield web3.eth.getAccounts();
                    const coinDenom = assetAmount.asset.address;
                    const amount = assetAmount.toBigInt().toString();
                    const fromAddress = account || accounts[0];
                    const sendArgs = {
                        from: fromAddress,
                        value: 0,
                        gas: 150000,
                    };
                    bridgeBankContract.methods
                        .burn(cosmosRecipient, coinDenom, amount)
                        .send(sendArgs)
                        .on("transactionHash", (hash) => {
                        console.log("burnToSifchain: bridgeBankContract.burn TX", hash);
                        pegTx.setTxHash(hash);
                    })
                        .on("error", (err) => {
                        console.log("burnToSifchain: bridgeBankContract.burn ERROR", err);
                        handleError(err);
                    });
                });
            })().catch((err) => {
                handleError(err);
            });
            return pegTx;
        },
    };
}
exports.default = createEthbridgeService;
//# sourceMappingURL=EthbridgeService.js.map