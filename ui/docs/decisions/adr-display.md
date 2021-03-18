# Amount Display API

```ts
type IFormatOptionsNumber = {
  mode?: "number"; // default is number
  mantissa?: number; // number of decimals after point default 0
  separator?: boolean; // Add thousand separators eg. 1,000 default false
  trimMantissa?: boolean; // Remove 0s from the mantissa default false
  forceSign?: boolean; // Ensure we have a + sign at the start of the value default false
  exponent?: number; // multiplier when undefined exponent will be set by (amount as IAssetAmount).decimals ?? 1
  significant?: number; // When not undefined this will give us significant digits. default undefined
};

type IFormatOptionsPercent = {
  mode: "percent";
  mantissa?: number; // number of decimals after point default 2
  trimMantissa?: boolean; // Remove 0s from the mantissa default false
  exponent?: number; // multiplier when undefined exponent will be 2
  space?: boolean; // separate prefix and suffix with spaces default false
};

type IFormatOptions = IFormatOptionsNumber | IFormatOptionsPercent;
```

```ts
function format(
  amount: IAmount | IAssetAmount,
  options: IFormatOptions
): string;
```

# Examples

| BigInt             | (as IAssetAmount).decimals | IFormatOptions                                          | output               |
| ------------------ | -------------------------- | ------------------------------------------------------- | -------------------- |
| 100000000000       | undefined                  | { mantissa: 2, separator: true }                        | `100,000,000,000.00` |
| 990000000000000000 | 18                         | { mantissa: 6 }                                         | `0.990000`           |
| 990000000000000000 | 18                         | { mantissa: 6, trimMantissa: true }                     | `0.99`               |
| 999999800000000000 | 18                         | { mantissa: 8 }                                         | `0.9999998`          |
| 100                | undefined                  | { mode:"percent", mantissa: 1 }                         | `1.0%`               |
| 1000               | undefined                  | { mode:"percent", mantissa: 2 }                         | `10.00%`             |
| 12345              | undefined                  | { mode:"percent", mantissa: 3, exponent: 3}             | `12.345%`            |
| 12345              | undefined                  | { mode:"percent", mantissa: 3, exponent: 3, space:true} | `12.345 %`           |
