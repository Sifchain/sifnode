import { Coin } from "./Coin";
import { Token } from "./Token";

export type Asset = Token | Coin;

const assetMap = new Map<string, Asset>();

export const Asset = {
  set(key: string, value: Asset) {
    assetMap.set(key, value);
  },
  get(key: string): Asset | undefined {
    return assetMap.get(key);
  },
};
