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
      url: "http://localhost:8545/",
      accounts: [
        'ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80', '59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d', '5de4111afa1a4b94908f83103eb1f1706367c2e68ca870fc3fb9a804cdab365a', '7c852118294e51e653712a81e05800f419141751be58f605c371e15141b007a6', '47e179ec197488593b187f80a00eb0da91f1b9d0b13f8733639f19c30a34926a', '8b3a350cf5c34c9194ca85829a2df0ec3153be0318b5e2d3348e872092edffba', '92db14e403b83dfe3df233f83dfa3a0d7096f21ca9b0d6d6b8d88b2b4ec1564e', '4bbbf85ce3377467afe5d46f804f221813b2bb87f24d81f60f1fcdbf7cbf4356', 'dbda1821b80551c9d65939329250298aa3472ba22feea921c0cf5d620ea67b97', '2a871d0798f97d79848a013d4936a73bf4cc922c825d33c1cf7073dff6d409c6'
      ],
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
