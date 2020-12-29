import { provider, TransactionReceipt } from "web3-core";
import Web3 from "web3";
import { getBridgeBankContract } from "./bridgebankContract";
import { AssetAmount } from "../../entities";

export type EthbridgeServiceContext = {
  // sifApiUrl: string;
  bridgeBankContractAddress: string;
  getWeb3Provider: () => Promise<provider>;
};

const NULL_ADDRESS = "0x0000000000000000000000000000000000000000";

export default function createEthbridgeService({
  bridgeBankContractAddress,
  getWeb3Provider,
}: EthbridgeServiceContext) {
  // Pull this out to a util?
  let _web3: Web3 | null = null;
  async function ensureWeb3(): Promise<Web3> {
    if (!_web3) {
      _web3 = new Web3(await getWeb3Provider());
    }
    return _web3;
  }

  return {
    async lock(cosmosRecipient: string, assetAmount: AssetAmount) {
      const web3 = await ensureWeb3();
      const bridgeBankContract = getBridgeBankContract(
        web3,
        bridgeBankContractAddress
      );
      const accounts = await web3.eth.getAccounts();
      const coinDenom = assetAmount.asset.symbol;
      const amount = assetAmount.amount.toString();
      const fromAddress = accounts[0];

      return new Promise((resolve, reject) => {
        let hash: string;
        let receipt: TransactionReceipt;

        function resolvePromise() {
          if (receipt && hash) resolve(hash);
        }

        bridgeBankContract.methods
          .lock(cosmosRecipient, coinDenom, amount, {
            from: fromAddress,
            value: coinDenom === NULL_ADDRESS ? amount : 0,
            gas: 300000,
          })
          .send()
          .on("transactionHash", (_hash: string) => {
            hash = _hash;
            resolvePromise();
          })
          .on("receipt", (_receipt: any) => {
            receipt = _receipt;
            resolvePromise();
          })
          .on("error", (err: any) => {
            reject(err);
          });
      });
    },
  };
}
