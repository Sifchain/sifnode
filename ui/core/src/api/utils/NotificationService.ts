import { reactive } from "@vue/reactivity";
import {
  Notifications,
  Notification
} from "../../entities";

import {notifications} from "../../store/notifications"
// i am pulling this from entitty, but id?
// type NotificationServiceContext = {
//   type: "error" | "success" | "inform"
//   message: String
// }

// I dont think this will return anything? success, fail? boolean?
type INotificationService = Boolean
/**
 * Constructor for Notification service
 *
 * Handles communication between core and in-app notifications related to success, error, and information. 
 */

 // I'm thinking of createNotificationService, but why? What do you call it where you wrap everything like that?
 // this doesn't need to be in createApi. So I put in utils

// could do notifyError, notifySuccess, but adding type allows for more uniform interface

// todo rename this file to notifyService? nah

// opts could include timeout remove,or  manual remove (rmNotify())
export default function notify({
    type,
    message,
    detail
  }: Notification
  ): INotificationService {

  if (!type) throw "Notification type required: \"error\", \"success\", \"inform\""
  if (!message) throw "Message string required"

  // Reactive state for communicating state changes
  // yes I think event emitter here would be smart
  // think i'll import state
  // const state: {
  //   notifications: Notifications
  // } = reactive({
  //   notifications: []
  // });

  // add to first place in notification state
  notifications.unshift({type, message, detail})

  return true
}
