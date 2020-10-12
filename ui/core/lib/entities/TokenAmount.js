import { AssetAmount } from './AssetAmount';
export class TokenAmount extends AssetAmount {
    constructor(asset, amount) {
        super(asset, amount);
        this.asset = asset;
        this.amount = amount;
    }
}
export function createTokenAmount(amount, token) {
    return new TokenAmount(token, amount);
}
//# sourceMappingURL=TokenAmount.js.map