"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
const eventemitter2_1 = require("eventemitter2");
function createEventBusService({}) {
    const emitter = new eventemitter2_1.EventEmitter2();
    return {
        /**
         * Listen to all events
         * @param handler
         */
        onAny(handler) {
            emitter.onAny((_, value) => handler(value));
        },
        /**
         * Listen for a specific event or a list of specific events
         * @param eventType string or array of strings list of eventTypes
         * @param handler accepts an Event
         */
        on(eventType, handler) {
            const types = !Array.isArray(eventType) ? [eventType] : eventType;
            types.forEach((type) => {
                emitter.on(type, handler);
            });
        },
        /**
         * Emit an event
         * @param event the AppEvent to emit
         */
        dispatch(event) {
            emitter.emit(event.type, event);
        },
    };
}
exports.default = createEventBusService;
//# sourceMappingURL=EventBusService.js.map