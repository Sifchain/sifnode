import * as dotenv from "dotenv"
import { HardhatUserConfig } from "hardhat/config"
import "@nomiclabs/hardhat-ethers"
import "@nomiclabs/hardhat-etherscan"
import '@openzeppelin/hardhat-upgrades'
import "reflect-metadata"; // needed by tsyringe
import "@typechain/hardhat"
import '@float-capital/solidity-coverage'

// require("hardhat-gas-reporter");
// require('hardhat-contract-sizer');

const envconfig = dotenv.config()

const mainnetUrl = process.env["MAINNET_URL"] ?? "https://example.com"
const ropstenUrl = process.env['ROPSTEN_URL'] ?? "https://example.com"

const activePrivateKey = process.env[process.env['ACTIVE_PRIVATE_KEY'] ?? "0xabcd"] ?? "0xabcd";

const runCoverage = !!process.env['RUN_COVERAGE'];

console.log('COVERAGE: ', runCoverage)

const config: HardhatUserConfig = {
  networks: {
    hardhat: {
      allowUnlimitedContractSize: false,
      initialBaseFeePerGas: runCoverage ? 0 : 875000000
    },
    mainnetFork: {
      url: mainnetUrl,
      allowUnlimitedContractSize: false,
      chainId: 1,
      forking: {
        url: mainnetUrl,
        blockNumber: 13200632,
      }
    },
    ropsten: {
      url: ropstenUrl,
      accounts: [activePrivateKey],
      gas: 2000000
    },
    mainnet: {
      url: mainnetUrl,
      accounts: [activePrivateKey],
      gas: 2000000
    }
  },
  solidity: {
    compilers: [
      {
        version: "0.8.0",
        settings: {
          optimizer: {
            enabled: true,
            runs: 200
          },
        }
      },
      {
        version: "0.5.16",
        settings: {
          optimizer: {
            enabled: true,
            runs: 200
          },
        }
      },
    ],
  },
  typechain: {
    outDir: "build",
    target: "ethers-v5"
  },
  mocha: {
    timeout: 200000
  },
  etherscan: {
    // Your API key for Etherscan
    // Obtain one at https://etherscan.io/
    apiKey: process.env['ETHERSCAN_API_KEY']
  }
}

export default config
