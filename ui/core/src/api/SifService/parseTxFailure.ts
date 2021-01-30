import { BroadcastTxFailure } from "@cosmjs/launchpad";
import { TransactionStatus } from "../../entities";

export function parseTxFailure(
  txFailure: BroadcastTxFailure
): TransactionStatus {
  // This is rough and a little brittle for now
  // TODO: synchronise with backend and use error codes at the service level
  // and provide a localized error lookup on frontend for messages
  if (txFailure.rawLog.toLowerCase().includes("below expected")) {
    return {
      hash: txFailure.transactionHash,
      memo: "Swap failed - Received amount is below expected",
      state: "failed",
    };
  }

  if (txFailure.rawLog.toLowerCase().includes("swap_failed")) {
    return {
      hash: txFailure.transactionHash,
      memo: "Swap failed",
      state: "failed",
    };
  }

  if (txFailure.rawLog.toLowerCase().includes("request rejected")) {
    return {
      hash: txFailure.transactionHash,
      memo: "Request Rejected",
      state: "rejected",
    };
  }

  return {
    hash: txFailure.transactionHash,
    memo: "Unknown failure",
    state: "failed",
  };
}
