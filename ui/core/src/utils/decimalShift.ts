import { Amount, Asset, IAmount, IAsset } from "../entities";
import { AmountNotAssetAmount, format, trimMantissa } from "./format";

/**
 * Function to shift the magnitude of a string without using any Math libs
 * This helps us move between decimals and integers
 * @param decimal the decimal string
 * @param shift the shift in the decimal point required
 * @returns string decimal
 */
export function decimalShift(decimal: string, shift: number) {
  if (!decimal.match(/^[+-]?(\d+)?\.?\d+$/)) {
    throw new Error(`Cannot recognise number format: ${decimal}`);
  }

  // Extract the sign(+-) and the unsigned content from the decimal
  const [, sign = "", unsignedDecimal = decimal] =
    decimal.match(/^([+-]?)(.+)$/) || [];

  // Eg, "12.34" > ['12','34']
  const [origCharacter, origMantissa] = unsignedDecimal.split(".");

  const dotIndex = origCharacter.length;
  const targetIndex = dotIndex + shift;

  // Significand is all the digits "1234"
  const significand = [origCharacter, origMantissa].join("");

  const character =
    targetIndex >= 0
      ? significand
          .slice(0, targetIndex)
          .padEnd(targetIndex, "0")
          .replace(/^0+/, "") || "0"
      : "0";

  const mantissa =
    targetIndex >= 0
      ? significand.slice(targetIndex)
      : significand.padStart(Math.abs(targetIndex) + significand.length, "0");

  return `${sign}${[character, mantissa].filter(Boolean).join(".")}`;
}

/**
 * Utility for converting to the base units of an asset
 * @param decimal the decimal string
 * @param asset the asset we want to get the base unit amount for
 * @returns amount as a string
 */
export function toBaseUnits(decimal: string, asset: Asset): string {
  return decimalShift(decimal, asset.decimals);
}

/**
 * Utility for converting from the base units of an asset
 * @param integer the integer string
 * @param asset the asset we want to get the base unit amount for
 * @returns amount as a string
 */
export function fromBaseUnits(integer: string, asset: Asset): string {
  return decimalShift(integer, -1 * asset.decimals);
}

/**
 * Remove the decimal component of a string representation of a decimal number
 * @param decimal decimal to floor
 * @returns string with everything before the decimal point
 */
export function floorDecimal(decimal: string) {
  return decimal.split(".")[0];
}

/**
 * Utility to get the length of the trimmed mantissa from the amount
 * @param amount an IAmount
 * @returns length of mantissa
 */
export function getMantissaLength<T extends IAmount>(
  amount: AmountNotAssetAmount<T>,
): number {
  const number = format(amount, { mantissa: 18, trimMantissa: true });
  return number.length - number.indexOf(".") - 1;
}
