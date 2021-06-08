import { Ref } from "@vue/reactivity";
import { Asset, IAssetAmount, LiquidityProvider, Pool } from "../../../entities";
export declare enum PoolState {
    SELECT_TOKENS = 0,
    ZERO_AMOUNTS = 1,
    INSUFFICIENT_FUNDS = 2,
    VALID_INPUT = 3,
    NO_LIQUIDITY = 4,
    ZERO_AMOUNTS_NEW_POOL = 5
}
export declare function usePoolCalculator(input: {
    tokenAAmount: Ref<string>;
    tokenASymbol: Ref<string | null>;
    tokenBAmount: Ref<string>;
    tokenBSymbol: Ref<string | null>;
    balances: Ref<IAssetAmount[]>;
    liquidityProvider: Ref<LiquidityProvider | null>;
    poolFinder: (a: Asset | string, b: Asset | string) => Ref<Pool> | null;
    asyncPooling: Ref<boolean>;
    lastFocusedTokenField: Ref<"A" | "B" | null>;
}): {
    state: import("@vue/reactivity").ComputedRef<PoolState.SELECT_TOKENS | PoolState.ZERO_AMOUNTS | PoolState.INSUFFICIENT_FUNDS | PoolState.VALID_INPUT | PoolState.ZERO_AMOUNTS_NEW_POOL>;
    aPerBRatioMessage: import("@vue/reactivity").ComputedRef<string>;
    bPerARatioMessage: import("@vue/reactivity").ComputedRef<string>;
    aPerBRatioProjectedMessage: import("@vue/reactivity").ComputedRef<string>;
    bPerARatioProjectedMessage: import("@vue/reactivity").ComputedRef<string>;
    shareOfPool: import("@vue/reactivity").ComputedRef<import("../../../entities").IAmount>;
    shareOfPoolPercent: import("@vue/reactivity").ComputedRef<string>;
    preExistingPool: import("@vue/reactivity").ComputedRef<{
        amounts: [IAssetAmount, IAssetAmount];
        otherAsset: (asset: import("../../../entities").IAsset) => IAssetAmount;
        symbol: () => string;
        contains: (...assets: import("../../../entities").IAsset[]) => boolean;
        toString: () => string;
        getAmount: (asset: string | import("../../../entities").IAsset) => IAssetAmount;
        poolUnits: import("../../../entities").IAmount;
        priceAsset(asset: import("../../../entities").IAsset): IAssetAmount;
        calcProviderFee(x: IAssetAmount): IAssetAmount;
        calcPriceImpact(x: IAssetAmount): import("../../../entities").IAmount;
        calcSwapResult(x: IAssetAmount): IAssetAmount;
        calcReverseSwapResult(Sa: IAssetAmount): IAssetAmount;
        calculatePoolUnits(nativeAssetAmount: IAssetAmount, externalAssetAmount: IAssetAmount): import("../../../entities").IAmount[];
    } | null>;
    totalLiquidityProviderUnits: import("@vue/reactivity").ComputedRef<string>;
    totalPoolUnits: import("@vue/reactivity").ComputedRef<string>;
    poolAmounts: import("@vue/reactivity").ComputedRef<IAssetAmount[] | null>;
    tokenAFieldAmount: import("@vue/reactivity").ComputedRef<IAssetAmount | null>;
    tokenBFieldAmount: import("@vue/reactivity").ComputedRef<IAssetAmount | null>;
};
