import { task } from "hardhat/config";
import '@openzeppelin/hardhat-upgrades'
import "@nomiclabs/hardhat-waffle"
import "hardhat-deploy"

// require('hardhat-local-networks-config-plugin')
// require("hardhat-typechain");

// This is a sample Hardhat task. To learn how to create your own go to
// https://hardhat.org/guides/create-task.html
task("accounts", "Prints the list of accounts", async (args, hre) => {
    const accounts = await hre.ethers.getSigners();

    for (const account of accounts) {
        console.log(await account.address);
    }
});

/**
 * @type import('hardhat/config').HardhatUserConfig
 */
module.exports = {
    localNetworksConfig: '~/.hardhat/networks.ts',
    networks: {
        hardhat: {
            allowUnlimitedContractSize: false,
        },
    },
    namedAccounts: {
        owner: {
            hardhat: 0
        },
        operator: {
            hardhat: 1
        },
        pauser: {
            hardhat: 2
        },
        validator1: {
            hardhat: 3
        }
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
};
