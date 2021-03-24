import { LcdClient } from "@cosmjs/launchpad";
import { Network } from "../../../../../entities";
import { ClpExtension, setupClpExtension } from "./index";

const removeLiquidityParams = {
  asymmetry: "10000",
  base_req: {
    chain_id: "sifchain",
    from: "",
  },
  external_asset: {
    symbol: "",
    source_chain: "",
    ticker: "",
  },
  signer: "",
  w_basis_points: "",
};

const swapParams = {
  sent_asset: {
    symbol: "",
    ticker: "",
    source_chain: "",
  },
  received_asset: {
    symbol: "",
    ticker: "",
    source_chain: "",
  },
  base_req: {
    from: "",
    chain_id: "",
  },
  signer: "",
  sent_amount: "",
  min_receiving_amount: "",
};

const liquidityParams = {
  base_req: {
    from: "",
    chain_id: "",
  },
  external_asset: {
    source_chain: "",
    symbol: "",
    ticker: "",
  },
  native_asset_amount: "",
  external_asset_amount: "",
  signer: "",
};
const getLiquidityProvider = {
  symbol: "foo",
  lpAddress: "bar",
};
let base: LcdClient;
let clp: ClpExtension["clp"];

beforeEach(() => {
  base = ({
    get: jest.fn(async () => "1234"),
    post: jest.fn(async () => "1234"),
  } as any) as LcdClient;
  clp = setupClpExtension(base).clp;
});

test("removeLiquidity", async () => {
  await clp.removeLiquidity(removeLiquidityParams);
  expect(base.post).toHaveBeenCalledWith(
    "/clp/removeLiquidity",
    removeLiquidityParams,
  );
});

test("getPools", async () => {
  await clp.getPools();
  expect(base.get).toHaveBeenCalledWith(`/clp/getPools`);
});

test("getAssets", async () => {
  await clp.getAssets("abc1234");
  expect(base.get).toHaveBeenCalledWith(`/clp/getAssets?lpAddress=abc1234`);
});

test("swap", async () => {
  await clp.swap(swapParams);
  expect(base.post).toHaveBeenCalledWith(`/clp/swap`, swapParams);
});

test("addLiquidity", async () => {
  await clp.addLiquidity(liquidityParams);
  expect(base.post).toHaveBeenCalledWith(`/clp/addLiquidity`, liquidityParams);
});

test("createPool", async () => {
  await clp.createPool(liquidityParams);
  expect(base.post).toHaveBeenCalledWith(`/clp/createPool`, liquidityParams);
});

test("getLiquidityProvider", async () => {
  await clp.getLiquidityProvider(getLiquidityProvider);
  expect(base.get).toHaveBeenCalledWith(
    `/clp/getLiquidityProvider?symbol=foo&lpAddress=bar`,
  );
});
