// based on draft swagger spec https://raw.githubusercontent.com/Sifchain/sifnode/c1bb5a268da8b519d0fc90f81fa194d31c0f82b3/api/openapi/swagger.yml?token=AAJSXWM6CDXYAEETSC6BJ2S7Q2JLS

import {
  BroadcastingResult,
  EncodedTransaction,
  Transaction,
} from "../entities/Transaction";

export const transactionService = {
  // GET /txs/{hash}
  async getByhash(hash: string) {},

  // GET /txs/{hash}
  async search(
    actions: string,
    sender: string,
    page: number,
    limit: number,
    txheight: number
  ): Promise<Transaction[]> {
    return [];
  },

  // POST /txs
  async broadcast(tx: Transaction): Promise<BroadcastingResult> {},

  // POST /txs/encode
  async encode(tx: Transaction): Promise<EncodedTransaction> {
    return { tx: "somehash" };
  },

  async decode(tx: EncodedTransaction): Promise<Transaction> {
    return {};
  },
};
