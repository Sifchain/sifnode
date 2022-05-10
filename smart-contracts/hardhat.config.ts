import * as dotenv from "dotenv"
import { HardhatUserConfig } from "hardhat/config"
import "@nomiclabs/hardhat-ethers"
import "@nomiclabs/hardhat-etherscan"
import "@openzeppelin/hardhat-upgrades"
import "@float-capital/solidity-coverage"
import "hardhat-contract-sizer"
import "hardhat-gas-reporter"
import "reflect-metadata" // needed by tsyringe
import "@typechain/hardhat"
import "@nomiclabs/hardhat-waffle";

import "./tasks/task_blocklist";


import { print } from "./scripts/helpers/utils";
import { parseEther } from "ethers/lib/utils"

const networkUrl = process.env["NETWORK_URL"] ?? "http://needToSetNETWORK_URL.nothing";
const privateKey = process.env["PRIVATE_KEY"] ?? "0xe749379ae9430aee3494464263e0a6f66a3dde9d64eaf134d1b9990c1c006b0f";
const adminPrivateKey = process.env["ADMIN_KEY"] ?? "0x23007009f1ee4d6b5213dc6cfc7bf29c70694b94cf31c53a93af805f84e6eba0";
const operatorPrivateKey = process.env["OPERATOR_KEY"] ?? "0x81503fee57667e46debf70ece39326a0d36daeaa256af5e48f6e5578cfb2c16d";
const pauserPrivateKey = process.env["PAUSER_KEY"] ?? "0x04e31b83fdec75f8f845ec3390a9b02d056078044cdc6db1d5c14641d0299cca";
const keyList = [privateKey, adminPrivateKey, operatorPrivateKey, pauserPrivateKey];
const forkingAccounts = Number(process.env["IMPERSONATE_ACCOUNTS"]) == 1 ? keyList.map((key) => ({privateKey: key, balance: "0"})) : undefined;
const runCoverage = process.env["RUN_COVERAGE"] ? true : false
if (runCoverage) print("warn", "HARDHAT :: Test coverage mode is ON")

const reportGas = process.env["REPORT_GAS"] ? true : false
if (reportGas) print("warn", "HARDHAT :: Gas reporter is ON")

// Works only for 'hardhat' network:
const useForking = process.env["USE_FORKING"] ? true : false
if (useForking) print("warn", "HARDHAT :: Forking is ON")

const config: HardhatUserConfig = {
  networks: {
    hardhat: {
      allowUnlimitedContractSize: false,
      initialBaseFeePerGas: runCoverage ? 0 : 875000000,
      chainId: 1,
      mining: {
        auto: true,
        interval: 200,
      },
      accounts: forkingAccounts,
      forking: {
        enabled: useForking,
        url: networkUrl,
        blockNumber: 14630403,
      },
    },
    ropsten: {
      url: networkUrl,
      accounts: keyList,
      gas: 2000000,
    },
    mainnet: {
      url: networkUrl,
      accounts: keyList,
      gas: 6000000,
      gasPrice: "auto",
      gasMultiplier: 1.2,
    },
  },
  solidity: {
    compilers: [
      {
        version: "0.8.4",
        settings: {
          optimizer: {
            enabled: true,
            runs: 200,
          },
        },
      },
      {
        version: "0.5.16",
        settings: {
          optimizer: {
            enabled: true,
            runs: 200,
          },
        },
      },
    ],
  },
  typechain: {
    outDir: "build",
    target: "ethers-v5",
  },
  mocha: {
    timeout: 200000,
  },
  etherscan: {
    // Your API key for Etherscan
    // Obtain one at https://etherscan.io/
    apiKey: process.env["ETHERSCAN_API_KEY"],
  },
  paths: {
    sources: "./contracts",
    tests: "./test",
    cache: "./cache",
    artifacts: "./artifacts",
  },
  gasReporter: {
    enabled: reportGas,
  },
}

export default config
