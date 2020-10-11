import { Context } from ".";
import { TokenAmount } from "../entities";
declare const _default: ({ api, store }: Context) => {
    setQuantityOfToken(tokenAmount: TokenAmount): Promise<void>;
};
export default _default;
