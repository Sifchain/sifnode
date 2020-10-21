// These have mostly come off the API VOs and we may or may not need them
// As we flesh out API calls within the interface we can specify this further

import JSBI from "jsbi";
import { Address } from "./Address";
import { Asset } from "./Asset";

export type Transaction = any;
export type EncodedTransaction = { tx: string };
export type BroadcastingResult = any;

export type TxParams = {
  asset?: Asset;
  amount: JSBI;
  recipient: Address;
  feeRate?: number; // optional feeRate
  memo?: string; // optional memo to pass
};

export type TxHash = string;
