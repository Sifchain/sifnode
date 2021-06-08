import { Ref } from "@vue/reactivity";
import { Asset, IPool, IAssetAmount } from "../../../entities";
export declare enum SwapState {
    SELECT_TOKENS = 0,
    ZERO_AMOUNTS = 1,
    INSUFFICIENT_FUNDS = 2,
    VALID_INPUT = 3,
    INVALID_AMOUNT = 4,
    INSUFFICIENT_LIQUIDITY = 5
}
export declare function useSwapCalculator(input: {
    fromAmount: Ref<string>;
    fromSymbol: Ref<string | null>;
    toAmount: Ref<string>;
    toSymbol: Ref<string | null>;
    balances: Ref<IAssetAmount[]>;
    selectedField: Ref<"from" | "to" | null>;
    slippage: Ref<string>;
    poolFinder: (a: Asset | string, b: Asset | string) => Ref<IPool> | null;
}): {
    priceMessage: import("@vue/reactivity").ComputedRef<string>;
    state: import("@vue/reactivity").ComputedRef<SwapState>;
    fromFieldAmount: import("@vue/reactivity").ComputedRef<IAssetAmount | null>;
    toFieldAmount: import("@vue/reactivity").ComputedRef<IAssetAmount | null>;
    toAmount: Ref<string>;
    fromAmount: Ref<string>;
    priceImpact: import("@vue/reactivity").ComputedRef<string | null>;
    providerFee: import("@vue/reactivity").ComputedRef<string | null>;
    minimumReceived: import("@vue/reactivity").ComputedRef<IAssetAmount | null>;
    swapResult: Ref<{
        readonly address?: string | undefined;
        readonly decimals: number;
        readonly imageUrl?: string | undefined;
        readonly name: string;
        readonly network: import("../../../entities").Network;
        readonly symbol: string;
        readonly label: string;
        readonly asset: {
            address?: string | undefined;
            decimals: number;
            imageUrl?: string | undefined;
            name: string;
            network: import("../../../entities").Network;
            symbol: string;
            label: string;
        };
        readonly amount: {
            toBigInt: () => import("jsbi").default;
            toString: (detailed?: boolean | undefined) => string;
            add: (other: string | import("../../../entities").IAmount) => import("../../../entities").IAmount;
            divide: (other: string | import("../../../entities").IAmount) => import("../../../entities").IAmount;
            equalTo: (other: string | import("../../../entities").IAmount) => boolean;
            greaterThan: (other: string | import("../../../entities").IAmount) => boolean;
            greaterThanOrEqual: (other: string | import("../../../entities").IAmount) => boolean;
            lessThan: (other: string | import("../../../entities").IAmount) => boolean;
            lessThanOrEqual: (other: string | import("../../../entities").IAmount) => boolean;
            multiply: (other: string | import("../../../entities").IAmount) => import("../../../entities").IAmount;
            sqrt: () => import("../../../entities").IAmount;
            subtract: (other: string | import("../../../entities").IAmount) => import("../../../entities").IAmount;
        };
        toBigInt: () => import("jsbi").default;
        toString: (() => string) & ((detailed?: boolean | undefined) => string);
        toDerived: () => import("../../../entities").IAmount;
        add: (other: string | import("../../../entities").IAmount) => import("../../../entities").IAmount;
        divide: (other: string | import("../../../entities").IAmount) => import("../../../entities").IAmount;
        equalTo: (other: string | import("../../../entities").IAmount) => boolean;
        greaterThan: (other: string | import("../../../entities").IAmount) => boolean;
        greaterThanOrEqual: (other: string | import("../../../entities").IAmount) => boolean;
        lessThan: (other: string | import("../../../entities").IAmount) => boolean;
        lessThanOrEqual: (other: string | import("../../../entities").IAmount) => boolean;
        multiply: (other: string | import("../../../entities").IAmount) => import("../../../entities").IAmount;
        sqrt: () => import("../../../entities").IAmount;
        subtract: (other: string | import("../../../entities").IAmount) => import("../../../entities").IAmount;
    } | null>;
    reverseSwapResult: Ref<{
        readonly address?: string | undefined;
        readonly decimals: number;
        readonly imageUrl?: string | undefined;
        readonly name: string;
        readonly network: import("../../../entities").Network;
        readonly symbol: string;
        readonly label: string;
        readonly asset: {
            address?: string | undefined;
            decimals: number;
            imageUrl?: string | undefined;
            name: string;
            network: import("../../../entities").Network;
            symbol: string;
            label: string;
        };
        readonly amount: {
            toBigInt: () => import("jsbi").default;
            toString: (detailed?: boolean | undefined) => string;
            add: (other: string | import("../../../entities").IAmount) => import("../../../entities").IAmount;
            divide: (other: string | import("../../../entities").IAmount) => import("../../../entities").IAmount;
            equalTo: (other: string | import("../../../entities").IAmount) => boolean;
            greaterThan: (other: string | import("../../../entities").IAmount) => boolean;
            greaterThanOrEqual: (other: string | import("../../../entities").IAmount) => boolean;
            lessThan: (other: string | import("../../../entities").IAmount) => boolean;
            lessThanOrEqual: (other: string | import("../../../entities").IAmount) => boolean;
            multiply: (other: string | import("../../../entities").IAmount) => import("../../../entities").IAmount;
            sqrt: () => import("../../../entities").IAmount;
            subtract: (other: string | import("../../../entities").IAmount) => import("../../../entities").IAmount;
        };
        toBigInt: () => import("jsbi").default;
        toString: (() => string) & ((detailed?: boolean | undefined) => string);
        toDerived: () => import("../../../entities").IAmount;
        add: (other: string | import("../../../entities").IAmount) => import("../../../entities").IAmount;
        divide: (other: string | import("../../../entities").IAmount) => import("../../../entities").IAmount;
        equalTo: (other: string | import("../../../entities").IAmount) => boolean;
        greaterThan: (other: string | import("../../../entities").IAmount) => boolean;
        greaterThanOrEqual: (other: string | import("../../../entities").IAmount) => boolean;
        lessThan: (other: string | import("../../../entities").IAmount) => boolean;
        lessThanOrEqual: (other: string | import("../../../entities").IAmount) => boolean;
        multiply: (other: string | import("../../../entities").IAmount) => import("../../../entities").IAmount;
        sqrt: () => import("../../../entities").IAmount;
        subtract: (other: string | import("../../../entities").IAmount) => import("../../../entities").IAmount;
    } | null>;
};
