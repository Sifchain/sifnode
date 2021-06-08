export declare type TendermintSocketSubscriber = ReturnType<typeof TendermintSocketSubscriber>;
export declare function TendermintSocketSubscriber({ wsUrl }: {
    wsUrl: string;
}): {
    on(event: "Tx" | "NewBlock" | "error", handler: (event: any) => void): void;
    off(event: "Tx" | "NewBlock" | "error", handler: (event: any) => void): void;
};
export declare function createTendermintSocketSubscriber(wsUrl: string): {
    on(event: "error" | "Tx" | "NewBlock", handler: (event: any) => void): void;
    off(event: "error" | "Tx" | "NewBlock", handler: (event: any) => void): void;
};
