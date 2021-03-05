import { TransactionStatus } from "../../entities";
import { EventEmitter2 } from "eventemitter2";

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
  | PegTransactionPendingEvent
  | PegTransactionCompletedEvent
  | NoLiquidityPoolsFoundEvent
  | TransactionErrorEvent
  | PegTransactionErrorEvent;

export type NotificationServiceContext = {};

export type EventHandler = (event: AppEvent) => void;

// TODO: 1. Surface EventEmitter
// TODO: 2. Create view layer component to present events from surfaced emitter
// TODO: 3. Create view layer GA listener to surface events
// TODO: 4. Possibly type the events
export default function createNotificationsService({}: NotificationServiceContext) {
  const emitter = new EventEmitter2();
  return {
    on(eventType: AppEvent["type"], handler: EventHandler) {
      emitter.on(eventType, handler);
    },
    notify({ type, payload }: AppEvent) {
      emitter.emit(type, payload);
      // if (!type)
      //   throw 'Notification type required: "error", "success", "inform"';
      // if (!message) throw "Message string required";
      // notifications.unshift({ type, message, detail, loader });
      // return true;
    },
  };
}
