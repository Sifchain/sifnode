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

export type EventHandler = (event: AppEvent) => void;
export type AppEventType = AppEvent["type"];
export type AppEventTypes = AppEventType[];

// TODO: 1. Surface EventEmitter DONE
// TODO: 2. Possibly type the events DONE
// TODO: 3. Create view layer component to present events from surfaced emitter DONE
// TODO: 4. Create view layer GA listener to surface events

export type NotificationServiceContext = {};
export default function createNotificationsService({}: NotificationServiceContext) {
  const emitter = new EventEmitter2();

  return {
    /**
     * Listen to all events
     * @param handler
     */
    onAny(handler: EventHandler) {
      emitter.onAny((_, value: AppEvent) => handler(value));
    },

    /**
     * Listen for a specific event or a list of specific events
     * @param eventType string or array of strings list of eventTypes
     * @param handler accepts an Event
     */
    on(eventType: AppEventType | AppEventTypes, handler: EventHandler) {
      const types = !Array.isArray(eventType) ? [eventType] : eventType;
      types.forEach(type => {
        emitter.on(type, handler);
      });
    },

    /**
     * Emit an event
     * @param event the AppEvent to emit
     */
    notify(event: AppEvent) {
      emitter.emit(event.type, event);
    },
  };
}
