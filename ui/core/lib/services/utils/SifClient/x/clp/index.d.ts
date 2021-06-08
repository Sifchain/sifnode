import { LcdClient, Msg } from "@cosmjs/launchpad";
export declare type SwapParams = {
    sent_asset: {
        symbol: string;
        ticker: string;
        source_chain: string;
    };
    received_asset: {
        symbol: string;
        ticker: string;
        source_chain: string;
    };
    base_req: {
        from: string;
        chain_id: string;
    };
    signer: string;
    sent_amount: string;
    min_receiving_amount: string;
};
export declare type LiquidityParams = {
    base_req: {
        from: string;
        chain_id: string;
    };
    external_asset: {
        source_chain: string;
        symbol: string;
        ticker: string;
    };
    native_asset_amount: string;
    external_asset_amount: string;
    signer: string;
};
export declare type RemoveLiquidityParams = {
    base_req: {
        from: string;
        chain_id: string;
    };
    external_asset: {
        source_chain: string;
        symbol: string;
        ticker: string;
    };
    w_basis_points: string;
    asymmetry: string;
    signer: string;
};
export declare type RawPool = {
    external_asset: {
        source_chain: string;
        symbol: string;
        ticker: string;
    };
    native_asset_balance: string;
    external_asset_balance: string;
    pool_units: string;
    pool_address: string;
};
declare type LiquidityDetailsResponse = {
    result: {
        external_asset_balance: string;
        native_asset_balance: string;
        LiquidityProvider: {
            liquidity_provider_units: string;
            liquidity_provider_address: string;
            asset: {
                symbol: string;
                ticker: string;
                source_chain: string;
            };
        };
    };
    height: string;
};
declare type ClpCmdSwap = (params: SwapParams) => Promise<Msg>;
declare type ClpQueryPools = () => Promise<RawPool[]>;
declare type ClpQueryPool = (params: {
    ticker: string;
}) => Promise<RawPool>;
declare type ClpQueryAssets = (address: string) => Promise<{
    symbol: string;
}[]>;
declare type ClpAddLiquidity = (params: LiquidityParams) => Promise<any>;
declare type ClpCreatePool = (params: LiquidityParams) => Promise<any>;
declare type ClpGetLiquidityProvider = (params: {
    symbol: string;
    lpAddress: string;
}) => Promise<LiquidityDetailsResponse>;
declare type ClpRemoveLiquidity = (param: RemoveLiquidityParams) => Promise<any>;
export interface ClpExtension {
    readonly clp: {
        swap: ClpCmdSwap;
        getPools: ClpQueryPools;
        getAssets: ClpQueryAssets;
        addLiquidity: ClpAddLiquidity;
        createPool: ClpCreatePool;
        getLiquidityProvider: ClpGetLiquidityProvider;
        removeLiquidity: ClpRemoveLiquidity;
        getPool: ClpQueryPool;
    };
}
export declare function setupClpExtension(base: LcdClient): ClpExtension;
export {};
