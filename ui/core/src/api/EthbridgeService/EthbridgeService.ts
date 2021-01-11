import { provider } from "web3-core";
import Web3 from "web3";
import { getBridgeBankContract } from "./bridgebankContract";
import { AssetAmount, Token } from "../../entities";
import { createPegTxEventEmitter } from "./PegTxEventEmitter";
import { confirmTx } from "./utils/confirmTx";
import { SifUnSignedClient } from "../utils/SifClient";

export type EthbridgeServiceContext = {
  sifApiUrl: string;
  sifWsUrl: string;
  sifChainId: string;
  bridgebankContractAddress: string;
  getWeb3Provider: () => Promise<provider>;
  sifUnsignedClient?: SifUnSignedClient;
};

const ETH_ADDRESS = "0x0000000000000000000000000000000000000000";

export default function createEthbridgeService({
  sifApiUrl,
  sifWsUrl,
  sifChainId,
  bridgebankContractAddress,
  getWeb3Provider,
  sifUnsignedClient = new SifUnSignedClient(sifApiUrl, sifWsUrl),
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
    async burn(params: {
      fromAddress: string;
      ethereumRecipient: string;
      assetAmount: AssetAmount;
    }) {
      const web3 = await ensureWeb3();
      const ethereumChainId = await web3.eth.net.getId();
      const tokenAddress =
        (params.assetAmount.asset as Token).address ?? ETH_ADDRESS;

      return await sifUnsignedClient.burn({
        ethereum_receiver: params.ethereumRecipient,
        base_req: {
          chain_id: sifChainId,
          from: params.fromAddress,
        },
        amount: params.assetAmount.toFixed(0),
        symbol: params.assetAmount.asset.symbol,
        cosmos_sender: params.fromAddress,
        ethereum_chain_id: `${ethereumChainId}`,
        token_contract_address: tokenAddress,
      });
    },

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
