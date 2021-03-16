export enum ErrorCode {
  TX_FAILED_SLIPPAGE,
  TX_FAILED,
  USER_REJECTED,
  UNKNOWN_FAILURE,
  INSUFFICIENT_FUNDS,
}

const ErrorMessages = {
  [ErrorCode.TX_FAILED_SLIPPAGE]:
    "Your transaction has failed - Received amount is below expected",
  [ErrorCode.TX_FAILED]: "Your transaction has failed",
  [ErrorCode.USER_REJECTED]: "You have rejected the transaction",
  [ErrorCode.INSUFFICIENT_FUNDS]: "You have insufficient funds",
  [ErrorCode.UNKNOWN_FAILURE]: "There was an unknown failure",
};

export function getErrorMessage(code: ErrorCode): string {
  return ErrorMessages[code];
}
