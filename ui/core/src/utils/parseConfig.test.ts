import { Network } from "../entities";
import { parseConfig } from "./parseConfig";
const config = {
  sifAddrPrefix: "sif",
  sifChainId: "sifchain",
  sifApiUrl: "http://127.0.0.1:1317",
  sifWsUrl: "ws://localhost:26657/websocket",
  sifRpcUrl: "http://localhost:26657",
  web3Provider: "metamask",
  nativeAsset: "rowan",
  bridgebankContractAddress: "0x2C2B9C9a4a25e24B174f26114e8926a9f2128FE4",
  keplrChainConfig: {
    chainName: "Sifchain",
    chainId: "",
    rpc: "",
    rest: "",
    stakeCurrency: {
      coinDenom: "ROWAN",
      coinMinimalDenom: "rowan",
      coinDecimals: 18,
    },
    bip44: {
      coinType: 118,
    },
    bech32Config: {
      bech32PrefixAccAddr: "sif",
      bech32PrefixAccPub: "sifpub",
      bech32PrefixValAddr: "sifvaloper",
      bech32PrefixValPub: "sifvaloperpub",
      bech32PrefixConsAddr: "sifvalcons",
      bech32PrefixConsPub: "sifvalconspub",
    },
    currencies: [],
    feeCurrencies: [
      {
        coinDenom: "ROWAN",
        coinMinimalDenom: "rowan",
        coinDecimals: 18,
      },
    ],
    coinType: 118,
    gasPriceStep: {
      low: 0.5,
      average: 0.65,
      high: 0.8,
    },
  },
};

const expected = {
  assets: [
    {
      address: "0x2C2B9C9a4a25e24B174f26114e8926a9f2128FE4",
      decimals: 12,
      name: "123",
      network: "sifchain",
      symbol: "rowan",
    },
    {
      address: "0x2C2B9C9a4a25e24B174f26114e8926a9f2128FE4",
      decimals: 12,
      name: "123",
      network: "ethereum",
      symbol: "erowan",
    },
  ],
  bridgebankContractAddress: "0x2C2B9C9a4a25e24B174f26114e8926a9f2128FE4",
  bridgetokenContractAddress: "0x2C2B9C9a4a25e24B174f26114e8926a9f2128FE4",

  keplrChainConfig: {
    bech32Config: {
      bech32PrefixAccAddr: "sif",
      bech32PrefixAccPub: "sifpub",
      bech32PrefixConsAddr: "sifvalcons",
      bech32PrefixConsPub: "sifvalconspub",
      bech32PrefixValAddr: "sifvaloper",
      bech32PrefixValPub: "sifvaloperpub",
    },
    bip44: {
      coinType: 118,
    },
    chainId: "sifchain",
    chainName: "Sifchain",
    coinType: 118,
    currencies: [
      {
        coinDecimals: 12,
        coinDenom: "rowan",
        coinMinimalDenom: "rowan",
      },
    ],
    feeCurrencies: [
      {
        coinDecimals: 18,
        coinDenom: "ROWAN",
        coinMinimalDenom: "rowan",
      },
    ],
    gasPriceStep: {
      average: 0.65,
      high: 0.8,
      low: 0.5,
    },
    rest: "http://127.0.0.1:1317",
    rpc: "http://localhost:26657",
    stakeCurrency: {
      coinDecimals: 18,
      coinDenom: "ROWAN",
      coinMinimalDenom: "rowan",
    },
  },
  nativeAsset: {
    address: "0x2C2B9C9a4a25e24B174f26114e8926a9f2128FE4",
    decimals: 12,
    name: "123",
    network: "sifchain",
    symbol: "rowan",
  },
  sifAddrPrefix: "sif",
  sifApiUrl: "http://127.0.0.1:1317",
  sifChainId: "sifchain",
  sifWsUrl: "ws://localhost:26657/websocket",
};

test("parseConfig", () => {
  expect(
    parseConfig(config, [
      {
        address: "0x2C2B9C9a4a25e24B174f26114e8926a9f2128FE4",
        decimals: 12,
        name: "123",
        network: Network.SIFCHAIN,
        symbol: "rowan",
      },
      {
        address: "0x2C2B9C9a4a25e24B174f26114e8926a9f2128FE4",
        decimals: 12,
        name: "123",
        network: Network.ETHEREUM,
        symbol: "erowan",
      },
    ]),
  ).toMatchObject(expected);

  expect(() => {
    parseConfig(config, [
      {
        address: "0x2C2B9C9a4a25e24B174f26114e8926a9f2128FE4",
        decimals: 12,
        name: "123",
        network: Network.ETHEREUM,
        symbol: "rowan",
      },
    ]);
  }).toThrow();
  expect(() => {
    parseConfig(config, [
      {
        address: "0x2C2B9C9a4a25e24B174f26114e8926a9f2128FE4",
        decimals: 12,
        name: "123",
        network: Network.ETHEREUM,
        symbol: "erowan",
      },
    ]);
  }).toThrow();
});
