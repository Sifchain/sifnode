import Web3 from "web3";

/**
 * Returns a web3 instance that is connected to our test ganache system
 * Also sets up out snapshotting system for tests that use web3
 */
export async function getWeb3Provider() {
  return new Web3.providers.HttpProvider(
    process.env.WEB3_PROVIDER || "http://localhost:7545",
  );
}
