# Create a display API to universally display values throughout frontend

* Status: proposed
* Deciders: Michael Pierce, Thomas Davis, Rudi Yardley
* Date: 2021-03-18

# Amount Display API

```ts
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

type IFormatOptions = IFormatOptionsMantissa | IFormatOptionsShorthandTotalLength;
```

```ts
function format(
  amount: IAmount | IAssetAmount,
  options: IFormatOptions
): string;
```

# Examples

| BigInt              | (as IAssetAmount).decimals | IFormatOptions                                          | output               |
| ------------------- | -------------------------- | ------------------------------------------------------- | -------------------- |
| 100000000000        | undefined                  | { mantissa: 2, separator: true }                        | `100,000,000,000.00` |
| 100000000000        | undefined                  | { shorthand: true }                                         | `100b`               |
| 100000000000        | undefined                  | { shorthand: true, mantissa:6 }                             | `100.000000b`        |
| 990000000000000000  | 18                         | { mantissa: 6 }                                         | `0.990000`           |
| 990000000000000000  | 18                         | { mantissa: 6, trimMantissa: true }                     | `0.99`               |
| 999999800000000000  | 18                         | { mantissa: 8 }                                         | `0.9999998`          |
| 100                 | undefined                  | { mode:"percent", mantissa: 1 }                         | `1.0%`               |
| 1000                | undefined                  | { mode:"percent", mantissa: 2 }                         | `10.00%`             |
| 12345               | undefined                  | { mode:"percent", mantissa: 3, exponent: 3}             | `12.345%`            |
| 12345               | undefined                  | { mode:"percent", mantissa: 3, exponent: 3, space:true} | `12.345 %`           |
| -990000000000000000 | 18                         | { mantissa: 6, trimMantissa: true }                     | `-0.99`              |
| 999999800000000000  | 18                         | { mantissa: 8, forceSign:true }                         | `+0.9999998`         |
| 999999800000000000  | 18                         | { mantissa: 8, forceSign:true, space:true }             | `+ 0.9999998`        |
