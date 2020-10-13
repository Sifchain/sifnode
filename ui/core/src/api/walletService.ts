import { AssetAmount, Token } from "../entities";
import detectMetaMaskProvider from "@metamask/detect-provider";

import Web3 from "web3";
import { AbstractProvider } from "web3-core";
import { ETH } from "../constants";

// const SUPPORTED_TOKENS = [ATK, BTK];

type WindowWithPossibleMetaMask = typeof window & {
  ethereum?: MetaMaskProvider;
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
  if (win.ethereum) {
    const web3 = new Web3(win.ethereum);
    await win.ethereum.enable();
    return web3;
  }

  if (win.web3) {
    return new Web3(win.web3.currentProvider);
  }

  return null;
}

export function createWalletService(
  getWeb3: () => Promise<Web3 | null>,
  supportedTokens: Token[]
) {
  return {
    async getAssetBalances(): Promise<AssetAmount[]> {
      const web3 = await getWeb3();
      if (!web3) {
        alert("Cannot connect to wallet");
        return [];
      }
      const { eth } = web3;
      const [account] = await eth.getAccounts();

      const assetAmounts: AssetAmount[] = [];

      // This is going to give us all the acounts on the node.
      // Not sure if this is the right thing to do here.
      // So Commenting it out for now
      // for (const account of mainAccount) {
      const ethBalance = await eth.getBalance(account);

      assetAmounts.push(AssetAmount.create(ETH, ethBalance));

      for (const token of supportedTokens) {
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

        assetAmounts.push(AssetAmount.create(token, balanceOfErc));
      }
      // }

      return assetAmounts;
    },
  };
}

export const walletService = createWalletService(getWeb3, []);
