import '@openzeppelin/hardhat-upgrades'
import "@nomiclabs/hardhat-ethers"
import "@typechain/hardhat"
import "reflect-metadata"; // needed by tsyringe
import {HardhatUserConfig} from "hardhat/config";

const config:HardhatUserConfig ={
    networks: {
        hardhat: {
            allowUnlimitedContractSize: false,
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
