# Amount Api

An amount is an integer representation of a value used in calculation in the system

The idea here is to provide a wrapper for underlying library (JSBI/Fraction) to handle internal representation.

All methods that accept string as an amount convert that string amount via `Amount(str)`

```ts
interface IAmount {
  toBigInt(): JSBI;
  toString(): string; // string representation of integer
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
}

interface IAsset {
  symbol: string; // Eg. ceth
  label: string; // Eg. cETH
  name: string; // Eg. Ethereum
  decimals: number; // 18
  network: "ethereum" | "sifchain"; // | "bitcoin"; etc
  address: string; // All assets must have an address with  (0x00000000000000000000000000000000 for a native coin)
  imageUrl?: string;
}

// I think by doing this we will make it so instead of `someAmount.asset.symbol` we can just have `someAmount.symbol`
interface IAssetAmount extends IAmount, IAsset {}

function Amount(source: BigintIsh | IAmount): IAmount;
function AssetAmount(asset: Asset, amount: IAmount): IAmount;
```

We need to be able to create an amount easily from a string or JSBI

```ts
const amount = Amount("1236479876134");
const hundred = Amount(JSBI.BigInt("100"));
```

We can use some basic static values as convenience

```ts
Amount.ZERO;
Amount.TEN;
Amount._100;
Amount._1000;
Amount.ONE;
```

Despite using the Fraction floating point internal representation we deliberately don't have a feature for providing a floating point number to Amount - this should work so long as all values are in base units.

- [ ] Create new `Amount` wrapper - implement all Fraction methods - write tests
- [ ] Create new `AssetAmount` wrapper - write tests
- [ ] Add `label` on `Asset`
- [ ] Create display lib to handle `toFixed` `toSignificant` etc. -
- [ ] Move all application internal representations to baseunits
- [ ] Replace `Fraction` with `Amount` throughout codebase
- [ ] Replace `AssetAmount` with new `AssetAmount` throughout codebase
- [ ] Remove `Coin` and `Token` distinction they are not helpful and make `address` non optional on `Asset` initializing with `0x00000...` where appropriate
