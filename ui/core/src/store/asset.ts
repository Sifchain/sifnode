import { reactive } from "@vue/reactivity";
import { Asset } from "../entities";

export type AssetStore = {
  assetMap: Map<string, Asset>; // to look up assets based on symbol
  topTokens: Asset[];
};

export const asset = reactive({
  assetMap: new Map(),
  topTokens: [],
}) as AssetStore;
