import Web3 from "web3";

import timeMachine from "ganache-time-traveler";

/**
 * Returns a web3 instance that is connected to our test ganache system
 * Also sets up out snapshotting system for tests that use web3
 */
export async function getWeb3() {
  if (web3) return web3;
  runBeforeEachHooks();
  return new Web3(new Web3.providers.HttpProvider("http://localhost:8545"));
}

let web3: Web3;
let snapshotId: any;

function runBeforeEachHooks() {
  beforeEach(async () => {
    (globalThis as any).web3 = await getWeb3();
    let snapshot = await timeMachine.takeSnapshot();
    snapshotId = snapshot["result"];
  });

  afterEach(async () => {
    await timeMachine.revertToSnapshot(snapshotId);
  });
}
