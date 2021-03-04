import { Notification } from "../../entities";
import { notifications } from "../../store/notifications";

export type NotificationServiceContext = {};
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
