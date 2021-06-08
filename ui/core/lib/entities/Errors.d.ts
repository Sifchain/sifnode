export declare enum ErrorCode {
    TX_FAILED_SLIPPAGE = 0,
    TX_FAILED = 1,
    USER_REJECTED = 2,
    UNKNOWN_FAILURE = 3,
    INSUFFICIENT_FUNDS = 4,
    TX_FAILED_OUT_OF_GAS = 5,
    TX_FAILED_NOT_ENOUGH_ROWAN_TO_COVER_GAS = 6,
    TX_FAILED_USER_NOT_ENOUGH_BALANCE = 7
}
export declare function getErrorMessage(code: ErrorCode): string;
