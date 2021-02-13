import { Notification } from "../../entities";
import { notifications } from "../../store/notifications";

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
  loader
}: Notification): boolean {
  if (!type) throw 'Notification type required: "error", "success", "inform"';
  if (!message) throw "Message string required";

  // check if notification already exists
  if(notifications.find((item) => {
    return item.message === message
  })) { return false }

  notifications.unshift({ type, message, detail, loader });

  return true;
}
