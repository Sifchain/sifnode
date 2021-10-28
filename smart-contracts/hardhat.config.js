require("hardhat/config");
require("@nomiclabs/hardhat-ethers");
require("@nomiclabs/hardhat-etherscan");
require("@openzeppelin/hardhat-upgrades");
require("@float-capital/solidity-coverage");
require("hardhat-contract-sizer");
require("hardhat-gas-reporter");
require("dotenv").config();

const { print } = require("./scripts/helpers/utils");

const networkUrl = process.env["NETWORK_URL"];
const activePrivateKey = process.env[process.env.ACTIVE_PRIVATE_KEY];

if (!networkUrl) {
  print("error", "ABORTED! Missing NETWORK_URL env variable");
  throw new Error("INVALID_NETWORK_URL");
}
if (!activePrivateKey) {
  print("error", "ABORTED! Missing ACTIVE_PRIVATE_KEY env variable");
  throw new Error("INVALID_PRIVATE_KEY");
}

const runCoverage = process.env["RUN_COVERAGE"] ? true : false;
if (runCoverage) print("warn", "HARDHAT :: Test coverage mode is ON");

const reportGas = process.env["REPORT_GAS"] ? true : false;
if (reportGas) print("warn", "HARDHAT :: Gas reporter is ON");

// Works only for 'hardhat' network:
const useForking = process.env["USE_FORKING"] ? true : false;
if (useForking) print("warn", "HARDHAT :: Forking is ON");

module.exports = {
  networks: {
    hardhat: {
      allowUnlimitedContractSize: false,
      chainId: 1,
      initialBaseFeePerGas: runCoverage ? 0 : 875000000,
      forking: {
        enabled: useForking,
        url: networkUrl,
        blockNumber: 13469882,
      },
    },
    ropsten: {
      url: networkUrl,
      accounts: [activePrivateKey],
      gas: 2000000,
    },
    mainnet: {
      url: networkUrl,
      accounts: [activePrivateKey],
      gas: 2000000,
      gasPrice: "auto",
      gasMultiplier: 1.2,
    },
  },
  solidity: {
    compilers: [
      {
        version: "0.8.0",
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
  mocha: {
    timeout: 20000,
  },
  gasReporter: {
    enabled: reportGas,
  },
};
