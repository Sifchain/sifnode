import { Notification } from "../../entities";
import { notifications } from "../../store/notifications";

type INotificationService = Boolean;

/**
 *  Notification service
 *
 * Handles communication between core and in-app notifications related to success, error, and information.
 * Mostly from the API level
 */
export default function notify({
  type,
  message,
  detail,
}: Notification): INotificationService {
  if (!type) throw 'Notification type required: "error", "success", "inform"';
  if (!message) throw "Message string required";

  // add to first place in notification state
  notifications.unshift({ type, message, detail });

  return true;
}
