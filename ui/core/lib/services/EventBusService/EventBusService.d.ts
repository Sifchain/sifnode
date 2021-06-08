import { AppEvent } from "./Events";
export declare type EventHandler = (event: AppEvent) => void;
export declare type AppEventType = AppEvent["type"];
export declare type AppEventTypes = AppEventType[];
export declare type EventBusServiceContext = {};
export default function createEventBusService({}: EventBusServiceContext): {
    /**
     * Listen to all events
     * @param handler
     */
    onAny(handler: EventHandler): void;
    /**
     * Listen for a specific event or a list of specific events
     * @param eventType string or array of strings list of eventTypes
     * @param handler accepts an Event
     */
    on(eventType: AppEventType | AppEventTypes, handler: EventHandler): void;
    /**
     * Emit an event
     * @param event the AppEvent to emit
     */
    dispatch(event: AppEvent): void;
};
