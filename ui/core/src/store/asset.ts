import { reactive } from "@vue/reactivity";
import { Asset } from "../entities";

export type AssetStore = {
  assetMap: Map<string, Asset>; // to look up assets based on symbol
  top20Tokens: Asset[];
};

export const asset = reactive({
  assetMap: new Map(),
  top20Tokens: [],
}) as AssetStore;
