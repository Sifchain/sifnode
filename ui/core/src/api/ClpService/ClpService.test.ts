import createClpService from ".";
import { AssetAmount } from "../../entities";
import { getTestingTokens } from "../../test/utils/getTestingToken";

const [ROWAN, CATK, CBTK] = getTestingTokens(["ROWAN", "CATK", "CBTK"]);

let service: ReturnType<typeof createClpService>;

beforeEach(() => {
  service = createClpService({
    nativeAsset: ROWAN,
    sifChainId: "sifchain",
    sifApiUrl: "http://localhost:1317",
    sifWsUrl: "ws://localhost:26667/websocket",
  });
});

test("getPools()", async () => {
  const pools = await service.getPools();

  expect(pools.map(pool => pool.toString())).toEqual([
    "10000000.000000000000000000 ROWAN | 10000000.000000000000000000 CATK",
    "10000000.000000000000000000 ROWAN | 10000000.000000000000000000 CBTK",
    "10000000.000000000000000000 ROWAN | 8300.000000000000000000 CETH",
    "10000000.000000000000000000 ROWAN | 588235.000000000000000000 CLINK",
    "10000000.000000000000000000 ROWAN | 10000000.000000000000000000 CUSDC",
  ]);
});

test("addLiquidity", async () => {
  const message = await service.addLiquidity({
    fromAddress: "sif1l7hypmqk2yc334vc6vmdwzp5sdefygj2ad93p5",
    externalAssetAmount: AssetAmount(CATK, "1000"),
    nativeAssetAmount: AssetAmount(ROWAN, "1000"),
  });
  expect(message).toEqual({
    type: "cosmos-sdk/StdTx",
    value: {
      fee: { amount: [], gas: "200000" },
      memo: "",
      msg: [
        {
          type: "clp/AddLiquidity",
          value: {
            ExternalAsset: { symbol: "catk" },
            ExternalAssetAmount: "1000000000000000000000",
            NativeAssetAmount: "1000000000000000000000",
            Signer: "sif1l7hypmqk2yc334vc6vmdwzp5sdefygj2ad93p5",
          },
        },
      ],
      signatures: null,
    },
  });
});

test("removeLiquidity()", async () => {
  const message = await service.removeLiquidity({
    fromAddress: "sif1l7hypmqk2yc334vc6vmdwzp5sdefygj2ad93p5",
    asset: CATK,
    asymmetry: "0",
    wBasisPoints: "10000",
  });

  expect(message).toEqual({
    type: "cosmos-sdk/StdTx",
    value: {
      fee: { amount: [], gas: "200000" },
      memo: "",
      msg: [
        {
          type: "clp/RemoveLiquidity",
          value: {
            Asymmetry: "0",
            ExternalAsset: { symbol: "catk" },
            Signer: "sif1l7hypmqk2yc334vc6vmdwzp5sdefygj2ad93p5",
            WBasisPoints: "10000",
          },
        },
      ],
      signatures: null,
    },
  });
});

test("createPool()", async () => {
  const message = await service.createPool({
    fromAddress: "sif1l7hypmqk2yc334vc6vmdwzp5sdefygj2ad93p5",
    externalAssetAmount: AssetAmount(CATK, "1000"),
    nativeAssetAmount: AssetAmount(ROWAN, "1000"),
  });

  expect(message).toEqual({
    type: "cosmos-sdk/StdTx",
    value: {
      fee: { amount: [], gas: "200000" },
      memo: "",
      msg: [
        {
          type: "clp/CreatePool",
          value: {
            ExternalAsset: { symbol: "catk" },
            ExternalAssetAmount: "1000000000000000000000",
            NativeAssetAmount: "1000000000000000000000",
            Signer: "sif1l7hypmqk2yc334vc6vmdwzp5sdefygj2ad93p5",
          },
        },
      ],
      signatures: null,
    },
  });
});

test("swap()", async () => {
  const message = await service.swap({
    fromAddress: "sif1l7hypmqk2yc334vc6vmdwzp5sdefygj2ad93p5",
    receivedAsset: CATK,
    sentAmount: AssetAmount(CBTK, "1000"),
    minimumReceived: AssetAmount(CATK, "1000", { inBaseUnit: true }),
  });

  expect(message).toEqual({
    type: "cosmos-sdk/StdTx",
    value: {
      fee: { amount: [], gas: "200000" },
      memo: "",
      msg: [
        {
          type: "clp/Swap",
          value: {
            ReceivedAsset: { symbol: "catk" },
            SentAmount: "1000000000000000000000",
            MinReceivingAmount: "1000",
            SentAsset: { symbol: "cbtk" },
            Signer: "sif1l7hypmqk2yc334vc6vmdwzp5sdefygj2ad93p5",
          },
        },
      ],
      signatures: null,
    },
  });
});
test("getLiquidityProvider()", async () => {
  const lp = await service.getLiquidityProvider({
    lpAddress: "sif1l7hypmqk2yc334vc6vmdwzp5sdefygj2ad93p5",
    symbol: "catk",
  });

  expect(lp?.asset.symbol).toEqual("catk");
  expect(lp?.address).toEqual("sif1l7hypmqk2yc334vc6vmdwzp5sdefygj2ad93p5");
  expect(lp?.units.toFixed(0)).toEqual("10000000000000000000000000");
});
