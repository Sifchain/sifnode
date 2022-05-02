#!/usr/bin/env node

const fs = require("fs");

async function createAssetRewardsFile(periodId, startBlock, endBlock) {
  const csv = fs.readFileSync(`./${periodId}.csv`, "utf-8");
  const entries = JSON.parse(fs.readFileSync("./entries.json", "utf-8")).result
    .registry.entries;
  const lines = csv.split("\r\n").filter((line) => line.split(",")[1] !== "");

  let [, allocation] = lines[1].split('"');
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

  const rewardPeriod = {
    reward_period_id: periodId,
    reward_period_start_block: startBlock,
    reward_period_end_block: endBlock,
    reward_period_allocation: allocation,
    reward_period_pool_multipliers: multipliers,
    reward_period_default_multiplier: "0.0",
  };

  return rewardPeriod;
}

async function start() {
  const rewardPeriods = [
    await createAssetRewardsFile("RP_2", 6586931, 6687730),
    await createAssetRewardsFile("RP_1", 6486131, 6586930),
  ];

  fs.writeFileSync("./rewards.json", JSON.stringify(rewardPeriods, null, 2));
}

start();
