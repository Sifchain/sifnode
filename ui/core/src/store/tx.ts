import { reactive } from "@vue/reactivity";

import { TransactionStatus } from "../entities";

export type TxStore = {
  hash: { [hash: string]: TransactionStatus };
};

export const tx = reactive<TxStore>({ hash: {} }) as TxStore;
