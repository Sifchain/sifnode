import { Token } from './Token';
import JSBI from 'jsbi';
import { AssetAmount } from './AssetAmount';
import { BigintIsh } from './fraction/Fraction';
export declare class TokenAmount extends AssetAmount {
    asset: Token;
    amount: BigintIsh;
    constructor(asset: Token, amount: BigintIsh);
}
export declare function createTokenAmount(amount: JSBI, token: Token): TokenAmount;
