import { Ref } from "@vue/reactivity";
import { Store } from ".";
import { Asset, Pool } from "../entities";
declare type PoolFinderFn = (s: Store) => (a: Asset | string, b: Asset | string) => Ref<Pool> | null;
export declare const createPoolFinder: PoolFinderFn;
export {};
