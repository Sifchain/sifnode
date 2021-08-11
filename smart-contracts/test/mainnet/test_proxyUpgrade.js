require("dotenv").config();
const Web3 = require("web3");
const fs = require('fs');
const path = require('path');

const { ethers, upgrades } = require("hardhat");

const state = {};

async function setup() {
  setupEnv();
  await setupNewContracts();
  await setupDeployedContracts();
}

function setupEnv() {
  state.provider = new Web3.providers.HttpProvider(process.env.LOCAL_PROVIDER);
  state.chainId = 31337;
  state.web3 = new Web3(state.provider);
  state.owner = process.env.OWNER;
  state.pauser = process.env.PAUSER;
  state.operator = process.env.OPERATOR;
  state.mnemonic = process.env.MNEMONIC;

  state.consensusThreshold = process.env.CONSENSUS_THRESHOLD;
  state.initialValidatorAddresses = process.env.INITIAL_VALIDATOR_ADDRESSES.split(",");
  state.initialValidatorPowers = process.env.INITIAL_VALIDATOR_POWERS.split(",");
  state.mainnetGasPrice = process.env.MAINNET_GAS_PRICE;
  state.erowanAddress = process.env.EROWAN_ADDRESS;
  state.ethereumPrivateKey = process.env.ETHEREUM_PRIVATE_KEY;
}

async function setupNewContracts() {
  state.BridgeBank = {
    new: {
      contract: await ethers.getContractFactory("BridgeBank")
    }
  }

  state.CosmosBridge = {
    new: {
      contract: await ethers.getContractFactory("CosmosBridge")
    }
  }
}

async function setupDeployedContracts() {
  await setDeployedContract('BridgeBank');
  await setDeployedContract('CosmosBridge');
}

async function setDeployedContract(contractName) {
  const filename = path.join(process.cwd(), `deployments/sifchain/${contractName}.json`);
  const artifactContents = fs.readFileSync(filename, { encoding: "utf-8" });
  const parsedArtifactContents = JSON.parse(artifactContents);

  state[contractName].deployed = {
    address: parsedArtifactContents.networks['1'].address,
    abi: parsedArtifactContents.abi,
  }

  const contract = await ethers.getContractAt(parsedArtifactContents.abi, state[contractName].deployed.address);
  state[contractName].deployed.contract = contract;

  console.log(`${contractName} is deployed at address: ${state[contractName].deployed.address}`);
}

describe("Mainnet Upgrade Test", function () {
  before(async function () {
    await setup();
  });

  it("should allow us to prepare the update without any errors", async function () {
    await upgrades.prepareUpgrade(state.BridgeBank.deployed.address, state.BridgeBank.new.contract);
    const operator = await state.BridgeBank.deployed.contract.operator();
    const owner = await state.BridgeBank.deployed.contract.owner();
    console.log(`operator: ${operator}`);
    console.log(`owner: ${owner}`);
  });
});