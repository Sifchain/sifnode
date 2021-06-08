import { Ref } from "@vue/reactivity";
import { Asset, LiquidityProvider, Pool } from "../entities";
import { PoolState } from "./addLiquidityCalculator";
export declare function useRemoveLiquidityCalculator(input: {
    externalAssetSymbol: Ref<string | null>;
    nativeAssetSymbol: Ref<string | null>;
    wBasisPoints: Ref<string | null>;
    asymmetry: Ref<string | null>;
    poolFinder: (a: Asset | string, b: Asset | string) => Ref<Pool> | null;
    liquidityProvider: Ref<LiquidityProvider | null>;
    sifAddress: Ref<string>;
}): {
    withdrawExternalAssetAmount: string;
    withdrawNativeAssetAmount: string;
    state: PoolState;
};
