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

const networkUrl = process.env["NETWORK_URL"] ?? "http://needToSetNETWORK_URL.nothing"
const activePrivateKey = process.env["ETHEREUM_PRIVATE_KEY"] ?? "0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
const keyList = [activePrivateKey]

if (!networkUrl) {
  print("error", "ABORTED! Missing NETWORK_URL env variable")
  throw new Error("INVALID_NETWORK_URL")
}
if (!activePrivateKey) {
  print("error", "ABORTED! Missing ACTIVE_PRIVATE_KEY env variable")
  throw new Error("INVALID_PRIVATE_KEY")
}

const runCoverage = process.env["RUN_COVERAGE"] ? true : false
if (runCoverage) print("warn", "HARDHAT :: Test coverage mode is ON")

const reportGas = process.env["REPORT_GAS"] ? true : false
if (reportGas) print("warn", "HARDHAT :: Gas reporter is ON")

// Works only for 'hardhat' network:
const useForking = process.env["USE_FORKING"] ? true : false
if (useForking) print("warn", "HARDHAT :: Forking is ON")

var accounts: string[] = []
if (process.env["ETH_ACCOUNTS"] ? true : false) {
    accounts = (process.env["ETH_ACCOUNTS"] || "").split(",")
}

const config: HardhatUserConfig = {
  networks: {
    hardhat: {
      allowUnlimitedContractSize: false,
      initialBaseFeePerGas: runCoverage ? 0 : 875000000,
      chainId: 9999,
      mining: {
        auto: true,
        interval: 200,
      },
      forking: {
        enabled: useForking,
        url: networkUrl,
        blockNumber: 13469882,
      },
    },
    ropsten: {
      url: networkUrl,
      accounts: keyList,
      gas: 2000000,
    },
    geth: {
      url: networkUrl,
      accounts: accounts,
      gas: 6000000,
      gasPrice: "auto",
      gasMultiplier: 1.2,
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
        version: "0.8.17",
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
