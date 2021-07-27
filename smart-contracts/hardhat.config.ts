import '@openzeppelin/hardhat-upgrades'
import "@nomiclabs/hardhat-ethers"
import "@typechain/hardhat"
import "reflect-metadata"; // needed by tsyringe
import {HardhatUserConfig} from "hardhat/config";
// require('solidity-coverage');
// require("hardhat-gas-reporter");
// require('hardhat-contract-sizer');

const alchemyUrl = process.env["ALCHEMY_URL"] ?? "no url specified"

const config: HardhatUserConfig = {
    networks: {
        hardhat: {
            allowUnlimitedContractSize: false,
            chainId: 1,
            forking: {
                url: alchemyUrl,
                blockNumber: 12865480,
            }
        },
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
    }
}

export default config
