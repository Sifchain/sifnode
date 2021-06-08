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
const utils_1 = require("../utils");
const subscribeToUnconfirmedPegTxs_1 = require("./subscribeToUnconfirmedPegTxs");
const subscribeToTx_1 = require("./utils/subscribeToTx");
function isOriginallySifchainNativeToken(asset) {
    return ["erowan", "rowan"].includes(asset.symbol);
}
exports.default = ({ services, store, }) => {
    const config = {
        // listen for 50 confirmations
        // Eventually this should be set on ebrelayer
        // to centralize the business logic
        ethConfirmations: 50,
    };
    // Create the context for passing to commands, queries and subscriptions
    const ctx = { services, store, config };
    /*
      TODO: suggestion externalize all interactors injecting ctx would look like the following
  
      const commands = {
        unpeg: Unpeg(ctx),
        peg: Peg(ctx),
      }
  
      const queries = {
        getSifTokens: GetSifTokens(ctx),
        getEthTokens: GetEthTokens(ctx),
        calculateUnpegFee: CalculateUnpegFee(ctx),
      }
  
      const subscriptions = {
        subscribeToUnconfirmedPegTxs: SubscribeToUnconfirmedPegTxs(ctx),
      }
    */
    // Rename and split this up to subscriptions, commands, queries
    const actions = {
        subscribeToUnconfirmedPegTxs: subscribeToUnconfirmedPegTxs_1.SubscribeToUnconfirmedPegTxs(ctx),
        getSifTokens() {
            return services.sif.getSupportedTokens();
        },
        getEthTokens() {
            return services.eth.getSupportedTokens();
        },
        calculateUnpegFee(asset) {
            const feeNumber = isOriginallySifchainNativeToken(asset)
                ? "70000000000000000"
                : "70000000000000000";
            return entities_1.AssetAmount(entities_1.Asset.get("ceth"), feeNumber);
        },
        unpeg(assetAmount) {
            return __awaiter(this, void 0, void 0, function* () {
                const lockOrBurnFn = isOriginallySifchainNativeToken(assetAmount.asset)
                    ? services.ethbridge.lockToEthereum
                    : services.ethbridge.burnToEthereum;
                const feeAmount = this.calculateUnpegFee(assetAmount.asset);
                const tx = yield lockOrBurnFn({
                    assetAmount,
                    ethereumRecipient: store.wallet.eth.address,
                    fromAddress: store.wallet.sif.address,
                    feeAmount,
                });
                const txStatus = yield services.sif.signAndBroadcast(tx.value.msg);
                if (txStatus.state !== "accepted") {
                    services.bus.dispatch({
                        type: "PegTransactionErrorEvent",
                        payload: {
                            txStatus,
                            message: txStatus.memo || "There was an error while unpegging",
                        },
                    });
                }
                console.log("unpeg txStatus.state", txStatus.state, txStatus.memo, txStatus.code, tx.value.msg);
                return txStatus;
            });
        },
        // TODO: Move this approval command to within peg and report status via store or some other means
        //       This has been done for convenience but we should not have to know in the view that
        //       approval is required before pegging as that is very much business domain knowledge
        approve(address, assetAmount) {
            return __awaiter(this, void 0, void 0, function* () {
                return yield services.ethbridge.approveBridgeBankSpend(address, assetAmount);
            });
        },
        peg(assetAmount) {
            return __awaiter(this, void 0, void 0, function* () {
                if (assetAmount.asset.network === entities_1.Network.ETHEREUM &&
                    !utils_1.isSupportedEVMChain(store.wallet.eth.chainId)) {
                    services.bus.dispatch({
                        type: "ErrorEvent",
                        payload: {
                            message: "EVM Network not supported!",
                        },
                    });
                    return {
                        hash: "",
                        state: "failed",
                    };
                }
                const subscribeToTx = subscribeToTx_1.SubscribeToTx(ctx);
                const lockOrBurnFn = isOriginallySifchainNativeToken(assetAmount.asset)
                    ? services.ethbridge.burnToSifchain
                    : services.ethbridge.lockToSifchain;
                return yield new Promise((done) => {
                    const pegTx = lockOrBurnFn(store.wallet.sif.address, assetAmount, config.ethConfirmations);
                    subscribeToTx(pegTx);
                    pegTx.onTxHash((hash) => {
                        done({
                            hash: hash.txHash,
                            memo: "Transaction Accepted",
                            state: "accepted",
                        });
                    });
                });
            });
        },
    };
    return actions;
};
//# sourceMappingURL=index.js.map