import { reactive } from "@vue/reactivity";

import { TransactionStatus } from "../entities";

// Store for reporting on current tx status
export type TxStore = {
  // txs as required by blockchain address
  eth: {
    [address: string]: {
      [hash: string]: TransactionStatus;
    };
  };
};

export const tx = reactive<TxStore>({ eth: {} }) as TxStore;
