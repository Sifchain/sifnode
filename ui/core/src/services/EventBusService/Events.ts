import { TransactionStatus } from "../../entities";

// Add more wallet types here as they come up
type WalletType = "sif" | "eth";

type ErrorEvent = {
  type: "ErrorEvent";
  payload: {
    message: string;
    detail?: {
      type: "etherscan" | "info";
      message: string;
    };
  };
};

type TransactionErrorEvent = {
  type: "TransactionErrorEvent";
  payload: {
    txStatus: TransactionStatus;
    message: string;
  };
};

type WalletConnectedEvent = {
  type: "WalletConnectedEvent";
  payload: { walletType: WalletType; address: string };
};

type WalletDisconnectedEvent = {
  type: "WalletDisconnectedEvent";
  payload: { walletType: WalletType; address: string };
};

type WalletConnectionErrorEvent = {
  type: "WalletConnectionErrorEvent";
  payload: { walletType: WalletType; message: string };
};

type PegTransactionPendingEvent = {
  type: "PegTransactionPendingEvent";
  payload: { hash: string };
};

type PegTransactionCompletedEvent = {
  type: "PegTransactionCompletedEvent";
  payload: {
    hash: string;
  };
};

type PegTransactionErrorEvent = {
  type: "PegTransactionErrorEvent";
  payload: {
    txStatus: TransactionStatus;
    message: string;
  };
};

type NoLiquidityPoolsFoundEvent = {
  type: "NoLiquidityPoolsFoundEvent";
  payload: {};
};

export type AppEvent =
  | ErrorEvent
  | WalletConnectedEvent
  | WalletDisconnectedEvent
  | WalletConnectionErrorEvent
  | PegTransactionPendingEvent
  | PegTransactionCompletedEvent
  | NoLiquidityPoolsFoundEvent
  | TransactionErrorEvent
  | PegTransactionErrorEvent;
