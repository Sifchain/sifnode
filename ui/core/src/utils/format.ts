import { Amount, IAmount } from "../entities/Amount";
import { IAssetAmount } from "../entities/AssetAmount";
import numbro from "numbro";
import { IAsset } from "../entities";
import { decimalShift } from "./decimalShift";

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

function isAsset(val: any): val is IAsset {
  return typeof val?.symbol === "string";
}

export function round(decimal: string, places: number) {
  return decimalShift(
    Amount(decimal)
      .multiply(Amount(decimalShift("1", places)))
      .toBigInt() // apply rounding
      .toString(),
    -1 * places,
  );
}

type AmountNotAssetAmount<T extends IAmount> = T extends IAssetAmount
  ? never
  : T;

export function format<T extends IAmount>(
  amount: AmountNotAssetAmount<T>,
  options: IFormatOptions,
): string;
export function format<T extends IAmount>(
  amount: AmountNotAssetAmount<T>,
  asset: Exclude<IAsset, IAssetAmount>,
  options: IFormatOptions,
): string;
export function format<T extends IAmount>(
  _amount: AmountNotAssetAmount<T>,
  _asset: Exclude<IAsset, IAssetAmount> | IFormatOptions,
  _options?: IFormatOptions,
): string {
  const amount = _amount;
  const options = isAsset(_asset) ? _options! : _asset;
  const asset = isAsset(_asset) ? _asset : undefined;

  let decimal = asset
    ? decimalShift(amount.toBigInt().toString(), -1 * asset.decimals)
    : amount.toString();

  let postfix = options.prefix ?? "";
  let prefix = options.postfix ?? "";
  let space = "";

  if (options.zeroFormat && amount.equalTo("0")) {
    return options.zeroFormat;
  }

  if (options.shorthand) {
    return numbro(decimal).format(createNumbroConfig(options));
  }

  if (options.space) {
    space = " ";
  }

  if (options.mode === "percent") {
    decimal = decimalShift(decimal, 2);
    postfix = "%";
  }

  if (typeof options.mantissa === "number") {
    decimal = applyMantissa(decimal, options.mantissa);
  }

  if (options.trimMantissa) {
    decimal = trimMantissa(decimal);
  }

  if (options.separator) {
    decimal = applySeparator(decimal);
  }

  return `${prefix}${decimal}${space}${postfix}`;
}

function trimMantissa(decimal: string) {
  return decimal.replace(/(0+)$/, "").replace(/\.$/, ".0");
}

function applySeparator(decimal: string) {
  const [char, mant] = decimal.split(".");
  return [char.replace(/\B(?<!\.\d*)(?=(\d{3})+(?!\d))/g, ","), mant].join(".");
}

function applyMantissa(decimal: string, mantissa: number) {
  return round(decimal, mantissa);
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
