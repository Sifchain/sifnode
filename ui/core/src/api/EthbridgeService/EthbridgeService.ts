import { provider } from "web3-core";
import Web3 from "web3";
import { getBridgeBankContract } from "./bridgebankContract";
import { AssetAmount, Token } from "../../entities";
import { createPegTxEventEmitter } from "./PegTxEventEmitter";
import { confirmTx } from "./utils/confirmTx";

export type EthbridgeServiceContext = {
  // sifApiUrl: string;
  bridgebankContractAddress: string;
  getWeb3Provider: () => Promise<provider>;
};

const ETH_ADDRESS = "0x0000000000000000000000000000000000000000";

export default function createEthbridgeService({
  bridgebankContractAddress,
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
    lock(
      sifRecipient: string,
      assetAmount: AssetAmount,
      confirmations: number
    ) {
      const emitter = createPegTxEventEmitter();

      function handleError(err: any) {
        emitter.emit({ type: "Error", payload: err });
      }

      (async function() {
        const web3 = await ensureWeb3();
        const cosmosRecipient = Web3.utils.utf8ToHex(sifRecipient);

        const bridgeBankContract = await getBridgeBankContract(
          web3,
          bridgebankContractAddress
        );
        const accounts = await web3.eth.getAccounts();
        const coinDenom = (assetAmount.asset as Token).address ?? ETH_ADDRESS;
        const amount = assetAmount.numerator.toString();
        const fromAddress = accounts[0];

        const sendArgs = {
          from: fromAddress,
          value: coinDenom === ETH_ADDRESS ? amount : 0,
          gas: 5000000,
        };

        bridgeBankContract.methods
          .lock(cosmosRecipient, coinDenom, amount)
          .send(sendArgs)
          .on("transactionHash", (hash: string) => {
            emitter.setTxHash(hash);
          })
          .on("error", (err: any) => {
            handleError(err);
          });

        emitter.onTxHash(({ payload: txHash }) => {
          confirmTx({
            web3,
            txHash,
            confirmations,
            onSuccess() {
              emitter.emit({ type: "Complete", payload: null });
            },
            onCheckConfirmation(count) {
              emitter.emit({ type: "EthConfCountChanged", payload: count });
            },
          });
        });
      })();

      return emitter;
    },
  };
}
