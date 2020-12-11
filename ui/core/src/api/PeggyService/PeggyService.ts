import { AssetAmount } from "../../entities";
import { IPeggyService } from "./IPeggyService";
import { createTxEventEmitter } from "./TxEventEmitter";
import { TxEventEmitter } from "./types";

// MOCK SEQUENCES
const sleep = (ms: number) => new Promise((resolve) => setTimeout(resolve, ms));
function mockLockSequence(emitter: TxEventEmitter) {
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

function mockBurnSequence(emitter: TxEventEmitter) {
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

export default function createPeggyService(): IPeggyService {
  return {
    burn(ethereumRecipient: string, assetAmount: AssetAmount) {
      // Some random string for now
      const txHash = "abcd1234";
      // Create an emitter
      const e = createTxEventEmitter(txHash);
      // Direct that emitter through a mock sequence
      return mockBurnSequence(e);
    },
    lock(cosmosRecipient: string, assetAmount: AssetAmount) {
      // Some random string for now
      const txHash = "abcd1234";
      // Create an emitter
      const e = createTxEventEmitter(txHash);
      // Direct that emitter through a mock sequence
      return mockLockSequence(e);
    },
  };
}
