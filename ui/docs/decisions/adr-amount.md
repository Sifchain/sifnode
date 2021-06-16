# Create Amount API to universally handle values throughout frontend

* Status: proposed                                                   
* Deciders: Michael Pierce, Thomas Davis, Rudi Yardley
* Date: 2021-03-09                                                 

## Problem

We need to standardize how we handle amounts throughout the frontend. We started this project with a supertight deadline and so we looked to prior art and grabbed pieces from Uniswap to learn about the domain and throw together a working piece of software. With all of this we have a fair bit of baggage and API inconsistencies which we have identified as the most important and urgent but also most time consuming piece of technical debt.

## Goals

- Standardize how we treat amounts throughout the app.
- Use big integers to store and work with amounts in base units (ie. wei/satoshi) and remove reference to amounts in native format (ie. ether)
- Only convert upon display within a [standard companion display](adr-display.md) lib which will be addressed in another spec
- Focus on erganomics and using sensible memorable shorthands

## Amount

An amount is an integer representation of a value used in calculation within our system. The idea here is to provide a wrapper for underlying library (JSBI/Fraction/Big.js) to handle internal representation.

```ts
type IAmount = {
  // for use by display lib and in testing
  toBigInt(): JSBI;
  toString(): string;

  // for use elsewhere
  add(other: IAmount | string): IAmount;
  subtract(other: IAmount | string): IAmount;
  lessThan(other: IAmount | string): boolean;
  lessThanOrEqual(other: IAmount | string): boolean;
  equalTo(other: IAmount | string): boolean;
  greaterThan(other: IAmount | string): boolean;
  greaterThanOrEqual(other: IAmount | string): boolean;
  multiply(other: IAmount | string): IAmount;
  divide(other: IAmount | string): IAmount;
  sqrt(): IAmount;
};
```

There is a companion constructor function called `Amount()` that will generate an amount from various sources

```ts
function Amount(source: JSBI | bigint | string | IAmount): IAmount;
```

Here you can provide either a bigint, a JSBI, a string or another IAmount (say an AssetAmount) to convert to an amount.

All methods that accept string as an amount convert that string amount via `Amount(str)`

```ts
const hundred = Amount(JSBI.BigInt("100"));
const amount = Amount("1236479876134");
const amount = Amount(AssetAmount("eth", "100"));
```

We can use some basic static values as convenience

```ts
Amount.ZERO = Amount("0");
Amount.TEN = Amount("10");
Amount._100 = Amount("100");
Amount._1000 = Amount("1000");
Amount.ONE = Amount("1");
```

## Asset

An Asset represents a token denomination in our system. This ADR suggests two changes from our current setup.

1. Remove the `Token` and `Coin` distinction as they are generally not helpful instead all Coins are Assets with an address field that contains the string `0x0000000000000000000000000000000000000000`
2. Add a label field to use for displaying token labels. This would normally be the symbol with correct capitalization. Eg. `eROWAN` for the symbol `erowan`

Following is an interface which assets should conform

```ts
interface IAsset {
  readonly symbol: string; // Eg. ceth
  readonly label: string; // Eg. cETH
  readonly name: string; // Eg. Ethereum
  readonly decimals: number; // 18
  readonly network: "ethereum" | "sifchain"; // | "bitcoin"; etc
  readonly address: string; // All assets must have an address with (0x0000000000000000000000000000000000000000 for a native coin)
  readonly imageUrl?: string;
}
```

There is a companion constructor function called `Asset()` that will generate an asset from source data

```ts
function Asset(source: IAsset | string): IAsset;
```

Assets are initialized on app init and are cached if an asset has been cached you can use a string to identify it.

```ts
Asset("eth"); // returns the cached 'eth' asset
```

## AssetAmount

An `AssetAmount` is a convenience struct that represents an `Amount` that is connected to an `Asset`. Most of the values our application deals with will be `AssetAmount`s so it makes sense to have this construct. This is simply going to combine the two data constructs together and conform to the following interface:

```ts
interface IAssetAmount extends IAmount, IAsset {
  readonly asset: IAsset;
  readonly amount: IAmount;
}
```

When expanded this is effectively the same as the following:

```ts
interface IAssetAmout {
  // Getters for source structures
  readonly asset: IAsset;
  readonly amount: IAmount;

  // For use by display lib and in testing
  toBigInt(): JSBI;
  toString(): string;

  // For use within and outside core
  add(other: IAmount | string): IAmount;
  subtract(other: IAmount | string): IAmount;
  lessThan(other: IAmount | string): boolean;
  lessThanOrEqual(other: IAmount | string): boolean;
  equalTo(other: IAmount | string): boolean;
  greaterThan(other: IAmount | string): boolean;
  greaterThanOrEqual(other: IAmount | string): boolean;
  multiply(other: IAmount | string): IAmount;
  divide(other: IAmount | string): IAmount;
  sqrt(): IAmount;

  // Asset props
  readonly symbol: string;
  readonly label: string;
  readonly name: string;
  readonly decimals: number;
  readonly network: Network;
  readonly address: string;
  readonly imageUrl?: string;
}
```

This will be better than the system we have as currently we often need to reach into the assetAmount to get things like the symbol of the assetamount.

```ts
const amount: IAssetAmount = getAmount();
amount.symbol; // before this would have to be amount.asset.symbol
```

There is a companion constructor function called `AssetAmount()` that will generate an `AssetAmount` from source data

```ts
function AssetAmount(
  asset: IAsset | string,
  amount: IAmount | string
): IAssetAmount;
```

We can then use the string shorthands to initialize an asset amount:

```ts
const amount = AssetAmount("eth", "100"); // 100 wei
```

You can use the full Asset and Amount API to do so as well:

```ts
// Using string initiializers and constructor functions
const amount = AssetAmount(Asset("eth"), Amount("100")); // 100 wei

// Using raw data initiializers
const amount = AssetAmount(
  {
    symbol: "eth",
    decimals: 18,
    address: "0x0000000000000000000000000000000000000000",
    network: "ethereum",
    name: "Ethereum",
    label: "ETH",
    imageUrl:
      "https://assets.coingecko.com/coins/images/279/small/ethereum.png?1595348880",
  },
  Amount(JSBI.BigInt("100"))
); // 100 wei
```

Despite using the Fraction floating point internal representation we deliberately don't have a feature for providing a floating point number to Amount - this should work so long as all values are in base units. We will however need to refactor certain parts of the app that do this conversion within their calculations.

```ts
Amount("100.1234"); // will throw an error
```

## Migration Phases

- [ ] Create new `Amount` wrapper - implement all Fraction methods - write tests
- [ ] Create new `AssetAmount` wrapper - write tests
- [ ] Add `label` on `Asset`
- [ ] Create display lib to handle `toFixed` `toSignificant` etc. -
- [ ] Move all application internal representations to baseunits
- [ ] Replace `Fraction` with `Amount` throughout codebase
- [ ] Replace `AssetAmount` with new `AssetAmount` throughout codebase
- [ ] Remove `Coin` and `Token` distinction they are not helpful and make `address` non optional on `Asset` initializing with `0x00000...` where appropriate
