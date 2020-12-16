// messageStore
import { reactive } from "@vue/reactivity";

// see get set for asset, similar
import { Notification } from "../entities";

// entity:
export type INotificationsStore = Array<Notification>

export const notifications = reactive([]) as INotificationsStore;