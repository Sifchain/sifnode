const getEnv = require("./helpers/envLoader").loadEnv;

const hre = require("hardhat");
const fs = require("fs");
const { ethers, upgrades } = require("hardhat");

require("dotenv").config();

async function main() {
    const accounts = await ethers.getSigners();

    const BridgeBank = await hre.ethers.getContractFactory("BridgeBank");
    const CosmosBridge = await hre.ethers.getContractFactory("CosmosBridge");
    const BridgeRegistry = await hre.ethers.getContractFactory("BridgeRegistry");

    const {
        consensusThreshold,
        operator,
        initialValidators,
        initialPowers,
        owner,
        pauser
    } = getEnv();

    const cosmosBridge = await upgrades.deployProxy(CosmosBridge, [
        operator,
        consensusThreshold,
        initialValidators,
        initialPowers
    ]);

    await cosmosBridge.deployed();
    console.log("deployed cosmos bridge to: ", cosmosBridge.address);

    const bank = await upgrades.deployProxy(BridgeBank, [
        operator,
        cosmosBridge.address,
        owner,
        pauser,
    ]);

    await bank.deployed();
    console.log("Bridge bank deployed to:", bank.address);
    
    const bridgeRegistry = await upgrades.deployProxy(BridgeRegistry, [
        cosmosBridge.address,
        bank.address
    ]);

    await bridgeRegistry.deployed();
    console.log("Bridge registry deployed to:", bridgeRegistry.address);
    
    fs.writeFileSync("mainnet.json", JSON.stringify({
        bridgebank: bank.address,
        cosmosbridge: cosmosBridge.address,
        bridgeregistry: bridgeRegistry.address,
    }, null, 4));
}

main();