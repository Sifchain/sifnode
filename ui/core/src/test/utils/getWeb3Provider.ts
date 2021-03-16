import Web3 from "web3";

import timeMachine from "ganache-time-traveler";

/**
 * Returns a web3 instance that is connected to our test ganache system
 * Also sets up out snapshotting system for tests that use web3
 */
export async function getWeb3Provider() {
  return new Web3.providers.HttpProvider(
    process.env.WEB3_PROVIDER || "http://localhost:7545",
  );
}

let snapshotId: any;

beforeEach(async () => {
  // Unfortunately ganache-time-traveler requires web3 on the global object 🤦‍♂️
  (globalThis as any).web3 = new Web3(await getWeb3Provider());
  let snapshot = await timeMachine.takeSnapshot();
  snapshotId = snapshot["result"];
});

afterEach(async () => {
  await timeMachine.revertToSnapshot(snapshotId);
});
