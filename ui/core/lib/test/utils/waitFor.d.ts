declare type ToStringable = {
    toString: () => string;
};
export declare function waitFor(getter: () => Promise<ToStringable>, expected: ToStringable, name: string): Promise<void>;
export {};
