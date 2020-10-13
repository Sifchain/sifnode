import { AssetAmount, Token } from "../entities";
import { ETH } from "../constants";
import { Web3Getter } from "./utils/getWeb3";

export type WalletServiceContext = {
  getWeb3: Web3Getter;
  getSupportedTokens: () => Promise<Map<string, Token>>;
};

export default function createWalletService({
  getWeb3,
  getSupportedTokens,
}: WalletServiceContext) {
  return {
    async getAssetBalances(): Promise<AssetAmount[]> {
      const web3 = await getWeb3();
      const supportedTokens = await getSupportedTokens();

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

      for (const [_, token] of supportedTokens) {
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
