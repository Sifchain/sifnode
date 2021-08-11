import '@openzeppelin/hardhat-upgrades'
import "@nomiclabs/hardhat-ethers"
import "@typechain/hardhat"
import "reflect-metadata"; // needed by tsyringe
import { HardhatUserConfig } from "hardhat/config";
require('solidity-coverage');
require("hardhat-gas-reporter");
require('hardhat-contract-sizer');

const config: HardhatUserConfig = {
  networks: {
    hardhat: {
      allowUnlimitedContractSize: false,
    },
    localRpc: {
      allowUnlimitedContractSize: false,
      chainId: 31337,
      url: 'http://127.0.0.1:8545/',
    },
  },
  solidity: {
    version: "0.8.0",
    settings: {
      optimizer: {
        enabled: true,
        runs: 200
      },
    },
  },
  typechain: {
    outDir: "build",
    target: "ethers-v5"
  },
}

export default config
