import invariant from "tiny-invariant";
import _Big from "big.js";
import toFormat from "toformat";
import { Fraction, parseBigintIsh, Rounding, TEN, } from "./fraction/Fraction";
import JSBI from "jsbi";
const Big = toFormat(_Big);
export class AssetAmount extends Fraction {
    constructor(asset, amount) {
        super(parseBigintIsh(amount), JSBI.exponentiate(TEN, JSBI.BigInt(asset.decimals)));
        this.asset = asset;
        this.amount = amount;
    }
    toSignificant(significantDigits = 6, format, rounding = Rounding.ROUND_DOWN) {
        return super.toSignificant(significantDigits, format, rounding);
    }
    toFixed(decimalPlaces = this.asset.decimals, format, rounding = Rounding.ROUND_DOWN) {
        invariant(decimalPlaces <= this.asset.decimals, "DECIMALS");
        return super.toFixed(decimalPlaces, format, rounding);
    }
    toExact(format = { groupSeparator: "" }) {
        Big.DP = this.asset.decimals;
        return new Big(this.numerator.toString())
            .div(this.denominator.toString())
            .toFormat(format);
    }
    static create(asset, amount) {
        return new AssetAmount(asset, amount);
    }
}
//# sourceMappingURL=AssetAmount.js.map