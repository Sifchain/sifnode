import detectMetaMaskProvider from "@metamask/detect-provider";
import Web3 from "web3";
import { AbstractProvider, provider } from "web3-core";

type MetaMaskProvider = AbstractProvider & {
  request?: (a: any) => Promise<void>;
};

type WindowWithPossibleMetaMask = typeof window & {
  ethereum?: MetaMaskProvider;
  web3?: Web3;
};

// Detect mossible metamask provider from browser
export const getMetamaskProvider = async (): Promise<provider> => {
  const mmp = await detectMetaMaskProvider();
  const win = window as WindowWithPossibleMetaMask;

  if (!mmp || !win) return null;

  if (mmp) {
    return mmp as provider;
  }

  // if a wallet has left web3 on the page we can use the current provider
  if (win.web3) {
    return win.web3.currentProvider;
  }

  return null;
};
