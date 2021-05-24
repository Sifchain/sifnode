import { BroadcastTxFailure } from "@cosmjs/launchpad";
import { TransactionStatus } from "../../entities";
import { ErrorCode, getErrorMessage } from "../../entities/Errors";

export function parseTxFailure(txFailure: {
  transactionHash: string;
  rawLog: string;
}): TransactionStatus {
  console.log({ "txFailure.rawLog": txFailure.rawLog });
  // TODO: synchronise with backend and use error codes at the service level
  // and provide a localized error lookup on frontend for messages
  if (txFailure.rawLog.toLowerCase().includes("below expected")) {
    return {
      code: ErrorCode.TX_FAILED_SLIPPAGE,
      hash: txFailure.transactionHash,
      memo: getErrorMessage(ErrorCode.TX_FAILED_SLIPPAGE),
      state: "failed",
    };
  }

  if (txFailure.rawLog.toLowerCase().includes("swap_failed")) {
    return {
      code: ErrorCode.TX_FAILED,
      hash: txFailure.transactionHash,
      memo: getErrorMessage(ErrorCode.TX_FAILED),
      state: "failed",
    };
  }

  if (txFailure.rawLog.toLowerCase().includes("request rejected")) {
    return {
      code: ErrorCode.USER_REJECTED,
      hash: txFailure.transactionHash,
      memo: getErrorMessage(ErrorCode.USER_REJECTED),
      state: "rejected",
    };
  }

  if (txFailure.rawLog.toLowerCase().includes("out of gas")) {
    return {
      code: ErrorCode.TX_FAILED_OUT_OF_GAS,
      hash: txFailure.transactionHash,
      memo: getErrorMessage(ErrorCode.TX_FAILED_OUT_OF_GAS),
      state: "out_of_gas",
    };
  }

  if (txFailure.rawLog.toLowerCase().includes("insufficient funds")) {
    return {
      code: ErrorCode.INSUFFICIENT_FUNDS,
      hash: txFailure.transactionHash,
      memo: getErrorMessage(ErrorCode.INSUFFICIENT_FUNDS),
      state: "failed",
    };
  }

  if (
    txFailure.rawLog.toLowerCase().includes("user does not have enough balance")
  ) {
    return {
      code: ErrorCode.TX_FAILED_USER_NOT_ENOUGH_BALANCE,
      hash: txFailure.transactionHash,
      memo: getErrorMessage(ErrorCode.TX_FAILED_USER_NOT_ENOUGH_BALANCE),
      state: "failed",
    };
  }

  return {
    code: ErrorCode.UNKNOWN_FAILURE,
    hash: txFailure.transactionHash,
    memo: getErrorMessage(ErrorCode.UNKNOWN_FAILURE),
    state: "failed",
  };
}
