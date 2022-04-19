#!/usr/bin/env node

const fs = require("fs");

async function createAssetRewardsFile() {
  const csv = fs.readFileSync("./pools.csv", "utf-8");
  const entries = JSON.parse(fs.readFileSync("./entries.json", "utf-8")).result
    .registry.entries;
  const lines = csv.split("\r\n").filter((line) => line.split(",")[1] !== "");

  let [, allocation] = lines[0].split('"');
  allocation = `${allocation.trim().split(",").join("")}${"0".repeat(18)}`;

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
      pool_multiplier_asset: asset,
      multiplier,
    };
  });

  const rewards = [
    {
      reward_period_id: "RP_1",
      reward_period_start_block: 1,
      reward_period_end_block: 100,
      reward_period_allocation: allocation,
      reward_period_pool_multipliers: multipliers,
      reward_period_default_multiplier: "0.0",
    },
  ];

  fs.writeFileSync("./rewards.json", JSON.stringify(rewards, null, 2));
}

async function start() {
  await createAssetRewardsFile();
}

start();
