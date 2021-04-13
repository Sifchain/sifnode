#!/usr/bin/env node

const fs = require("fs");
const { resolve } = require("path");

function loadData(location) {
  return JSON.parse(
    fs.readFileSync(location, {
      encoding: "utf-8",
    }),
  );
}

function saveData(location, data) {
  const parsed = Buffer.from(JSON.stringify(data, null, 2) + "\n"); // Add line break for linting to be happy
  fs.writeFileSync(location, parsed);
}

function updateERowan(asset) {
  const location = resolve(
    __dirname,
    "../../../smart-contracts/build/contracts/BridgeToken.json",
  );

  const {
    networks: {
      5777: { address },
    },
  } = loadData(location);

  return { ...asset, address };
}

function updateToken(contractName, asset) {
  const location = resolve(
    __dirname,
    `../../chains/ethereum/build/contracts/${contractName}.json`,
  );

  const {
    networks: {
      5777: { address },
    },
  } = loadData(location);
  return { ...asset, address };
}

// ASSET ADDRESSES

const assetsEthereumLocation = resolve(
  __dirname,
  "../src/assets.ethereum.localnet.json",
);

const data = loadData(assetsEthereumLocation);

data.assets = data.assets.map((asset) => {
  switch (asset.symbol) {
    case "atk":
      return updateToken("AliceToken", asset);
    case "btk":
      return updateToken("BobToken", asset);
    case "usdc":
      return updateToken("UsdCoin", asset);
    case "link":
      return updateToken("LinkCoin", asset);
    case "erowan":
      return updateERowan(asset);
  }
  return asset;
});

saveData(assetsEthereumLocation, data);

// BRIDGEBANK ADDRESS

function updateBridgeBankLocation() {
  // update bridgeBank location
  const configLocalnetLocation = resolve(
    __dirname,
    "../src/config.localnet.json",
  );

  const configData = loadData(configLocalnetLocation);
  const bridgeBankLocation = resolve(
    __dirname,
    "../../../smart-contracts/build/contracts/BridgeBank.json",
  );

  const {
    networks: {
      5777: { address: bridgeBankAddress },
    },
  } = loadData(bridgeBankLocation);

  configData.bridgebankContractAddress = bridgeBankAddress;

  saveData(configLocalnetLocation, configData);
}

updateBridgeBankLocation();
