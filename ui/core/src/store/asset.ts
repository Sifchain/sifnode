import { reactive } from "@vue/reactivity";
import { Asset } from "../entities";

export type AssetStore = {
  topTokens: Asset[];
};

export const asset = reactive({
  topTokens: [],
}) as AssetStore;
