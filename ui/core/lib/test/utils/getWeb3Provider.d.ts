/**
 * Returns a web3 instance that is connected to our test ganache system
 * Also sets up out snapshotting system for tests that use web3
 */
export declare function getWeb3Provider(): Promise<import("web3-core").HttpProvider>;
