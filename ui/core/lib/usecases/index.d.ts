import { Api, FullApi } from "../api/types";
import { State, store } from "../store";
export declare type Context<T extends keyof FullApi = keyof FullApi> = Api<T, {
    state: State;
    store: typeof store;
}>;
export declare function createUsecases(context: Context): {
    swapTokens(token0Quantity: import("../entities").TokenAmount, token0: import("../entities").Token, token1: import("../entities").Token): Promise<void>;
    setQuantityOfToken(tokenAmount: import("../entities").TokenAmount): Promise<void>;
    updateListOfAvailableTokens(): Promise<void>;
    intializeRemoveLiquidity(): void;
    renderRemoveLiquidityPageData: (liquidityPool: import("../entities").Pair, token: import("../entities").Token, tokenAmount: import("../entities").TokenAmount) => {
        canRemoveLiquidity: boolean;
        amount: import("../entities").TokenAmount;
        gasFees: import("../entities").TokenAmount;
        shareOfPool: number;
        amountToRemoveIsTooHigh: boolean;
    };
    removeLiquidity(liquidityPool: import("../entities").Pair, token: import("../entities").Token, tokenAmount: import("../entities").TokenAmount): Promise<void>;
    renderDestroyPool: (isAdmin: boolean) => {
        destroyPoolButtonAvailable: boolean;
    };
    destroyPool(): Promise<void>;
    intializeCreatePoolUseCase(): void;
    renderCreatePoolData: (amountA: import("../entities").TokenAmount, amountB: import("../entities").TokenAmount) => {
        tokenAPerBRatio: number;
        tokenBPerARatio: number;
        tokenAAmountOwned: import("../entities").TokenAmount;
        tokenBAmountOwned: import("../entities").TokenAmount;
        shareOfPool: number;
        canCreatePool: boolean;
        isInsufficientFunds: boolean;
    };
    createPool(amountA: import("../entities").TokenAmount, amountB: import("../entities").TokenAmount): Promise<void>;
    connectToEthWallet(ethWallet: any): Promise<void>;
    connectToCosmosWallet(cosmosWallet: any): Promise<void>;
    broadcastTx(tx: any): Promise<void>;
    intializeAddLiquidityUseCase(): void;
    renderLiquidityData: (amountA: import("../entities").TokenAmount, amountB?: import("../entities").TokenAmount) => {
        tokenAPerBRatio: number;
        tokenBPerARatio: number;
        tokenAAmountOwned: import("../entities").TokenAmount;
        tokenBAmountOwned: import("../entities").TokenAmount;
        shareOfPool: number;
        isInsufficientFunds: boolean;
    };
    addLiquidity(amountA: import("../entities").TokenAmount, amountB: import("../entities").TokenAmount): Promise<void>;
};
export declare const usecases: {
    swapTokens(token0Quantity: import("../entities").TokenAmount, token0: import("../entities").Token, token1: import("../entities").Token): Promise<void>;
    setQuantityOfToken(tokenAmount: import("../entities").TokenAmount): Promise<void>;
    updateListOfAvailableTokens(): Promise<void>;
    intializeRemoveLiquidity(): void;
    renderRemoveLiquidityPageData: (liquidityPool: import("../entities").Pair, token: import("../entities").Token, tokenAmount: import("../entities").TokenAmount) => {
        canRemoveLiquidity: boolean;
        amount: import("../entities").TokenAmount;
        gasFees: import("../entities").TokenAmount;
        shareOfPool: number;
        amountToRemoveIsTooHigh: boolean;
    };
    removeLiquidity(liquidityPool: import("../entities").Pair, token: import("../entities").Token, tokenAmount: import("../entities").TokenAmount): Promise<void>;
    renderDestroyPool: (isAdmin: boolean) => {
        destroyPoolButtonAvailable: boolean;
    };
    destroyPool(): Promise<void>;
    intializeCreatePoolUseCase(): void;
    renderCreatePoolData: (amountA: import("../entities").TokenAmount, amountB: import("../entities").TokenAmount) => {
        tokenAPerBRatio: number;
        tokenBPerARatio: number;
        tokenAAmountOwned: import("../entities").TokenAmount;
        tokenBAmountOwned: import("../entities").TokenAmount;
        shareOfPool: number;
        canCreatePool: boolean;
        isInsufficientFunds: boolean;
    };
    createPool(amountA: import("../entities").TokenAmount, amountB: import("../entities").TokenAmount): Promise<void>;
    connectToEthWallet(ethWallet: any): Promise<void>;
    connectToCosmosWallet(cosmosWallet: any): Promise<void>;
    broadcastTx(tx: any): Promise<void>;
    intializeAddLiquidityUseCase(): void;
    renderLiquidityData: (amountA: import("../entities").TokenAmount, amountB?: import("../entities").TokenAmount) => {
        tokenAPerBRatio: number;
        tokenBPerARatio: number;
        tokenAAmountOwned: import("../entities").TokenAmount;
        tokenBAmountOwned: import("../entities").TokenAmount;
        shareOfPool: number;
        isInsufficientFunds: boolean;
    };
    addLiquidity(amountA: import("../entities").TokenAmount, amountB: import("../entities").TokenAmount): Promise<void>;
};
export declare type UseCases = ReturnType<typeof createUsecases>;
