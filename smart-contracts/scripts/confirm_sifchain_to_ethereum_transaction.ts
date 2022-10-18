import * as hardhat from "hardhat";
import { container } from "tsyringe";
import { DeployedCosmosBridge } from "../src/contractSupport";
import { HardhatRuntimeEnvironmentToken } from "../src/tsyringe/injectionTokens";
import { setupDeployment } from "../src/hardhatFunctions";
import { LogNewOracleClaimEvent } from "../build/CosmosBridge";
import { BigNumber } from "ethers";

const startingBlockNumber = 15676218;
const endingBlockNumber = startingBlockNumber + 1000000;

interface ProphecyClaimState {
  [_index: string]: {
    count: number,
    complete: boolean,
    validators: string[]
  };
}

async function main() {
  container.register(HardhatRuntimeEnvironmentToken, { useValue: hardhat });

  await setupDeployment(container);

  const cosmosBridge = await container.resolve(DeployedCosmosBridge).contract;

  const validatorToBlock: any = {};

  const newProphecyClaims = await cosmosBridge.queryFilter(cosmosBridge.filters.LogNewOracleClaim(), startingBlockNumber, endingBlockNumber);
  const prophecyClaimState = newProphecyClaims.reduce((acc: ProphecyClaimState, x: LogNewOracleClaimEvent) => {
    const prophecyId = BigNumber.from(x.args._prophecyID).toString();
    const item = {
      [prophecyId]: {
        count: ((acc[prophecyId] ?? {}).count ?? 0) + 1,
        validators: ((acc[prophecyId] ?? {}).validators ?? []).concat(x.args._validatorAddress),
        complete: false,
        mostRecentBlockNumber: x.blockNumber
      }
    };
    validatorToBlock[x.args._validatorAddress] = x.blockNumber;
    return {
      ...acc,
      ...item
    };
  }, <ProphecyClaimState>{});

  const prophecyCompletedLogs = await cosmosBridge.queryFilter(cosmosBridge.filters.LogProphecyCompleted(), startingBlockNumber, endingBlockNumber);
  prophecyCompletedLogs.forEach(x => {
    const prophecyId = BigNumber.from((x.args as any)["_prophecyID"]).toString();
    prophecyClaimState[prophecyId].complete = true;
  });

  const activeValidators = Object.entries(prophecyClaimState).reduce((acc, x: any) => {
    const [prophecyId, data] = x;
    for (const v of data["validators"])
      acc.add(v);
    return acc;
  }, new Set());

  console.log(JSON.stringify({
    prophecyClaimState,
    activeValidators: Array.from(activeValidators),
    validatorToBlock
  }, undefined, 2));
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });
