import { provider } from "web3-core";
import Web3 from "web3";
import { getBridgeBankContract } from "./bridgebankContract";
import { getTokenContract } from "./tokenContract";
import { AssetAmount, Token } from "../../entities";
import { createPegTxEventEmitter } from "./PegTxEventEmitter";
import { confirmTx } from "./utils/confirmTx";
import { SifUnSignedClient } from "../utils/SifClient";
import { parseTxFailure } from "./parseTxFailure";

// TODO: Do we break this service out to ethbridge and cosmos?

export type EthbridgeServiceContext = {
  sifApiUrl: string;
  sifWsUrl: string;
  sifChainId: string;
  bridgebankContractAddress: string;
  bridgetokenContractAddress: string;
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

  async function approveBridgeBankSpend(account: string, amount: AssetAmount) {
    // This will popup an approval request in metamask
    const web3 = await ensureWeb3();
    const tokenContract = await getTokenContract(
      web3,
      (amount.asset as Token).address
    );

    const sendArgs = {
      from: account,
      value: 0,
    };

    // Hmm what happens when there is a signing failure but we have approved bridgebank
    return await new Promise((resolve, reject) => {
      tokenContract.methods
        .approve(bridgebankContractAddress, amount.toBaseUnits().toString())
        .send(sendArgs)
        .on("transactionHash", (hash: string) => {
          resolve(hash);
        })
        .on("error", (err: any) => {
          console.log("lockToSifchain: bridgeBankContract.lock ERROR", err);
          reject(err);
        });
    });
  }

  return {
    async burnToEthereum(params: {
      fromAddress: string;
      ethereumRecipient: string;
      assetAmount: AssetAmount;
      feeAmount: AssetAmount;
    }) {
      const web3 = await ensureWeb3();
      const ethereumChainId = await web3.eth.net.getId();
      const tokenAddress =
        (params.assetAmount.asset as Token).address ?? ETH_ADDRESS;

      const txReceipt = await sifUnsignedClient.burn({
        ethereum_receiver: params.ethereumRecipient,
        base_req: {
          chain_id: sifChainId,
          from: params.fromAddress,
        },
        amount: params.assetAmount.toBaseUnits().toString(),
        symbol: params.assetAmount.asset.symbol,
        cosmos_sender: params.fromAddress,
        ethereum_chain_id: `${ethereumChainId}`,
        token_contract_address: tokenAddress,
        ceth_amount: params.feeAmount.toBaseUnits().toString(),
      });

      console.log("VVVVV12", txReceipt, tokenAddress);

      return txReceipt;
    },

    lockToSifchain(
      sifRecipient: string,
      assetAmount: AssetAmount,
      confirmations: number
    ) {
      const emitter = createPegTxEventEmitter();

      function handleError(err: any) {
        console.log("VVVVV13", err, "OOPS");
        emitter.emit({
          type: "Error",
          payload: parseTxFailure({ hash: "", log: err.message.toString() }),
        });
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
        console.log("VVVVV14", fromAddress, coinDenom);

        const sendArgs = {
          from: fromAddress,
          value: coinDenom === ETH_ADDRESS ? amount : 0,
        };

        console.log(
          "lockToSifchain: bridgeBankContract.lock",
          JSON.stringify({ cosmosRecipient, coinDenom, amount, sendArgs })
        );

        if (coinDenom !== ETH_ADDRESS) {
          console.log("VVVVV1", coinDenom, ETH_ADDRESS);
          await approveBridgeBankSpend(fromAddress, assetAmount);
        }

        bridgeBankContract.methods
          .lock(cosmosRecipient, coinDenom, amount)
          .send(sendArgs)
          .on("transactionHash", (hash: string) => {
            emitter.setTxHash(hash);
          })
          .on("error", (err: any) => {
            console.log("lockToSifchain: bridgeBankContract.lock ERROR", err);
            handleError(err);
          });

        emitter.onTxHash(({ payload: txHash }) => {
          confirmTx({
            web3,
            txHash,
            confirmations,
            onSuccess() {
              console.log("lockToSifchain: bridgeBankContract.lock complete");
              console.log("VVVVV2", "PEG", "COMPLETE");
              emitter.emit({ type: "Complete", payload: null });
            },
            onCheckConfirmation(count) {
              console.log("VVVVV3", count, "COMFIRMATIONS");
              emitter.emit({ type: "EthConfCountChanged", payload: count });
            },
          });
        });
      })().catch(err => {
        console.log("VVVVV4", err, "SHOOT");
        handleError(err);
      });

      return emitter;
    },

    async lockToEthereum(params: {
      fromAddress: string;
      ethereumRecipient: string;
      assetAmount: AssetAmount;
      feeAmount: AssetAmount;
    }) {
      const web3 = await ensureWeb3();
      const ethereumChainId = await web3.eth.net.getId();
      const tokenAddress =
        (params.assetAmount.asset as Token).address ?? ETH_ADDRESS;

      const lockParams = {
        ethereum_receiver: params.ethereumRecipient,
        base_req: {
          chain_id: sifChainId,
          from: params.fromAddress,
        },
        amount: params.assetAmount.toBaseUnits().toString(),
        symbol: params.assetAmount.asset.symbol,
        cosmos_sender: params.fromAddress,
        ethereum_chain_id: `${ethereumChainId}`,
        token_contract_address: tokenAddress,
        ceth_amount: params.feeAmount.toBaseUnits().toString(),
      };

      const lockReceipt = await sifUnsignedClient.lock(lockParams);
      console.log("VVVVV5", lockReceipt, "LOCK");

      return lockReceipt;
    },

    burnToSifchain(
      sifRecipient: string,
      assetAmount: AssetAmount,
      confirmations: number,
      account?: string
    ) {
      const emitter = createPegTxEventEmitter();

      function handleError(err: any) {
        console.log("VVVVV6", err, "SHOOT");
        emitter.emit({
          type: "Error",
          payload: parseTxFailure({ hash: "", log: err }),
        });
      }

      (async function() {
        const web3 = await ensureWeb3();
        const cosmosRecipient = Web3.utils.utf8ToHex(sifRecipient);

        const bridgeBankContract = await getBridgeBankContract(
          web3,
          bridgebankContractAddress
        );
        const accounts = await web3.eth.getAccounts();
        const coinDenom = (assetAmount.asset as Token).address;
        const amount = assetAmount.numerator.toString();
        const fromAddress = account || accounts[0];

        const sendArgs = {
          from: fromAddress,
          value: 0,
        };
        console.log("VVVVV7", fromAddress, assetAmount);

        await approveBridgeBankSpend(fromAddress, assetAmount);

        bridgeBankContract.methods
          .burn(cosmosRecipient, coinDenom, amount)
          .send(sendArgs)
          .on("transactionHash", (hash: string) => {
            emitter.setTxHash(hash);
          })
          .on("error", (err: any) => {
            console.log("VVVVV8", err, "OGHHH");
            console.log("lockToSifchain: bridgeBankContract.burn ERROR", err);
            handleError(err);
          });

        emitter.onTxHash(({ payload: txHash }) => {
          console.log("Waiting for confirmation... ");
          confirmTx({
            web3,
            txHash,
            confirmations,
            onSuccess() {
              console.log("VVVVV9", "UNPEG", "COMPLETE");
              emitter.emit({ type: "Complete", payload: null });
            },
            onCheckConfirmation(count) {
              console.log("VVVVV10", count, "CONFIRMATIONS");
              emitter.emit({ type: "EthConfCountChanged", payload: count });
            },
          });
        });
      })().catch(err => {
        console.log("VVVVV11", err, "OUCH");
        handleError(err);
      });

      return emitter;
    },
  };
}
