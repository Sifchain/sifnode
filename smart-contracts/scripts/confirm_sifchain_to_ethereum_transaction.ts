import * as hardhat from "hardhat";
import { container } from "tsyringe";
import { DeployedCosmosBridge } from "../src/contractSupport";
import { HardhatRuntimeEnvironmentToken } from "../src/tsyringe/injectionTokens";
import { setupDeployment } from "../src/hardhatFunctions";
import { LogNewOracleClaimEvent } from "../build/CosmosBridge";
import { BigNumber } from "ethers";

const startingBlockNumber = 	15776327;
const endingBlockNumber = startingBlockNumber + 6000;

interface ProphecyClaimState {
  [_index: string]: {
    count: number,
    complete: boolean
  };
}

async function main() {
  container.register(HardhatRuntimeEnvironmentToken, { useValue: hardhat });

  await setupDeployment(container);

  const cosmosBridge = await container.resolve(DeployedCosmosBridge).contract;

  const newProphecyClaims = await cosmosBridge.queryFilter(cosmosBridge.filters.LogNewOracleClaim(), startingBlockNumber, endingBlockNumber);
  const prophecyClaimState = newProphecyClaims.reduce((acc: ProphecyClaimState, x: LogNewOracleClaimEvent) => {
    const prophecyId = BigNumber.from(x.args._prophecyID).toString();
    const item = {
      [prophecyId]: {
        count: ((acc[prophecyId] ?? {}).count ?? 0) + 1,
        complete: false,
        mostRecentBlockNumber: x.blockNumber
      }
    };
    return {
      ...acc,
      ...item
    };
  }, <ProphecyClaimState>{});

  const prophecyCompletedLogs = await cosmosBridge.queryFilter(cosmosBridge.filters.LogProphecyCompleted(), startingBlockNumber, endingBlockNumber);
  prophecyCompletedLogs.forEach(x => {
    const prophecyId = BigNumber.from((x.args as any)["_prophecyID"]).toString();
    prophecyClaimState[prophecyId].complete = true;
  })

  console.log(JSON.stringify({ prophecyClaimState }, undefined, 2));
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });
