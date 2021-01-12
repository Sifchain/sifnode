import { OfflineSigner } from "@cosmjs/launchpad";

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

// Detect mossible keplr provider from browser
export default (): provider | null => {
  const win = window as WindowWithPossibleKeplr;

  if (!win) return null;

  if (win.keplr && win.getOfflineSigner) {
    // assign offline signer (they use __proto__ for some reason), so this is not as pretty as i'd like)
    Object.getPrototypeOf(win.keplr).getOfflineSigner = win.getOfflineSigner;
    return win.keplr as Keplr;
  }

  return null;
};
