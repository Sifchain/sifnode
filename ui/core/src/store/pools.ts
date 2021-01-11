import { reactive } from "@vue/reactivity";

import { Pool } from "../entities";

export type PoolStore = {
  [s: string]: Pool;
};

export const pools = reactive<PoolStore>({}) as PoolStore;
