import { AssetAmount } from '../entities';
import detectMetaMaskProvider from '@metamask/detect-provider';

import Web3 from 'web3';
import { AbstractProvider } from 'web3-core';
import { ETH } from '../constants';
import JSBI from 'jsbi';

type WindowWithPossibleMetaMask = typeof window & {
  etherium?: MetaMaskProvider;
  web3: OldMetaMaskProvider;
};

type MetaMaskProvider = AbstractProvider & { enable: () => Promise<void> };
type OldMetaMaskProvider = AbstractProvider & {
  currentProvider: AbstractProvider;
};

// Not sure if this is
async function getWeb3(): Promise<Web3 | null> {
  const mmp = await detectMetaMaskProvider();
  const win = window as WindowWithPossibleMetaMask;

  if (!mmp || !win) return null;

  if (win.etherium) {
    const web3 = new Web3(win.etherium);
    await win.etherium.enable();
    return web3;
  }

  if (win.web3) {
    return new Web3(win.web3.currentProvider);
  }

  return null;
}

function createWalletService(getWeb3: () => Promise<Web3 | null>) {
  return {
    async getAssetBalances(): Promise<AssetAmount[]> {
      const web3 = await getWeb3();
      if (!web3) {
        alert('Cannot connect to wallet');
        return [];
      }
      const { eth } = web3;
      const accounts = await eth.getAccounts();
      const assetAmounts = [];
      for (let account of accounts) {
        const ethBalance = await eth.getBalance(account);

        assetAmounts.push(AssetAmount.create(ETH, ethBalance));
      }

      return assetAmounts;
    },
  };
}

export const walletService = createWalletService(getWeb3);
