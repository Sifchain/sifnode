// messageStore
import { reactive } from "@vue/reactivity";

// see get set for asset, similar
import { Notification } from "../entities";

// entity:
export type NotificationsStore = Array<Notification>;

export const notifications = reactive([]) as NotificationsStore;
