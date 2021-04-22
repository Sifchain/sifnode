import { TransactionStatus } from "../../entities";
import { EventEmitter2 } from "eventemitter2";
import { AppEvent } from "./Events";

export type EventHandler = (event: AppEvent) => void;
export type AppEventType = AppEvent["type"];
export type AppEventTypes = AppEventType[];

export type EventBusServiceContext = {};

export default function createEventBusService({}: EventBusServiceContext) {
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
      types.forEach((type) => {
        emitter.on(type, handler);
      });
    },

    /**
     * Emit an event
     * @param event the AppEvent to emit
     */
    dispatch(event: AppEvent) {
      emitter.emit(event.type, event);
    },
  };
}
