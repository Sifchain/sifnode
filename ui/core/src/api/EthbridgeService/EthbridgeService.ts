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
  sifRpcUrl: string;
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
  sifRpcUrl,
  sifChainId,
  bridgebankContractAddress,
  getWeb3Provider,
  sifUnsignedClient = new SifUnSignedClient(sifApiUrl, sifWsUrl, sifRpcUrl),
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
    async approveBridgeBankSpend(account: string, amount: AssetAmount) {
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

      // TODO - give interface option to approve unlimited spend via web3.utils.toTwosComplement(-1);
      // NOTE - We may want to move this out into its own separate function.
      // Although I couldn't think of a situation we'd call allowance separately from approve
      const hasAlreadyApprovedSpend = await tokenContract.methods
        .allowance(account, bridgebankContractAddress)
        .call();
      if (hasAlreadyApprovedSpend >= amount.toBaseUnits().toString()) {
        // dont request approve again
        console.log(
          "approveBridgeBankSpend: spend already approved",
          hasAlreadyApprovedSpend
        );
        return;
      }

      const res = await tokenContract.methods
        .approve(bridgebankContractAddress, amount.toBaseUnits().toString())
        .send(sendArgs);
      console.log("approveBridgeBankSpend:", res);
      return res;
    },

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
      console.log("burnToEthereum: start: ", tokenAddress);

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

      console.log("burnToEthereum: txReceipt: ", txReceipt, tokenAddress);
      return txReceipt;
    },

    lockToSifchain(
      sifRecipient: string,
      assetAmount: AssetAmount,
      confirmations: number
    ) {
      const emitter = createPegTxEventEmitter();

      function handleError(err: any) {
        console.log("lockToSifchain: handleError: ", err);
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

        const sendArgs = {
          from: fromAddress,
          value: coinDenom === ETH_ADDRESS ? amount : 0,
        };

        console.log(
          "lockToSifchain: bridgeBankContract.lock",
          JSON.stringify({ cosmosRecipient, coinDenom, amount, sendArgs })
        );

        bridgeBankContract.methods
          .lock(cosmosRecipient, coinDenom, amount)
          .send(sendArgs)
          .on("transactionHash", (hash: string) => {
            console.log("lockToSifchain: bridgeBankContract.lock TX", hash);
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
              console.log(
                "lockToSifchain: confirmTx SUCCESS",
                txHash,
                confirmations
              );
              emitter.emit({ type: "Complete", payload: null });
            },
            onCheckConfirmation(count) {
              console.log(
                "lockToSifchain: onCheckConfirmation PENDING",
                confirmations
              );
              emitter.emit({ type: "EthConfCountChanged", payload: count });
            },
          });
        });
      })().catch(err => {
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

      console.log("lockToEthereum: TRY LOCK", tokenAddress);
      const lockReceipt = await sifUnsignedClient.lock(lockParams);
      console.log("lockToEthereum: LOCKED", lockReceipt);

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
        console.log("burnToSifchain: handleError ERROR", err);
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
          gas: 150000, // Note: This chose in lieu of burn(params).estimateGas({from})
        };

        bridgeBankContract.methods
          .burn(cosmosRecipient, coinDenom, amount)
          .send(sendArgs)
          .on("transactionHash", (hash: string) => {
            console.log("burnToSifchain: bridgeBankContract.burn TX", hash);
            emitter.setTxHash(hash);
          })
          .on("error", (err: any) => {
            console.log("burnToSifchain: bridgeBankContract.burn ERROR", err);
            handleError(err);
          });

        emitter.onTxHash(({ payload: txHash }) => {
          console.log("Waiting for confirmation... ");
          confirmTx({
            web3,
            txHash,
            confirmations,
            onSuccess() {
              console.log(
                "burnToSifchain: commitTx SUCCESS",
                txHash,
                confirmations
              );
              emitter.emit({ type: "Complete", payload: null });
            },
            onCheckConfirmation(count) {
              console.log(
                "burnToSifchain: commitTx.checkConfirmation PENDING",
                confirmations
              );
              emitter.emit({ type: "EthConfCountChanged", payload: count });
            },
          });
        });
      })().catch(err => {
        handleError(err);
      });

      return emitter;
    },
  };
}
