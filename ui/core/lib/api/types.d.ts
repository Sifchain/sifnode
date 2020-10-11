import * as io from "./index";
export declare type FullApi = typeof io;
export declare type Api<T extends keyof FullApi = keyof FullApi, U extends object = {}> = {
    api: Pick<FullApi, T>;
} & U;
