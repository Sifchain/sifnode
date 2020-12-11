import { AssetAmount } from "../../entities";
import { TxEventEmitter } from "./types";

export type IPeggyService = {
  /**
   * Release funds from the Ethereum Smart Contract and burn the equivalent tokens in sifnode
   * @param ethereumRecipient Ethereum address to send funds to
   * @param assetAmount amount of funds and sif asset eg ceth
   */
  burn(ethereumRecipient: string, assetAmount: AssetAmount): TxEventEmitter;

  /**
   * Lock funds in the Ethereum Smart Contract and mint the equivalent tokens in sifnode
   * @param cosmosRecipient sif address to send funds to
   * @param assetAmount amount of funds and eth asset eg erowan
   */
  lock(cosmosRecipient: string, assetAmount: AssetAmount): TxEventEmitter;
};
