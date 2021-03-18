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
  const exponent =
    -1 *
    calculateExponent(
      options.mode === "percent",
      options.exponent,
      (amount as IAssetAmount).decimals,
    );

  const mantissa = extractMantissa(significand, exponent);

  const characteristic =
    exponent !== 0 ? significand.slice(0, exponent) : significand;
  const adjusted = [characteristic, mantissa].join(".");

  return numbro(adjusted).format(numbroConfig);
}

function extractMantissa(significand: string, exponent: number) {
  if (exponent !== 0) {
    const sliced = significand.slice(exponent);
    const diff = -1 * exponent;
    const padded = sliced.padStart(diff, "0"); // TODO: ES2017 do we need to polyfill?
    return padded;
  }
  return "";
}

function calculateExponent(
  isPercent: boolean,
  optionsExponent: number | undefined,
  amountDecimals: number | undefined,
): number {
  if (isPercent) {
    return (optionsExponent ?? 2) + 2;
  }
  return optionsExponent ?? amountDecimals ?? 0;
}

function isShorthandWithTotalLength(
  val: any,
): val is IFormatOptionsShorthandTotalLength {
  return val?.shorthand && val?.totalLength;
}

function createNumbroConfig(options: IFormatOptions) {
  const {
    forceSign = false,
    mode = "number",
    separator = false,
    space = false,
    prefix = "",
    postfix = "",
  } = options;

  if (isShorthandWithTotalLength(options)) {
    const { shorthand = false, totalLength = undefined } = options;
    return {
      forceSign: !!forceSign,
      output: mode,
      thousandSeparated: separator,
      spaceSeparated: space,
      average: shorthand,
      prefix,
      postfix,
      totalLength,
    };
  }
  const { shorthand = false, mantissa = 0, trimMantissa = false } = options;
  return {
    forceSign: !!forceSign,
    output: mode,
    thousandSeparated: separator,
    spaceSeparated: space,
    average: shorthand,
    prefix,
    postfix,
    mantissa,
    trimMantissa,
  };
}
