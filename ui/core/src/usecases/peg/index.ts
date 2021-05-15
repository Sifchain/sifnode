import { UsecaseContext } from "..";
import {
  Address,
  Asset,
  AssetAmount,
  IAsset,
  IAssetAmount,
  Network,
  TransactionStatus,
} from "../../entities";
import { isSupportedEVMChain } from "../utils";

import { SubscribeToUnconfirmedPegTxs } from "./subscribeToUnconfirmedPegTxs";
import { SubscribeToTx } from "./utils/subscribeToTx";

function isOriginallySifchainNativeToken(asset: Asset) {
  return ["erowan", "rowan"].includes(asset.symbol);
}

// TODO: Subscriptions, Commands and Queries should all be their own concepts and each exist within their
//       own files to manage complexity allow for team to grow and avoid refactoring
//       subtle complexity of dependency injection to maintain testability is required passing in ctx below

/**
 * Shared peg config for use throughout the peg feature
 */
export type PegConfig = { ethConfirmations: number };

export default ({
  services,
  store,
}: UsecaseContext<
  // Once we have moved all interactors to their own files this can be
  // UsecaseContext<any,any> or renamed to InteractorContext<any,any>
  "sif" | "ethbridge" | "bus" | "eth", // Select the services you want to access
  "wallet" | "tx" // Select the store keys you want to access
>) => {
  const config: PegConfig = {
    // listen for 50 confirmations
    // Eventually this should be set on ebrelayer
    // to centralize the business logic
    ethConfirmations: 50,
  };

  // Create the context for passing to commands, queries and subscriptions
  const ctx = { services, store, config };

  /* 
    TODO: suggestion externalize all interactors injecting ctx would look like the following

    const commands = {
      unpeg: Unpeg(ctx),
      peg: Peg(ctx),
    }

    const queries = {
      getSifTokens: GetSifTokens(ctx),
      getEthTokens: GetEthTokens(ctx),
      calculateUnpegFee: CalculateUnpegFee(ctx),
    }
    
    const subscriptions = {
      subscribeToUnconfirmedPegTxs: SubscribeToUnconfirmedPegTxs(ctx),
    }
  */

  // Rename and split this up to subscriptions, commands, queries
  const actions = {
    subscribeToUnconfirmedPegTxs: SubscribeToUnconfirmedPegTxs(ctx),

    getSifTokens() {
      return services.sif.getSupportedTokens();
    },

    getEthTokens() {
      return services.eth.getSupportedTokens();
    },

    calculateUnpegFee(asset: IAsset) {
      const feeNumber = isOriginallySifchainNativeToken(asset)
        ? "70000000000000000"
        : "70000000000000000";

      return AssetAmount(Asset.get("ceth"), feeNumber);
    },

    async unpeg(assetAmount: IAssetAmount) {
      const lockOrBurnFn = isOriginallySifchainNativeToken(assetAmount.asset)
        ? services.ethbridge.lockToEthereum
        : services.ethbridge.burnToEthereum;

      const feeAmount = this.calculateUnpegFee(assetAmount.asset);

      const tx = await lockOrBurnFn({
        assetAmount,
        ethereumRecipient: store.wallet.eth.address,
        fromAddress: store.wallet.sif.address,
        feeAmount,
      });

      console.log(
        "unpeg",
        tx,
        assetAmount,
        store.wallet.eth.address,
        store.wallet.sif.address,
        feeAmount,
      );

      const txStatus = await services.sif.signAndBroadcast(tx.value.msg);

      if (txStatus.state !== "accepted") {
        services.bus.dispatch({
          type: "PegTransactionErrorEvent",
          payload: {
            txStatus,
            message: txStatus.memo || "There was an error while unpegging",
          },
        });
      }

      console.log(
        "unpeg txStatus.state",
        txStatus.state,
        txStatus.memo,
        txStatus.code,
        tx.value.msg,
      );

      return txStatus;
    },

    // TODO: Move this approval command to within peg and report status via store or some other means
    //       This has been done for convenience but we should not have to know in the view that
    //       approval is required before pegging as that is very much business domain knowledge
    async approve(address: Address, assetAmount: IAssetAmount) {
      return await services.ethbridge.approveBridgeBankSpend(
        address,
        assetAmount,
      );
    },

    async peg(assetAmount: IAssetAmount): Promise<TransactionStatus> {
      if (
        assetAmount.asset.network === Network.ETHEREUM &&
        !isSupportedEVMChain(store.wallet.eth.chainId)
      ) {
        services.bus.dispatch({
          type: "ErrorEvent",
          payload: {
            message: "EVM Network not supported!",
          },
        });
        return {
          hash: "",
          state: "failed",
        };
      }

      const subscribeToTx = SubscribeToTx(ctx);

      const lockOrBurnFn = isOriginallySifchainNativeToken(assetAmount.asset)
        ? services.ethbridge.burnToSifchain
        : services.ethbridge.lockToSifchain;

      return await new Promise<TransactionStatus>((done) => {
        const pegTx = lockOrBurnFn(
          store.wallet.sif.address,
          assetAmount,
          config.ethConfirmations,
        );

        subscribeToTx(pegTx);

        pegTx.onTxHash((hash) => {
          done({
            hash: hash.txHash,
            memo: "Transaction Accepted",
            state: "accepted",
          });
        });
      });
    },
  };

  return actions;
};
