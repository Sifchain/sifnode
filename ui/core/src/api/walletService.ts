import { AssetAmount } from "../entities";
import detectMetaMaskProvider from "@metamask/detect-provider";

import Web3 from "web3";
import { AbstractProvider } from "web3-core";
import { ETH, USDC } from "../constants";

const SUPPORTED_TOKENS = [USDC];

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
        alert("Cannot connect to wallet");
        return [];
      }
      const { eth } = web3;
      const accounts = await eth.getAccounts();
      const assetAmounts: AssetAmount[] = [];
      for (const account of accounts) {
        const ethBalance = await eth.getBalance(account);

        assetAmounts.push(
          AssetAmount.create(ETH, web3.utils.fromWei(ethBalance, "microether"))
        );

        for (const token of SUPPORTED_TOKENS) {
          const contract = new eth.Contract(
            [
              // balanceOf
              {
                constant: true,
                inputs: [{ name: "_owner", type: "address" }],
                name: "balanceOf",
                outputs: [{ name: "balance", type: "uint256" }],
                type: "function",
              },
              // decimals
              {
                constant: true,
                inputs: [],
                name: "decimals",
                outputs: [{ name: "", type: "uint8" }],
                type: "function",
              },
            ],
            token.address
          );

          const balanceOfErc = await contract.methods.balanceOf(account).call();

          console.log({ balanceOfErc, token: token.symbol });
          assetAmounts.push(AssetAmount.create(token, balanceOfErc));
        }
      }

      return assetAmounts;
    },
  };
}

export const walletService = createWalletService(getWeb3);
