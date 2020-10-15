import { Address, Asset, Balance, Token } from "../../entities";

import { Web3Getter } from "../utils/getWeb3";
import { EthProvider } from "./EthProvider";
import { Balances } from "../../entities/Balance";

export type WalletServiceContext = {
  getWeb3: Web3Getter;
  getSupportedTokens: () => Promise<Token[]>;
};

export default function createWalletService({
  getWeb3,
  getSupportedTokens,
}: WalletServiceContext) {
  let ethWallet: EthProvider | undefined;

  return {
    async disconnect() {
      ethWallet = undefined;
    },
    async connect(): Promise<boolean> {
      const web3 = await getWeb3();
      if (!web3) {
        console.log("Cound not connect to wallet");
        return false;
      }
      const [address] = await web3.eth.getAccounts();
      const supportedTokens = await getSupportedTokens();
      ethWallet = new EthProvider(address, web3, supportedTokens);
      return true;
    },

    async getBalance(
      address?: Address,
      asset?: Asset | Token
    ): Promise<Balances> {
      if (!ethWallet) return [];
      const balances = await ethWallet.getBalance(address, asset);
      return balances || [];
    },

    isConnected() {
      return Boolean(ethWallet);
    },

    // FOLLOWING ARE YTI:

    // setPhrase(phrase: string): Address
    // setNetwork(net: Network): void
    // getNetwork(): Network

    // getExplorerAddressUrl(address: Address): string
    // getExplorerTxUrl(txID: string): string
    // getTransactions(params?: TxHistoryParams): Promise<TxsPage>

    // getFees(): Promise<Fees>

    // transfer(params: TxParams): Promise<TxHash>
    // deposit(params: TxParams): Promise<TxHash>

    // purgeClient(): void

    // async getAssetBalances({
    //   limit = -1,
    //   supportedTokens = [],
    // }: {
    //   supportedTokens?: Token[];
    //   limit?: number;
    // }): Promise<Balance[]> {
    //   if (!this.isConnected()) {
    //     console.log("Cannot connect to wallet");
    //     return [];
    //   }

    //   const assetAmounts: Balance[] = [];

    //   // This is going to give us all the acounts on the node.
    //   // Not sure if this is the right thing to do here.
    //   // So Commenting it out for now
    //   // for (const account of mainAccount) {

    //   assetAmounts.push(Balance.create(ETH, ethBalance));

    //   for (let i = 0; i < supportedTokens.length; i++) {
    //     if (limit > -1 && i > limit) break;

    //     const token = supportedTokens[i];
    //     const contract = new eth.Contract(generalTokenAbi, token.address);
    //     const balanceOfErc = await contract.methods.balanceOf(account).call();

    //     assetAmounts.push(Balance.create(token, balanceOfErc));
    //   }
    //   // }

    //   return assetAmounts;
    // },
  };
}
