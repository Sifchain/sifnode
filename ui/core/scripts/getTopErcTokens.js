#!/usr/bin/env node
const fs = require("fs");
const Web3 = require("web3");

const ABI = [
  {
    constant: true,
    inputs: [],
    name: "name",
    outputs: [
      {
        name: "",
        type: "string",
      },
    ],
    payable: false,
    stateMutability: "view",
    type: "function",
  },
  {
    constant: false,
    inputs: [
      {
        name: "_spender",
        type: "address",
      },
      {
        name: "_value",
        type: "uint256",
      },
    ],
    name: "approve",
    outputs: [
      {
        name: "",
        type: "bool",
      },
    ],
    payable: false,
    stateMutability: "nonpayable",
    type: "function",
  },
  {
    constant: true,
    inputs: [],
    name: "totalSupply",
    outputs: [
      {
        name: "",
        type: "uint256",
      },
    ],
    payable: false,
    stateMutability: "view",
    type: "function",
  },
  {
    constant: false,
    inputs: [
      {
        name: "_from",
        type: "address",
      },
      {
        name: "_to",
        type: "address",
      },
      {
        name: "_value",
        type: "uint256",
      },
    ],
    name: "transferFrom",
    outputs: [
      {
        name: "",
        type: "bool",
      },
    ],
    payable: false,
    stateMutability: "nonpayable",
    type: "function",
  },
  {
    constant: true,
    inputs: [],
    name: "decimals",
    outputs: [
      {
        name: "",
        type: "uint8",
      },
    ],
    payable: false,
    stateMutability: "view",
    type: "function",
  },
  {
    constant: true,
    inputs: [
      {
        name: "_owner",
        type: "address",
      },
    ],
    name: "balanceOf",
    outputs: [
      {
        name: "balance",
        type: "uint256",
      },
    ],
    payable: false,
    stateMutability: "view",
    type: "function",
  },
  {
    constant: true,
    inputs: [],
    name: "symbol",
    outputs: [
      {
        name: "",
        type: "string",
      },
    ],
    payable: false,
    stateMutability: "view",
    type: "function",
  },
  {
    constant: false,
    inputs: [
      {
        name: "_to",
        type: "address",
      },
      {
        name: "_value",
        type: "uint256",
      },
    ],
    name: "transfer",
    outputs: [
      {
        name: "",
        type: "bool",
      },
    ],
    payable: false,
    stateMutability: "nonpayable",
    type: "function",
  },
  {
    constant: true,
    inputs: [
      {
        name: "_owner",
        type: "address",
      },
      {
        name: "_spender",
        type: "address",
      },
    ],
    name: "allowance",
    outputs: [
      {
        name: "",
        type: "uint256",
      },
    ],
    payable: false,
    stateMutability: "view",
    type: "function",
  },
  {
    payable: true,
    stateMutability: "payable",
    type: "fallback",
  },
  {
    anonymous: false,
    inputs: [
      {
        indexed: true,
        name: "owner",
        type: "address",
      },
      {
        indexed: true,
        name: "spender",
        type: "address",
      },
      {
        indexed: false,
        name: "value",
        type: "uint256",
      },
    ],
    name: "Approval",
    type: "event",
  },
  {
    anonymous: false,
    inputs: [
      {
        indexed: true,
        name: "from",
        type: "address",
      },
      {
        indexed: true,
        name: "to",
        type: "address",
      },
      {
        indexed: false,
        name: "value",
        type: "uint256",
      },
    ],
    name: "Transfer",
    type: "event",
  },
];

const web3 = new Web3(
  new Web3.providers.HttpProvider(
    "https://mainnet.infura.io/v3/93cd052103fd44bd9cf855654e5804ac"
  )
);

async function updateCoins() {
  const CoinGecko = require("coingecko-api");
  const CoinGeckoClient = new CoinGecko();
  let data = await CoinGeckoClient.coins.all({
    page: 1,
    order: "market_cap_desc",
    per_page: 100,
  });

  const assets = data.data.map(
    ({ id, symbol, name, image, market_data: { market_cap_rank } }) => ({
      id,
      symbol,
      name,
      image,
      market_cap_rank,
    })
  );

  const coins = [];
  for (let asset of assets) {
    coin = await CoinGeckoClient.coins.fetch(asset.id);

    // Avoid rate limiting
    await new Promise((res) => setTimeout(res, 2000));

    console.log([coin.data.id, coin.data.market_cap_rank]);
    coins.push(coin.data);
  }

  fs.writeFileSync("data/coins.json", JSON.stringify(coins, null, 2));

  const erc20Tokens = coins.filter((coin) => {
    return coin.contract_address && coin.asset_platform_id === "ethereum";
  });

  const erc20 = erc20Tokens.map(
    ({
      id,
      name,
      symbol,
      image,
      contract_address,
      asset_platform_id,
      market_cap_rank,
    }) => ({
      id,
      name,
      symbol,
      image,
      contract_address,
      asset_platform_id,
      market_cap_rank,
    })
  );

  fs.writeFileSync("data/erc20.json", JSON.stringify(erc20, null, 2));
}

async function enrichERC() {
  const tokens = JSON.parse(fs.readFileSync("erc20.json"));
  const erc20 = [];
  for (let token of tokens) {
    const contract = new web3.eth.Contract(ABI, token.contract_address);
    const decimals = await contract.methods.decimals().call();
    erc20.push({ ...token, decimals: parseInt(decimals) });
  }
  fs.writeFileSync("data/topErc20Tokens.json", JSON.stringify(erc20, null, 2));
}

// updateCoins();
enrichERC();
