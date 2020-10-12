import { Context } from ".";
declare function renderDestroyPool(isAdmin: boolean): {
    destroyPoolButtonAvailable: boolean;
};
declare const _default: ({ api, store }: Context) => {
    renderDestroyPool: typeof renderDestroyPool;
    destroyPool(): Promise<void>;
};
export default _default;
