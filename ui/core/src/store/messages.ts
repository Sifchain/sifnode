// messageStore
import { reactive } from "@vue/reactivity";

// see get set for asset, similar
// import { Asset } from "../entities";

// entity:
type Message = {
  id: String
  type: "error" | "success" | "info"
  message: String
}

export type MessageStore = {
  messageMap: Map<string, Message>;
};

export const messages = reactive({
  messageMap: new Map()
}) as MessageStore;
