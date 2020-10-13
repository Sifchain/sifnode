import detectMetaMaskProvider from "@metamask/detect-provider";
import { AbstractProvider } from "web3-core";
import Web3 from "web3";

type MetaMaskProvider = AbstractProvider & { enable: () => Promise<void> };
type OldMetaMaskProvider = AbstractProvider & {
  currentProvider: AbstractProvider;
};

type WindowWithPossibleMetaMask = typeof window & {
  ethereum?: MetaMaskProvider;
  web3: OldMetaMaskProvider;
};

// Not sure if this is
export const getWeb3: Web3Getter = async () => {
  const mmp = await detectMetaMaskProvider();
  const win = window as WindowWithPossibleMetaMask;

  if (!mmp || !win) return null;
  if (win.ethereum) {
    const web3 = new Web3(win.ethereum);
    await win.ethereum.enable();
    return web3;
  }

  if (win.web3) {
    return new Web3(win.web3.currentProvider);
  }

  return null;
};

export type Web3Getter = () => Promise<Web3 | null>;
