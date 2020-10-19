// These have mostly come off the API VOs and we may or may not need them
// As we flesh out API calls within the interface we can specify this further

import { SifAddress } from "./Wallet";

export type Transaction = any;
export type EncodedTransaction = { tx: string };
export type BroadcastingResult = any;

export type SifTransaction = {
  amount?: Number
  denom?: string
  to_address?: SifAddress
  memo?: string
}