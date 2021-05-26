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

type IFormatOptionsMantissa<
  M = number | DynamicMantissa
> = IFormatOptionsBase & {
  shorthand?: boolean;
  mantissa?: M; // number of decimals after point default is exponent
  trimMantissa?: boolean | "integer"; // Remove 0s from the mantissa default false
};

type IFormatOptionsShorthandTotalLength = IFormatOptionsBase & {
  shorthand: true;
  totalLength?: number; // This will give us significant digits using abbreviations eg. `1.234k` it will override anything in mantissa
};

export type DynamicMantissa = Record<number | "infinity", number>;

export type IFormatOptions =
  | IFormatOptionsMantissa
  | IFormatOptionsShorthandTotalLength;

type IFormatOptionsFixedMantissa =
  | IFormatOptionsMantissa<number>
  | IFormatOptionsShorthandTotalLength;

function isAsset(val: any): val is IAsset {
  return !!val && typeof val?.symbol === "string";
}

/**
 * Takes an amount and a dynamic mantissa hash and returns the mantisaa value to use
 * @param amount amount given to format function
 * @param hash dynamic value hash to calculate mantissa from
 * @returns number of mantissa to send to formatter
 */
export function getMantissaFromDynamicMantissa(
  amount: IAmount,
  hash: DynamicMantissa,
) {
  const { infinity, ...numHash } = hash;

  const entries = Object.entries(numHash);

  entries.sort(([a], [b]) => {
    if (a > b) return 1;
    return -1;
  });

  for (const entry of entries) {
    const [range, mantissa] = entry;
    if (amount.lessThan(range)) {
      return mantissa;
    }
  }

  if (amount.lessThan("10000")) {
    return 2;
  }

  return infinity;
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

function isDynamicMantissa(
  value: undefined | number | DynamicMantissa,
): value is DynamicMantissa {
  return typeof value !== "number";
}

function isOptionsWithFixedMantissa(
  options: IFormatOptionsFixedMantissa | IFormatOptions,
): options is IFormatOptionsFixedMantissa {
  return options.shorthand || !isDynamicMantissa(options.mantissa);
}

/**
 * Options come with a dynamic or fixed mantissa. This function converts a dynamic mantissa value if it exists to a fixed number
 * @param amount
 * @param options
 * @returns
 */
function convertDynamicMantissaToFixedMantissa(
  amount: IAmount,
  options: IFormatOptions,
): IFormatOptionsFixedMantissa {
  if (
    !isOptionsWithFixedMantissa(options) &&
    typeof options.mantissa === "object"
  ) {
    return {
      ...options,
      mantissa: getMantissaFromDynamicMantissa(amount, options.mantissa),
    };
  }
  return options as IFormatOptionsFixedMantissa;
}

export type AmountNotAssetAmount<T extends IAmount> = T extends IAssetAmount
  ? never
  : T;

export function format<T extends IAmount>(
  amount: AmountNotAssetAmount<T>,
): string;
export function format<T extends IAmount>(
  amount: AmountNotAssetAmount<T>,
  asset: Exclude<IAsset, IAssetAmount>,
): string;
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
  _asset?: Exclude<IAsset, IAssetAmount> | IFormatOptions,
  _options?: IFormatOptions,
): string {
  const amount = _amount;
  const _optionsWithDynamicMantissa =
    (isAsset(_asset) ? _options : _asset) || {};
  const asset = isAsset(_asset) ? _asset : undefined;

  const options = convertDynamicMantissaToFixedMantissa(
    amount,
    _optionsWithDynamicMantissa,
  );

  // This should not happen in typed parts of the codebase
  if (typeof amount === "string") {
    // We need this in order to push developers to use the amount API right to the point at which we format values for display
    // Currently not using JSX means types are not necessarily propagated to every view so types guards
    // and there was a happy coincidence that format happened to work with a string and no asset
    //
    // We need to avoid this for the following reasons:
    //   * It encourages the status quo of not using JSX which has many poor knockon effects
    //   * One way api leads to simpler and easier to understand code
    //   * It reduces refactorability
    //   * It adds complexity to the codebase as it enables accidental amount -> string -> amount flows
    //   * It makes it more likely that developers accidentally try to format AssetAmounts as Amounts which
    //     is something this function attempts to solve using Types
    //   * It adds difficult to track down errors as strings of unknown format are passed to the format function
    //
    // Once JSX is used throughout the codebase it might be time to revisit this
    throw new Error(
      "Amount can only take an IAmount and must NOT be a string. If you have a string and need to format it you should first convert it to an IAmount. Eg. format(Amount('100'), myformat)",
    );
  }

  if (!amount) {
    // In theory this should not happen if we are using typescript correctly
    // This might happen due to a service response not being runtime checked
    // or in Vue because we are not using JSX templates
    console.error(`Amount "${amount}" supplied to format function is falsey`);
    return ""; // return empty string if there is an error
  }

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
    decimal = trimMantissa(decimal, options.trimMantissa === "integer");
  }

  if (options.separator) {
    decimal = applySeparator(decimal);
  }

  return `${prefix}${decimal}${space}${postfix}`;
}

export function trimMantissa(decimal: string, integer = false) {
  return decimal.replace(/(0+)$/, "").replace(/\.$/, integer ? "" : ".0");
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

function createNumbroConfig(options: IFormatOptionsFixedMantissa) {
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
          trimMantissa: !!options.trimMantissa,
        }),
  };
}
