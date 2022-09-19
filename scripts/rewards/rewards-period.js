import fs from "fs";
import fetch from "node-fetch";

function formatAllocationAmount(allocation) {
  return `${allocation.trim().split(",").join("")}${"0".repeat(18)}`;
}
function doesPoolExist(base_denom, pool) {
  return (
    base_denom === `c${pool}` ||
    base_denom === `u${pool}` ||
    base_denom === `e${pool}` ||
    base_denom === `a${pool}` ||
    base_denom === `${pool}`
  );
}

async function createAssetRewardsFileFromParams({
  periodId,
  startBlock,
  endBlock,
  allocation,
  defaultMultiplier = "0.0",
  pools = [],
}) {
  const entries = (
    await (
      await fetch("https://api.sifchain.finance/tokenregistry/entries")
    ).json()
  ).result.registry.entries;

  const multipliers = pools.map(({ pool, multiplier }) => {
    const entry = entries.find(({ base_denom }) =>
      doesPoolExist(base_denom, pool)
    );

    const asset = entry ? entry.denom : `???${pool}???`;

    return {
      pool_multiplier_asset: asset,
      multiplier,
    };
  });

  const rewardPeriod = {
    reward_period_id: periodId,
    reward_period_start_block: startBlock,
    reward_period_end_block: endBlock,
    reward_period_allocation: formatAllocationAmount(allocation),
    reward_period_pool_multipliers: multipliers,
    reward_period_default_multiplier: defaultMultiplier,
  };

  return rewardPeriod;
}

async function createAssetRewardsFile(periodId, startBlock, endBlock) {
  const csv = fs.readFileSync(`./${periodId}.csv`, "utf-8");
  const entries = JSON.parse(fs.readFileSync("./entries.json", "utf-8")).result
    .registry.entries;
  const lines = csv.split("\r\n").filter((line) => line.split(",")[0] !== "");

  let [, allocation] = lines[0].split('"');

  const multipliers = lines.slice(1).map((line) => {
    const [, poolName, multiplier] = line.split(",");

    const entry = entries.find(({ base_denom }) =>
      doesPoolExist(base_denom, poolName)
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
    reward_period_allocation: formatAllocationAmount(allocation),
    reward_period_pool_multipliers: multipliers,
    reward_period_default_multiplier: "0.0",
  };

  return rewardPeriod;
}

async function start() {
  const rewardPeriods = [
    await createAssetRewardsFileFromParams({
      periodId: "RP_7",
      startBlock: 7081749,
      endBlock: 7180679,
      allocation: "32938034",
      defaultMultiplier: "0.585",
      pools: [
        { pool: "atom", multiplier: "1.085" },
        { pool: "eth", multiplier: "1.085" },
        { pool: "usdc", multiplier: "1.085" },
        { pool: "juno", multiplier: "1.085" },
        { pool: "uscrt", multiplier: "1.085" },
        { pool: "luna", multiplier: "0.0" },
        { pool: "ust", multiplier: "0.0" },
        { pool: "usd", multiplier: "0.0" },
      ],
    }),
    await createAssetRewardsFileFromParams({
      periodId: "RP_6",
      startBlock: 6982961,
      endBlock: 7081748,
      allocation: "35040462",
      defaultMultiplier: "0.602",
      pools: [
        { pool: "atom", multiplier: "1.075" },
        { pool: "eth", multiplier: "1.075" },
        { pool: "usdc", multiplier: "1.075" },
        { pool: "juno", multiplier: "1.075" },
        { pool: "uscrt", multiplier: "1.075" },
        { pool: "luna", multiplier: "0.0" },
        { pool: "ust", multiplier: "0.0" },
        { pool: "usd", multiplier: "0.0" },
      ],
    }),
    // await createAssetRewardsFile("RP_5", 6885850, 6982960),
    // await createAssetRewardsFile("RP_4", 6788531, 6885841),
    // await createAssetRewardsFile("RP_3", 6687731, 6788530),
    // await createAssetRewardsFile("RP_2", 6586931, 6687730),
    // await createAssetRewardsFile("RP_1", 6486131, 6586930),
  ];

  fs.writeFileSync("./rewards.json", JSON.stringify(rewardPeriods, null, 2));
}

start();
