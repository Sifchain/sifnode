export declare type TendermintSocketPoll = ReturnType<typeof TendermintSocketPoll>;
export declare type BlockData = {
    result: {
        block: {
            header: {
                height: string;
            };
            data: {
                txs: null | string[];
            };
        };
    };
};
declare function fetchBlock(url: string): Promise<BlockData>;
export declare type ITendermintSocketPoll = ReturnType<typeof TendermintSocketPoll>;
export declare function TendermintSocketPoll({ apiUrl, fetcher, pollInterval, }: {
    apiUrl: string;
    fetcher?: typeof fetchBlock;
    pollInterval?: number;
}): {
    on(event: "Tx" | "NewBlock" | "error", handler: (event: any) => void): void;
    off(event: "Tx" | "NewBlock" | "error", handler: (event: any) => void): void;
};
export declare function createTendermintSocketPoll(apiUrl: string): {
    on(event: "error" | "Tx" | "NewBlock", handler: (event: any) => void): void;
    off(event: "error" | "Tx" | "NewBlock", handler: (event: any) => void): void;
};
export {};
