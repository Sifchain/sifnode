import {
  isBroadcastTxFailure,
  makeCosmoshubPath,
  Secp256k1HdWallet,
  SigningCosmosClient,
} from "@cosmjs/launchpad";
import axios from "axios";
import fetch from "cross-fetch";
import { juniper } from "./test/utils/accounts";

test("", async () => {
  const getPoolsResp = await fetch("http://127.0.0.1:1317/clp/getPools", {
    headers: {
      accept: "application/json, text/plain, */*",
      "accept-language": "en-US,en;q=0.9,la;q=0.8",
      "cache-control": "no-cache",
      pragma: "no-cache",
      "sec-ch-ua":
        '"Google Chrome";v="87", " Not;A Brand";v="99", "Chromium";v="87"',
      "sec-ch-ua-mobile": "?0",
      "sec-fetch-dest": "empty",
      "sec-fetch-mode": "cors",
      "sec-fetch-site": "cross-site",
    },
    referrer: "http://localhost:8080/",
    referrerPolicy: "strict-origin-when-cross-origin",
    body: null,
    method: "GET",
    mode: "cors",
    credentials: "omit",
  });

  const getPoolsRespJson = await getPoolsResp.json();

  expect(getPoolsRespJson.result.Pools).toEqual([
    {
      external_asset: {
        symbol: "catk",
      },
      native_asset_balance: "10000000000000000000000000",
      external_asset_balance: "10000000000000000000000000",
      pool_units: "10000000000000000000000000",
    },
    {
      external_asset: {
        symbol: "cbtk",
      },
      native_asset_balance: "10000000000000000000000000",
      external_asset_balance: "10000000000000000000000000",
      pool_units: "10000000000000000000000000",
    },
    {
      external_asset: {
        symbol: "ceth",
      },
      external_asset_balance: "8300000000000000000000",
      native_asset_balance: "10000000000000000000000000",
      pool_units: "10000000000000000000000000",
    },
    {
      external_asset: {
        symbol: "clink",
      },
      native_asset_balance: "10000000000000000000000000",
      external_asset_balance: "588235000000000000000000",
      pool_units: "10000000000000000000000000",
    },
    {
      external_asset: {
        symbol: "cusdc",
      },
      native_asset_balance: "10000000000000000000000000",
      external_asset_balance: "10000000000000000000000000",
      pool_units: "10000000000000000000000000",
    },
  ]);

  const addLiquidityBody = JSON.stringify({
    base_req: {
      chain_id: "sifchain",
      from: "sif1m625hcmnkc84cgmef6upzzyfu6mxd4jkpnfwwl",
    },
    external_asset: {
      source_chain: "sifchain",
      symbol: "ceth",
      ticker: "ceth",
    },
    external_asset_amount: "5000000000000000000",
    native_asset_amount: "1000000000000000000000",
    signer: "sif1m625hcmnkc84cgmef6upzzyfu6mxd4jkpnfwwl",
  });

  const addLiquidityResp = await fetch(
    "http://127.0.0.1:1317/clp/addLiquidity",
    {
      headers: {
        accept: "application/json, text/plain, */*",
        "accept-language": "en-US,en;q=0.9,la;q=0.8",
        "cache-control": "no-cache",
        "content-type": "application/json;charset=UTF-8",
        pragma: "no-cache",
        "sec-ch-ua":
          '"Google Chrome";v="87", " Not;A Brand";v="99", "Chromium";v="87"',
        "sec-ch-ua-mobile": "?0",
        "sec-fetch-dest": "empty",
        "sec-fetch-mode": "cors",
        "sec-fetch-site": "cross-site",
      },
      referrer: "http://localhost:8080/",
      referrerPolicy: "strict-origin-when-cross-origin",
      body: addLiquidityBody,
      method: "POST",
      mode: "cors",
      credentials: "omit",
    }
  );
  const stdtx = await addLiquidityResp.json();
  expect(stdtx).toEqual({
    type: "cosmos-sdk/StdTx",
    value: {
      msg: [
        {
          type: "clp/AddLiquidity",
          value: {
            Signer: "sif1m625hcmnkc84cgmef6upzzyfu6mxd4jkpnfwwl",
            ExternalAsset: { symbol: "ceth" },
            NativeAssetAmount: "1000000000000000000000",
            ExternalAssetAmount: "5000000000000000000",
          },
        },
      ],
      fee: { amount: [], gas: "200000" },
      signatures: null,
      memo: "",
    },
  });

  const lpResponse = await fetch(
    `http://127.0.0.1:1317/clp/getLiquidityProvider?symbol=ceth&lpAddress=sif1m625hcmnkc84cgmef6upzzyfu6mxd4jkpnfwwl`,
    {
      headers: {
        accept: "application/json, text/plain, */*",
        "accept-language": "en-US,en;q=0.9,la;q=0.8",
        "cache-control": "no-cache",
        pragma: "no-cache",
        "sec-ch-ua":
          '"Google Chrome";v="87", " Not;A Brand";v="99", "Chromium";v="87"',
        "sec-ch-ua-mobile": "?0",
        "sec-fetch-dest": "empty",
        "sec-fetch-mode": "cors",
        "sec-fetch-site": "cross-site",
      },
      referrer: "http://localhost:8080/",
      referrerPolicy: "strict-origin-when-cross-origin",
      body: null,
      method: "GET",
      mode: "cors",
      credentials: "omit",
    }
  );

  expect(await lpResponse.json()).toEqual({
    error: "liquidity Provider does not exist",
  });

  const wallet = await Secp256k1HdWallet.fromMnemonic(
    juniper.mnemonic,
    makeCosmoshubPath(0),
    "sif"
  );
  const client = new SigningCosmosClient(
    "http://127.0.0.1:1317",
    juniper.address,
    wallet
  );

  const receipt = await client.signAndBroadcast(stdtx.value.msg, {
    amount: [],
    gas: "200000",
  });

  expect(isBroadcastTxFailure(receipt)).toBe(false);

  const getPools2Resp = await fetch("http://127.0.0.1:1317/clp/getPools", {
    headers: {
      accept: "application/json, text/plain, */*",
      "accept-language": "en-US,en;q=0.9,la;q=0.8",
      "cache-control": "no-cache",
      pragma: "no-cache",
      "sec-ch-ua":
        '"Google Chrome";v="87", " Not;A Brand";v="99", "Chromium";v="87"',
      "sec-ch-ua-mobile": "?0",
      "sec-fetch-dest": "empty",
      "sec-fetch-mode": "cors",
      "sec-fetch-site": "cross-site",
    },
    referrer: "http://localhost:8080/",
    referrerPolicy: "strict-origin-when-cross-origin",
    body: null,
    method: "GET",
    mode: "cors",
    credentials: "omit",
  });

  const getPools2RespJson = await getPools2Resp.json();

  expect(getPools2RespJson.result.Pools).toEqual([
    {
      external_asset: {
        symbol: "catk",
      },
      native_asset_balance: "10000000000000000000000000",
      external_asset_balance: "10000000000000000000000000",
      pool_units: "10000000000000000000000000",
    },
    {
      external_asset: {
        symbol: "cbtk",
      },
      native_asset_balance: "10000000000000000000000000",
      external_asset_balance: "10000000000000000000000000",
      pool_units: "10000000000000000000000000",
    },
    {
      external_asset: {
        symbol: "ceth",
      },
      external_asset_balance: "8305000000000000000000",
      native_asset_balance: "10001000000000000000000000",
      pool_units: "10003510285120826309483873",
    },
    {
      external_asset: {
        symbol: "clink",
      },
      native_asset_balance: "10000000000000000000000000",
      external_asset_balance: "588235000000000000000000",
      pool_units: "10000000000000000000000000",
    },
    {
      external_asset: {
        symbol: "cusdc",
      },
      native_asset_balance: "10000000000000000000000000",
      external_asset_balance: "10000000000000000000000000",
      pool_units: "10000000000000000000000000",
    },
  ]);
  const lp2Response = await fetch(
    `http://127.0.0.1:1317/clp/getLiquidityProvider?symbol=ceth&lpAddress=sif1m625hcmnkc84cgmef6upzzyfu6mxd4jkpnfwwl`,
    {
      headers: {
        accept: "application/json, text/plain, */*",
        "accept-language": "en-US,en;q=0.9,la;q=0.8",
        "cache-control": "no-cache",
        pragma: "no-cache",
        "sec-ch-ua":
          '"Google Chrome";v="87", " Not;A Brand";v="99", "Chromium";v="87"',
        "sec-ch-ua-mobile": "?0",
        "sec-fetch-dest": "empty",
        "sec-fetch-mode": "cors",
        "sec-fetch-site": "cross-site",
      },
      referrer: "http://localhost:8080/",
      referrerPolicy: "strict-origin-when-cross-origin",
      body: null,
      method: "GET",
      mode: "cors",
      credentials: "omit",
    }
  );
  expect((await lp2Response.json()).result).toEqual({
    LiquidityProvider: {
      asset: {
        symbol: "ceth",
      },
      liquidity_provider_address: "sif1m625hcmnkc84cgmef6upzzyfu6mxd4jkpnfwwl",
      liquidity_provider_units: "3510285120826309483873",
    },
    external_asset_balance: "2914268801405083967",
    height: "22",
    native_asset_balance: "3509404248386784438078",
  });

  // According to https://github.com/Sifchain/sifnode/blob/develop/docs/clp/Liquidity%20Pools%20Architecture.md

  // (nativeAssetBalance + externalAssetBalance) * (lpNativeAssetAmount * externalAssetBalance + nativeAssetBalance * lpExternalAsset) / (4 * nativeAssetBalance * externalAssetBalance)

  // externalAssetBalance: "8300000000000000000000",
  // nativeAssetBalance: "10000000000000000000000000",
  // lpNativeAssetAmount: "1000000000000000000000",
  // lpExternalAsset: "5000000000000000000"
  /*
   10000000000000000000000000 + 8300000000000000000000
1.00083e+25
> term1 = 10000000000000000000000000 + 8300000000000000000000
1.00083e+25
> term2 = (1000000000000000000000 * 8300000000000000000000)+ (10000000000000000000000000 * 5000000000000000000)
5.8300000000000005e+43
> numer = term1 * term2
5.834838900000001e+68
> denom = 4 * 10000000000000000000000000 * 8300000000000000000000
3.3200000000000002e+47
> numer / denom 
1.7574815963855423e+21

SO the answer should be something roughly like 1757481596385542300000 NOT 3510285120826309483873
   */
});
