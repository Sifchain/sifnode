export enum ErrorCode {
  TX_FAILED_SLIPPAGE,
  TX_FAILED,
  USER_REJECTED,
  UNKNOWN_FAILURE,
  INSUFFICIENT_FUNDS,
  TX_FAILED_OUT_OF_GAS
}

const ErrorMessages = {
  [ErrorCode.TX_FAILED_SLIPPAGE]:
    "Your transaction has failed - Received amount is below expected",
  [ErrorCode.TX_FAILED]: "Your transaction has failed",
  [ErrorCode.USER_REJECTED]: "You have rejected the transaction",
  [ErrorCode.UNKNOWN_FAILURE]: "There was an unknown failure",
  [ErrorCode.INSUFFICIENT_FUNDS]: "You have insufficient funds",
  [ErrorCode.TX_FAILED_OUT_OF_GAS]: "Your transaction has failed - Out of gas",
};

export function getErrorMessage(code: ErrorCode): string {
  return ErrorMessages[code];
}
