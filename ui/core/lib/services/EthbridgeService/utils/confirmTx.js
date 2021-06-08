"use strict";
// Some of this was from https://medium.com/pixelpoint/track-blockchain-transactions-like-a-boss-with-web3-js-c149045ca9bf
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
exports.confirmTx = exports.getConfirmations = void 0;
function getConfirmations(web3, txHash) {
    return __awaiter(this, void 0, void 0, function* () {
        try {
            // Get transaction details
            const trx = yield web3.eth.getTransaction(txHash);
            // Get current block number
            const currentBlock = yield web3.eth.getBlockNumber();
            // When transaction is unconfirmed, its block number is null.
            // In this case we return 0 as number of confirmations
            return trx.blockNumber === null ? 0 : currentBlock - trx.blockNumber;
        }
        catch (error) {
            console.log(error);
            return 0;
        }
    });
}
exports.getConfirmations = getConfirmations;
function confirmTx({ web3, txHash, confirmations = 10, onSuccess = () => { }, onCheckConfirmation = () => { }, }) {
    let currentCount = 0;
    setTimeout(() => __awaiter(this, void 0, void 0, function* () {
        const confirmationCount = yield getConfirmations(web3, txHash);
        if (currentCount !== confirmationCount) {
            onCheckConfirmation && onCheckConfirmation(confirmationCount);
        }
        currentCount = confirmationCount;
        if (currentCount >= confirmations) {
            onSuccess && onSuccess();
            return;
        }
        confirmTx({
            web3,
            txHash,
            confirmations,
            onSuccess,
            onCheckConfirmation,
        });
    }), 500);
}
exports.confirmTx = confirmTx;
//# sourceMappingURL=confirmTx.js.map