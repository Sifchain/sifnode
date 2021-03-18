import { IAmount } from "../entities/Amount";
import { IAssetAmount } from "../entities/AssetAmount_";
import numbro from "numbro";

type IFormatOptionsBase = {
  exponent?: number; // display = (amount * 10^-exponent) when undefined exponent will be set by (amount as IAssetAmount).decimals ?? 0 - defaults to 2 for percent mode
  forceSign?: boolean; // Ensure we have a + sign at the start of the value default false
  mode?: "number" | "percent"; // defines the rendering strategy default "number"
  separator?: boolean; // Add thousand separators eg. 1,000 default false
  space?: boolean; // separate prefix and suffix with spaces default false
  prefix?: string; // Add a prefix
  postfix?: string; // Add a postfix
  zeroFormat?: string; // could be something like `N/A`
};

type IFormatOptionsMantissa = IFormatOptionsBase & {
  shorthand?: boolean;
  mantissa?: number; // number of decimals after point default is exponent
  trimMantissa?: boolean; // Remove 0s from the mantissa default false
};

type IFormatOptionsShorthandTotalLength = IFormatOptionsBase & {
  shorthand: true;
  totalLength?: number; // This will give us significant digits using abbreviations eg. `1.234k` it will override anything in mantissa
};

export type IFormatOptions =
  | IFormatOptionsMantissa
  | IFormatOptionsShorthandTotalLength;

export function format(
  amount: IAmount | IAssetAmount,
  options: IFormatOptionsMantissa,
): string;
export function format(
  amount: IAmount | IAssetAmount,
  options: IFormatOptionsShorthandTotalLength,
): string;
export function format(
  amount: IAmount | IAssetAmount,
  options: IFormatOptions,
): string {
  const numbroConfig = createNumbroConfig(options);
  const significand = amount.toBigInt().toString();

  const exponentShift = calculateExponentShift(
    options.mode === "percent",
    options.exponent,
    (amount as IAssetAmount).decimals,
  );

  const withPoint = exponentiateString(significand, exponentShift);

  return numbro(withPoint).format(numbroConfig);
}

export function exponentiateString(
  significand: string,
  exponentShift: number,
): string {
  const mantissa = extractMantissa(significand, exponentShift);
  const characteristic = extractCharacteristic(significand, exponentShift);
  return [characteristic, mantissa].join(".");
}

function extractMantissa(significand: string, exponentShift: number) {
  if (exponentShift !== 0) {
    const sliced = significand.slice(exponentShift);
    const diff = -1 * exponentShift;
    const padded = sliced.padStart(diff, "0"); // TODO: ES2017 do we need to polyfill?
    return padded;
  }
  return "";
}

function extractCharacteristic(significand: string, exponentShift: number) {
  return exponentShift !== 0
    ? significand.slice(0, exponentShift)
    : significand;
}

function calculateExponentShift(
  isPercent: boolean,
  optionsExponent: number | undefined,
  amountDecimals: number | undefined,
): number {
  if (isPercent) {
    return -1 * ((optionsExponent ?? 2) + 2);
  }
  return -1 * (optionsExponent ?? amountDecimals ?? 0);
}

function isShorthandWithTotalLength(
  val: any,
): val is IFormatOptionsShorthandTotalLength {
  return val?.shorthand && val?.totalLength;
}

function createNumbroConfig(options: IFormatOptions) {
  return {
    forceSign: options.forceSign ?? false,
    output: options.mode ?? "number",
    thousandSeparated: options.separator ?? false,
    spaceSeparated: options.space ?? false,
    prefix: options.prefix ?? "",
    postfix: options.postfix ?? "",
    ...(isShorthandWithTotalLength(options)
      ? {
          average: options.shorthand ?? false,
          totalLength: options.totalLength,
        }
      : {
          average: options.shorthand ?? false,
          mantissa: options.mantissa ?? 0,
          trimMantissa: options.trimMantissa ?? false,
        }),
  };
}
