import detectMetaMaskProvider from "@metamask/detect-provider";
import Web3 from "web3";
import { AbstractProvider } from "web3-core";

import { Token } from "../../../entities";

type MetaMaskProvider = AbstractProvider & {
  request?: (a: any) => Promise<void>;
};

type WindowWithPossibleMetaMask = typeof window & {
  ethereum?: MetaMaskProvider;
  web3?: Web3;
};

// Detect mossible metamask provider from browser
export const getWeb3Provider = async () => {
  const mmp = await detectMetaMaskProvider();
  const win = window as WindowWithPossibleMetaMask;

  if (!mmp || !win) return null;
  if (win.ethereum) {
    // Let's test for Metamask
    if (win.ethereum.request) {
      // If metamask lets try and connect
      try {
        await win.ethereum.request({ method: "eth_requestAccounts" });
      } catch (err) {
        console.error(err);
        return null;
      }
    }
    return win.ethereum;
  }

  // if a wallet has left web3 on the page we can use the current provider
  if (win.web3) {
    return win.web3.currentProvider;
  }

  return null;
};

// export type Web3Getter = () => Promise<Web3 | null>;
export type TokensGetter = () => Promise<Token[]>;
