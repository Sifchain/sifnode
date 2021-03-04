import { Notification } from "../../entities";
import { notifications } from "../../store/notifications";

export type NotificationServiceContext = {};

// TODO: 1. Surface EventEmitter
// TODO: 2. Create view layer component to present events from surfaced emitter
// TODO: 3. Create view layer GA listener to surface events
// TODO: 4. Possibly type the events
export default function createNotificationsService({}: NotificationServiceContext) {
  return {
    notify({ type, message, detail, loader }: Notification) {
      if (!type)
        throw 'Notification type required: "error", "success", "inform"';
      if (!message) throw "Message string required";

      notifications.unshift({ type, message, detail, loader });

      return true;
    },
  };
}
