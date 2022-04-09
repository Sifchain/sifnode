#!/usr/bin/env node

const fs = require("fs");

async function createAssetRewardsFile() {
  const csv = fs.readFileSync("./pools.csv", "utf-8");
  const entries = JSON.parse(fs.readFileSync("./entries.json", "utf-8")).result
    .registry.entries;
  const lines = csv.split("\r\n");

  let [, allocation] = lines[0].split('"');
  allocation = parseInt(allocation.trim().split(",").join(""));

  const multipliers = lines.slice(1).map((line) => {
    const [, poolName, multiplier] = line.split(",");

    const entry = entries.find(
      ({ base_denom }) =>
        base_denom === `c${poolName}` ||
        base_denom === `u${poolName}` ||
        base_denom === `e${poolName}` ||
        base_denom === `${poolName}`
    );

    const asset = entry ? entry.denom : `???${poolName}???`;

    return {
      asset,
      multiplier,
    };
  });

  const rewards = {
    id: "RP_1",
    start_block: 1,
    end_block: 100,
    allocation,
    multipliers,
  };

  fs.writeFileSync("./rewards.json", JSON.stringify(rewards, null, 2));
}

async function start() {
  await createAssetRewardsFile();
}

start();
