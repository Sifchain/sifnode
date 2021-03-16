import { TransactionStatus } from "../../entities";
import { ErrorCode, getErrorMessage } from "../../entities/Errors";

// TODO: Should this go in a shared ethereum client mimicking sifchain?
export function parseTxFailure({
  hash = "",
  log = "",
}: {
  hash: string;
  log: string;
}): TransactionStatus {
  // Ethereum events
  if (
    log
      .toString()
      .toLowerCase()
      .includes("request rejected") ||
    log
      .toString()
      .toLowerCase()
      .includes("user denied transaction")
  ) {
    return {
      code: ErrorCode.USER_REJECTED,
      memo: getErrorMessage(ErrorCode.USER_REJECTED),
      hash,
      state: "rejected",
    };
  }

  return {
    code: ErrorCode.UNKNOWN_FAILURE,
    memo: getErrorMessage(ErrorCode.UNKNOWN_FAILURE),
    hash,
    state: "failed",
  };
}
