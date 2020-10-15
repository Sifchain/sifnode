import { reactive } from "@vue/reactivity";
import { Asset } from "src/entities";

export type AssetStore = {
  assetMap: Map<string, Asset>; // to look up assets based on symbol
};

export const asset = reactive({
  assetMap: new Map(),
}) as AssetStore;
