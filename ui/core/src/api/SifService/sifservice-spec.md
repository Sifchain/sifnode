Our SifService will wrap the `@cosmjs/launchpad` lib

- We will create extensions that match entries in the `x` folder in the main go project.

* clp
* ethbridge
* oracle

## clp

#### Types

```ts
enum SourceChain {
  ETHEREUM
  DASH
  SIF
  BITCOIN
}
```

```ts
type Pool = {
  // externalAsset: Asset; // - Dont need this b/c it is already within the AssetAmount type
  externalAssetBalance: AssetAmount;
  nativeAssetBalance: AssetAmount;
};
```

#### `GET /clp/getPools`

```ts
type ClpGetPools = () => {
  result: Pool[];
  height: string; // block height
};
```

#### `GET /clp/getPool`

```ts
type ClpGetPool = ({
  ticker: string,
  sourceChain: string,
}) => {
  result: Pool[];
  height: string; // block height
};
```

#### `GET /clp/getLiquidityProvider`

```ts
type ClpGetLiquidityProvider = ({
  ticker: string,
  lpAddress: string,
}) => {
  result: {
    liquidityProviderUnits: string;
    liquidityProviderAddress: string;
    asset: Asset;
  };
  height: string; // block height
};
```

### `POST /clp/createPool`

```ts
type ClpCmdCreatePool = ({
  sourceChain: SourceChain,
  asset: Asset,
  externalAssetAmount: AssetBalance,
  nativeAssetAmount: AssetBalance,
}) => any;
```

### `POST /clp/decommissionPool`

```ts
type ClpCmdDecommissionPool = ({ asset: Asset }) => any;
```

### `POST /clp/removeLiquidity`

```ts
type ClpCmdRemoveLiquidity = ({
  asset: Asset,
  wBasis: string,
  asymmetry: number,
}) => any;
```

### `POST /clp/addLiquidity`

```ts
type ClpCmdAddLiquidity = ({
  externalAssetAmount: AssetAmount,
  nativeAssetAmount: AssetAmount,
}) => any;
```

### `POST /clp/swap`

```ts
type ClpCmdSwap = ({
  sentAsset: Asset,
  receivedAsset: Asset,
  amount: AssetAmount,
}) => any;
```
