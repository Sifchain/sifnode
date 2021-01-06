#!/usr/bin/env node

const fs = require("fs");
const { resolve } = require("path");

function loadData(location) {
  return JSON.parse(
    fs.readFileSync(location, {
      encoding: "utf-8",
    })
  );
}

function saveData(location, data) {
  const parsed = Buffer.from(JSON.stringify(data, null, 2));
  fs.writeFileSync(location, parsed);
}

function updateERowan(asset) {
  const location = resolve(
    __dirname,
    "../../../smart-contracts/build/contracts/BridgeToken.json"
  );

  const {
    networks: {
      5777: { address },
    },
  } = loadData(location);

  return { ...asset, address };
}

function updateAtk(asset) {
  const location = resolve(
    __dirname,
    "../../chains/ethereum/build/contracts/AliceToken.json"
  );

  const {
    networks: {
      5777: { address },
    },
  } = loadData(location);
  return { ...asset, address };
}

function updateBtk(asset) {
  const location = resolve(
    __dirname,
    "../../chains/ethereum/build/contracts/BobToken.json"
  );
  const {
    networks: {
      5777: { address },
    },
  } = loadData(location);
  return { ...asset, address };
}

const configLocation = resolve(
  __dirname,
  "../src/assets.ethereum.localnet.json"
);

const data = loadData(configLocation);

data.assets = data.assets.map((asset) => {
  switch (asset.symbol) {
    case "atk":
      return updateAtk(asset);
    case "btk":
      return updateBtk(asset);
    case "erowan":
      return updateERowan(asset);
  }
  return asset;
});

saveData(configLocation, data);
