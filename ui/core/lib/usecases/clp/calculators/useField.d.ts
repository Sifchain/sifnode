import { Ref } from "@vue/reactivity";
export declare function useField(amount: Ref<string>, symbol: Ref<string | null>): {
    fieldAmount: import("@vue/reactivity").ComputedRef<import("../../../entities").IAssetAmount | null>;
    asset: import("@vue/reactivity").ComputedRef<Readonly<import("../../../entities").IAsset> | null>;
};
