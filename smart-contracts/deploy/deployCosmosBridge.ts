import {DeployFunction} from 'hardhat-deploy/types';
import {HardhatRuntimeEnvironment} from 'hardhat/types';
import '@openzeppelin/hardhat-upgrades';

import {ContractNames} from "../src/contractNames";
import {loadDeploymentEnvWithDotenv} from "../src/deploymentEnv";

const deployCosmosBridge: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
    console.log("in deployCosmosBridge, " + JSON.stringify({

    }))
    const CosmosBridge = await hre.ethers.getContractFactory("CosmosBridge")
    const {operator} = await hre.getNamedAccounts()
    const {validator1} = await hre.getNamedAccounts()
    const deploymentEnv = loadDeploymentEnvWithDotenv()
    let cosmosBridgeArgs = [
        operator,
        deploymentEnv.consensusThreshold,
        [validator1],
        deploymentEnv.initialPowers
    ];
    console.log("in deployCosmosBridge, " + JSON.stringify({
        operator,
        args: cosmosBridgeArgs
    }, undefined, 2))
    const cosmosBridge = await hre.upgrades.deployProxy(CosmosBridge, cosmosBridgeArgs);
    await cosmosBridge.deployed()
    console.log("deployed cosmos bridge to: ", cosmosBridge.address, cosmosBridgeArgs);
};

deployCosmosBridge.tags = [ContractNames.CosmosBridge]

export default deployCosmosBridge;