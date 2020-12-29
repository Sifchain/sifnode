import { ActionContext } from "..";
import { createPegTxEventEmitter } from "./PegTxEventEmitter";
import { AssetAmount } from "../../entities";
import { PegTxEventEmitter } from "./types";

// MOCK SEQUENCES
const sleep = (ms: number) => new Promise((resolve) => setTimeout(resolve, ms));
function mockLockSequence(emitter: PegTxEventEmitter) {
  (async () => {
    await sleep(20);
    emitter.emit({ type: "EthTxInitiated", payload: {} });
    for (let count of Array.from(Array(30).keys())) {
      await sleep(10);
      emitter.emit({ type: "EthConfCountChanged", payload: count });
    }
    emitter.emit({ type: "EthTxConfirmed", payload: {} });
    await sleep(20);
    emitter.emit({ type: "SifTxInitiated", payload: {} });
    await sleep(50);
    emitter.emit({ type: "SifTxConfirmed", payload: {} });
    emitter.emit({ type: "Complete", payload: {} });
  })();
  return emitter;
}

function mockBurnSequence(emitter: PegTxEventEmitter) {
  (async () => {
    await sleep(20);
    emitter.emit({ type: "SifTxInitiated", payload: {} });
    for (let count of Array.from(Array(30).keys())) {
      await sleep(10);
      emitter.emit({ type: "SifConfCountChanged", payload: count });
    }
    emitter.emit({ type: "SifTxConfirmed", payload: {} });
    await sleep(20);
    emitter.emit({ type: "EthTxInitiated", payload: {} });
    await sleep(50);
    emitter.emit({ type: "EthTxConfirmed", payload: {} });
    emitter.emit({ type: "Complete", payload: {} });
  })();
  return emitter;
}
// end MOCK

export default ({
  api,
}: ActionContext<"SifService" | "EthereumService", "asset">) => {
  const actions = {
    getSifTokens() {
      return api.SifService.getSupportedTokens();
    },
    getEthTokens() {
      return api.EthereumService.getSupportedTokens();
    },
    burn(ethereumRecipient: string, assetAmount: AssetAmount) {
      // Some random string for now
      const txHash = "abcd1234";
      // Create an emitter
      const e = createPegTxEventEmitter(txHash);
      // Direct that emitter through a mock sequence
      return mockBurnSequence(e);
    },
    lock(cosmosRecipient: string, assetAmount: AssetAmount) {
      // 1. send tx to ethereum contract
      //
      //
      // // Some random string for now
      // const txHash = "abcd1234";
      // // Create an emitter
      // const e = createPegTxEventEmitter(txHash);
      // // add chaos
      // if (assetAmount.equalTo("100")) {
      //   setTimeout(() => {
      //     e.emit({ txHash, type: "Error", payload: "Boom!" });
      //   }, 345);
      // }
      // // Direct that emitter through a mock sequence
      // return mockLockSequence(e);
    },
  };

  return actions;
};
