import { Ref } from "@vue/reactivity";
import { IAssetAmount, IPool } from "../../../entities";
export declare function assetPriceMessage(amount: IAssetAmount | null, pair: IPool | null, decimals: number): string;
export declare function trimZeros(amount: string): string;
export declare function useBalances(balances: Ref<IAssetAmount[]>): import("@vue/reactivity").ComputedRef<Map<string, IAssetAmount>>;
