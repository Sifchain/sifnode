// messageStore
import { reactive } from "@vue/reactivity";

// see get set for asset, similar
// import { Asset } from "../entities";

// entity:
type Notification = {
  id: String
  type: "error" | "success" | "info"
  notification: String
}

export type NotificationStore = {
  NotificationMap: Map<string, Notification>;
};

export const notifications = reactive({
  NotificationMap: new Map()
}) as NotificationStore;