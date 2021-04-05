import JSBI from "jsbi";
import { Address } from "./Address";
import { Asset } from "./Asset";

export type TransactionStatus = {
  code?: number;
  hash: string;
  state:
    | "requested"
    | "accepted"
    | "failed"
    | "rejected"
    | "out_of_gas"
    | "completed"; // Do we need to differentiate between failed and rejected here?
  memo?: string;
  symbol?: string;
};

export type TxParams = {
  asset?: Asset;
  amount: JSBI;
  recipient: Address;
  feeRate?: number; // optional feeRate
  memo?: string; // optional memo to pass
};

export type TxHash = string;
