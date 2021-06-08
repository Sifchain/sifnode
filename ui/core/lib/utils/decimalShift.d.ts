import { Asset, IAmount } from "../entities";
import { AmountNotAssetAmount } from "./format";
/**
 * Function to shift the magnitude of a string without using any Math libs
 * This helps us move between decimals and integers
 * @param decimal the decimal string
 * @param shift the shift in the decimal point required
 * @returns string decimal
 */
export declare function decimalShift(decimal: string, shift: number): string;
/**
 * Utility for converting to the base units of an asset
 * @param decimal the decimal string
 * @param asset the asset we want to get the base unit amount for
 * @returns amount as a string
 */
export declare function toBaseUnits(decimal: string, asset: Asset): string;
/**
 * Utility for converting from the base units of an asset
 * @param integer the integer string
 * @param asset the asset we want to get the base unit amount for
 * @returns amount as a string
 */
export declare function fromBaseUnits(integer: string, asset: Asset): string;
/**
 * Remove the decimal component of a string representation of a decimal number
 * @param decimal decimal to floor
 * @returns string with everything before the decimal point
 */
export declare function floorDecimal(decimal: string): string;
/**
 * Utility to get the length of the trimmed mantissa from the amount
 * @param amount an IAmount
 * @returns length of mantissa
 */
export declare function getMantissaLength<T extends IAmount>(amount: AmountNotAssetAmount<T>): number;
