import { Ref, toRefs } from "@vue/reactivity";
import { Store } from ".";
import { Asset, Pool } from "../entities";

type PoolFinderFn = (
  s: Store,
) => (a: Asset | string, b: Asset | string) => Ref<Pool> | null;

export const createPoolFinder: PoolFinderFn = (s: Store) => (
  a: Asset | string, // externalAsset
  b: Asset | string, // nativeAsset
) => {
  const pools = toRefs(s.pools);
  const key = [a, b]
    .map(x => (typeof x === "string" ? x : x.symbol))
    .join("_") as keyof typeof pools;

  const poolRef = pools[key] as Ref<Pool> | undefined;
  return poolRef ?? null;
};
