import { OfflineSigner } from "@cosmjs/launchpad";
import { sleep } from "../../test/utils/sleep";

type WindowWithPossibleKeplr = typeof window & {
  keplr?: any;
  getOfflineSigner?: any;
};

// Mock out Keplr roughly. TODO import types
type Keplr = {
  experimentalSuggestChain?(chainInfo: any): Promise<any>;

  enable(chainId: string): Promise<void>;

  getKey(chainId: string): Promise<any>;

  getTxConfig(chainId: string, config: any): Promise<any>;

  sign(chainId: string, signer: string, message: Uint8Array): Promise<any>;

  sendTx(chainId: string, stdTx: any, mode: any): Promise<any>;

  suggestToken(chainId: string, contractAddress: string): Promise<void>;

  requestTx(
    chainId: string,
    txBytes: Uint8Array,
    mode: "sync" | "async" | "commit",
    isRestAPI: boolean
  ): Promise<void>;

  requestTxWithResult(
    chainId: string,
    txBytes: Uint8Array,
    mode: "sync" | "async" | "commit",
    isRestAPI: boolean
  ): Promise<any>;

  getSecret20ViewingKey(
    chainId: string,
    contractAddress: string
  ): Promise<string>;

  getOfflineSigner: (chainId?: string) => OfflineSigner;
};

// Todo
type provider = Keplr;

let numChecks = 0;

// Detect mossible keplr provider from browser
export default async function getKeplrProvider(): Promise<provider | null> {
  const win = window as WindowWithPossibleKeplr;

  if (!win) return null;
  console.log({
    "win.keplr": win.keplr,
    "win.getOfflineSigner": win.getOfflineSigner,
  });

  if (!win.keplr || !win.getOfflineSigner) {
    numChecks++;
    if (numChecks > 20) {
      return null;
    }
    await sleep(100);
    return getKeplrProvider();
  }

  // assign offline signer (they use __proto__ for some reason), so this is not as pretty as i'd like)
  Object.getPrototypeOf(win.keplr).getOfflineSigner = win.getOfflineSigner;
  console.log("Keplr wallet bootstraped");
  return win.keplr as Keplr;
}
