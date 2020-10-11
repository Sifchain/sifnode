import { Context } from ".";
import { Token, TokenAmount } from "../entities";
declare const _default: ({ api, store }: Context) => {
    swapTokens(token0Quantity: TokenAmount, token0: Token, token1: Token): Promise<void>;
};
export default _default;
