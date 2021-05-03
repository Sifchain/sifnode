import { provider } from "web3-core";
import Web3 from "web3";
import { getBridgeBankContract } from "./bridgebankContract";
import { getTokenContract } from "./tokenContract";
import { IAssetAmount } from "../../entities";
import {
  createPegTxEventEmitter,
  PegTxEventEmitter,
} from "./PegTxEventEmitter";
import { confirmTx } from "./utils/confirmTx";
import { SifUnSignedClient } from "../utils/SifClient";
import { parseTxFailure } from "./parseTxFailure";
import { Contract } from "web3-eth-contract";
import JSBI from "jsbi";

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
  // How to handle context/dependency injection?
  let _web3: Web3 | null = null;
  async function ensureWeb3(): Promise<Web3> {
    if (!_web3) {
      _web3 = new Web3(await getWeb3Provider());
    }
    return _web3;
  }

  /**
   * Create an event listener to report status of a peg transaction.
   * Usage:
   * const tx = createPegTx(50)
   * tx.setTxHash('0x52ds.....'); // set the hash to lookup and confirm on the blockchain
   * @param confirmations number of confirmations before pegtx is considered confirmed
   */
  function createPegTx(
    confirmations: number,
    symbol?: string,
    txHash?: string,
  ): PegTxEventEmitter {
    const emitter = createPegTxEventEmitter(txHash, symbol);

    // decorate pegtx to invert dependency to web3 and confirmations
    emitter.onTxHash(async ({ payload: txHash }) => {
      const web3 = await ensureWeb3();
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

    return emitter;
  }

  /**
   * Gets a list of transactionHashes found as _from keys within the given events within a given blockRange from the current block
   * @param {*} address eth address to correlate transactions with
   * @param {*} contract web3 contract
   * @param {*} eventList event name list of events (must have an addresskey)
   * @param {*} blockRange number of blocks from the current block header to search
   */
  async function getEventTxsInBlockrangeFromAddress(
    address: string,
    contract: Contract,
    eventList: string[],
    blockRange: number,
  ) {
    const web3 = await ensureWeb3();
    const latest = await web3.eth.getBlockNumber();
    const fromBlock = Math.max(latest - blockRange, 0);
    const allEvents = await contract.getPastEvents("allEvents", {
      // filter:{_from:address}, // if _from was indexed we could do this
      fromBlock,
      toBlock: "latest",
    });

    // unfortunately because _from is not an indexed key we have to manually filter
    // TODO: ask peggy team to index the _from field which would make this more efficient
    const txs: { symbol: string; hash: string }[] = [];
    for (let event of allEvents) {
      const isEventWeCareAbout = eventList.includes(event.event);

      const matchesInputAddress =
        event?.returnValues?._from?.toLowerCase() === address.toLowerCase();

      if (isEventWeCareAbout && matchesInputAddress && event.transactionHash) {
        txs.push({
          hash: event.transactionHash,
          symbol: event.returnValues?._symbol,
        });
      }
    }
    return txs;
  }

  return {
    async approveBridgeBankSpend(account: string, amount: IAssetAmount) {
      // This will popup an approval request in metamask
      const web3 = await ensureWeb3();
      const tokenContract = await getTokenContract(web3, amount.asset.address!);
      const sendArgs = {
        from: account,
        value: 0,
        gas: 100000,
      };

      // TODO - give interface option to approve unlimited spend via web3.utils.toTwosComplement(-1);
      // NOTE - We may want to move this out into its own separate function.
      // Although I couldn't think of a situation we'd call allowance separately from approve
      const hasAlreadyApprovedSpend = await tokenContract.methods
        .allowance(account, bridgebankContractAddress)
        .call();
      if (
        JSBI.lessThanOrEqual(
          amount.toBigInt(),
          JSBI.BigInt(hasAlreadyApprovedSpend),
        )
      ) {
        // dont request approve again
        console.log(
          "approveBridgeBankSpend: spend already approved",
          hasAlreadyApprovedSpend,
        );
        return;
      }

      const res = await tokenContract.methods
        .approve(bridgebankContractAddress, amount.toBigInt().toString())
        .send(sendArgs);
      console.log("approveBridgeBankSpend:", res);
      return res;
    },

    async burnToEthereum(params: {
      fromAddress: string;
      ethereumRecipient: string;
      assetAmount: IAssetAmount;
      feeAmount: IAssetAmount;
    }) {
      const web3 = await ensureWeb3();
      const ethereumChainId = await web3.eth.net.getId();
      const tokenAddress = params.assetAmount.asset.address ?? ETH_ADDRESS;
      console.log("burnToEthereum: start: ", tokenAddress);

      const txReceipt = await sifUnsignedClient.burn({
        ethereum_receiver: params.ethereumRecipient,
        base_req: {
          chain_id: sifChainId,
          from: params.fromAddress,
        },
        amount: params.assetAmount.toBigInt().toString(),
        symbol: params.assetAmount.asset.symbol,
        cosmos_sender: params.fromAddress,
        ethereum_chain_id: `${ethereumChainId}`,
        token_contract_address: tokenAddress,
        ceth_amount: params.feeAmount.toBigInt().toString(),
      });

      console.log("burnToEthereum: txReceipt: ", txReceipt, tokenAddress);
      return txReceipt;
    },

    lockToSifchain(
      sifRecipient: string,
      assetAmount: IAssetAmount,
      confirmations: number,
    ) {
      const pegTx = createPegTx(confirmations, assetAmount.asset.symbol);

      function handleError(err: any) {
        console.log("lockToSifchain: handleError: ", err);
        pegTx.emit({
          type: "Error",
          payload: parseTxFailure({ hash: "", log: err.message.toString() }),
        });
      }

      (async function () {
        const web3 = await ensureWeb3();
        const cosmosRecipient = Web3.utils.utf8ToHex(sifRecipient);

        const bridgeBankContract = await getBridgeBankContract(
          web3,
          bridgebankContractAddress,
        );
        const accounts = await web3.eth.getAccounts();
        const coinDenom = assetAmount.asset.address || ETH_ADDRESS; // eth address is ""
        const amount = assetAmount.toBigInt().toString();
        const fromAddress = accounts[0];

        const sendArgs = {
          from: fromAddress,
          value: coinDenom === ETH_ADDRESS ? amount : 0,
          gas: 150000,
        };

        console.log(
          "lockToSifchain: bridgeBankContract.lock",
          JSON.stringify({ cosmosRecipient, coinDenom, amount, sendArgs }),
        );

        bridgeBankContract.methods
          .lock(cosmosRecipient, coinDenom, amount)
          .send(sendArgs)
          .on("transactionHash", (hash: string) => {
            console.log("lockToSifchain: bridgeBankContract.lock TX", hash);
            pegTx.setTxHash(hash);
          })
          .on("error", (err: any) => {
            console.log("lockToSifchain: bridgeBankContract.lock ERROR", err);
            handleError(err);
          });
      })().catch((err) => {
        handleError(err);
      });

      return pegTx;
    },

    async lockToEthereum(params: {
      fromAddress: string;
      ethereumRecipient: string;
      assetAmount: IAssetAmount;
      feeAmount: IAssetAmount;
    }) {
      const web3 = await ensureWeb3();
      const ethereumChainId = await web3.eth.net.getId();
      const tokenAddress = params.assetAmount.asset.address ?? ETH_ADDRESS;

      const lockParams = {
        ethereum_receiver: params.ethereumRecipient,
        base_req: {
          chain_id: sifChainId,
          from: params.fromAddress,
        },
        amount: params.assetAmount.toBigInt().toString(),
        symbol: params.assetAmount.asset.symbol,
        cosmos_sender: params.fromAddress,
        ethereum_chain_id: `${ethereumChainId}`,
        token_contract_address: tokenAddress,
        ceth_amount: params.feeAmount.toBigInt().toString(),
      };

      console.log("lockToEthereum: TRY LOCK", tokenAddress);
      const lockReceipt = await sifUnsignedClient.lock(lockParams);
      console.log("lockToEthereum: LOCKED", lockReceipt);

      return lockReceipt;
    },

    /**
     * Get a list of unconfirmed transaction hashes associated with
     * a particular address and return pegTxs associated with that hash
     * @param address contract address
     * @param confirmations number of confirmations required
     */
    async fetchUnconfirmedLockBurnTxs(
      address: string,
      confirmations: number,
    ): Promise<PegTxEventEmitter[]> {
      const web3 = await ensureWeb3();

      const bridgeBankContract = await getBridgeBankContract(
        web3,
        bridgebankContractAddress,
      );

      const txs = await getEventTxsInBlockrangeFromAddress(
        address,
        bridgeBankContract,
        ["LogBurn", "LogLock"],
        confirmations,
      );

      return txs.map(({ hash, symbol }) =>
        createPegTx(confirmations, symbol, hash),
      );
    },

    burnToSifchain(
      sifRecipient: string,
      assetAmount: IAssetAmount,
      confirmations: number,
      account?: string,
    ) {
      const pegTx = createPegTx(confirmations, assetAmount.asset.symbol);

      function handleError(err: any) {
        console.log("burnToSifchain: handleError ERROR", err);
        pegTx.emit({
          type: "Error",
          payload: parseTxFailure({ hash: "", log: err }),
        });
      }

      (async function () {
        const web3 = await ensureWeb3();
        const cosmosRecipient = Web3.utils.utf8ToHex(sifRecipient);

        const bridgeBankContract = await getBridgeBankContract(
          web3,
          bridgebankContractAddress,
        );
        const accounts = await web3.eth.getAccounts();
        const coinDenom = assetAmount.asset.address;
        const amount = assetAmount.toBigInt().toString();
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
            pegTx.setTxHash(hash);
          })
          .on("error", (err: any) => {
            console.log("burnToSifchain: bridgeBankContract.burn ERROR", err);
            handleError(err);
          });
      })().catch((err) => {
        handleError(err);
      });

      return pegTx;
    },
  };
}
