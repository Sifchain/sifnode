/**
 * This script will upgrade BridgeBank and CosmosBridge in production and connect it to the Blocklist.
 * For instructions, please consult scripts/upgrades/peggy2-2022-02/runbook.md
 */
require("dotenv").config();

import hardhat from "hardhat";
import fs from "fs-extra";
import Web3 from "web3";
import support, { getDeployedContract } from "../../helpers/forkingSupport";
import { print } from "../../helpers/utils";
import { ethers } from "hardhat";
import { Blocklist, BridgeBank, CosmosBridge } from "../../../build";
import { BigNumber } from "ethers";
import { state } from "fp-ts";
import { SignerWithAddress } from "@nomiclabs/hardhat-ethers/signers";

// Helps converting an address to a checksum address
const addr = Web3.utils.toChecksumAddress;

// Defaults to the Ethereum Mainnet address
const BRIDGEBANK_ADDRESS = process.env.BRIDGEBANK_ADDRESS || addr("0xB5F54ac4466f5ce7E0d8A5cB9FE7b8c0F35B7Ba8");

// Defaults to empty string if not set
const BLOCKLIST_ADDRESS = process.env.BLOCKLIST_ADDRESS || ""

// If there is no DEPLOYMENT_NAME env var, we'll use the mainnet deployment
const DEPLOYMENT_NAME = process.env.DEPLOYMENT_NAME || "sifchain-1";

// If there is no FORKING_CHAIN_ID env var, we'll use the mainnet id
const CHAIN_ID = Number(process.env.CHAIN_ID) || 1;

// Are we running this script in test mode?
const USE_FORKING = !!process.env.USE_FORKING;
const IMPERSONATE_ACCOUNTS = !!process.env.IMPERSONATE_ACCOUNTS;

interface BridgeBankStorageSlots {
  pauser1: boolean;
  pauser2: boolean;
  owner: string;
  nonce: BigNumber;
  cosmosWhitelist: boolean;
  blocklist?: string,
};

interface CosmosBridgeStorageSlots {
  bridgeBank: string,
  consensusThreshold: BigNumber,
  currentValsetVersion: BigNumber,
  hasBridgeBank: boolean,
  totalPower: BigNumber,
  validatorCount: BigNumber,
  validator1IsValidator: boolean,
  validator2IsValidator: boolean,
  validator3IsValidator: boolean,
  validator4IsValidator: boolean,
  validator1Power: BigNumber,
  validator2Power: BigNumber,
  validator3Power: BigNumber,
  validator4Power: BigNumber,
  operator: string, 
};

interface StorageSlots extends BridgeBankStorageSlots, CosmosBridgeStorageSlots {};

type StorageSlotsKeys = keyof StorageSlots;

interface StateAddresses {
    pauser1: string,
    pauser2: string,
    validator1: string,
    validator2: string,
    validator3: string,
    validator4: string,
    cosmosWhitelistedtoken: string
}

interface StateSigners {
  admin: SignerWithAddress,
  operator: SignerWithAddress,
  pauser: SignerWithAddress 
}

interface StateContracts {
  bridgeBank: BridgeBank,
  cosmosBridge: CosmosBridge,
  blocklist: Blocklist,
  upgradedBridgeBank?: BridgeBank,
  upgradedCosmosBridge?: CosmosBridge
}

interface State {
  addresses: StateAddresses,
  signers: StateSigners,
  contracts: StateContracts,
  storageSlots: {
    before: StorageSlots,
    after: StorageSlots
  }
}

async function main() {
  print("highlight", "~~~ UPGRADE FROM PEGGY 1.0 TO PEGGY 2.0 ~~~");

  // We heeded the warnings and made sure the upgrade will work
  hardhat.upgrades.silenceWarnings();

  const addresses : StateAddresses = {
    pauser1: addr("c0a586fb260b2c14098a9d95b75f56f13cad2dd9"),
    pauser2: addr("0x9910ade709043d8b9ed2a31fdfcbfb6538f9a397"),
    validator1: addr("0x0D7dEF5C00a8B6ddc58A0255f0a94cc739C6d0B5"),
    validator2: addr("0x9B4002670C210A3b64e13807250BE62B8dEae201"),
    validator3: addr("0xbF45BFc92ebD305d4C0baf8395c4299bdFCE9EA2"),
    validator4: addr("0xeB29C7016eDd2D6B413fceE4C51474FED058005a"),
    cosmosWhitelistedtoken: addr("0x8D983cb9388EaC77af0474fA441C4815500Cb7BB"),
  }

  // Fetch the manifest and inject the new variables
  copyManifest(false);

  // Connect to each contract
  let contracts = await connectToContracts();

  // Get signers to send transactions
  const signers = await setupAccounts(addresses, contracts);

  // Fetch current values from the deployed contract
  const beforeStorageSlots: StorageSlots =  {
    ...await setBridgeBankStorageSlots(addresses, contracts.bridgeBank),
    ...await setCosmosBridgeStorageSlots(addresses, contracts.cosmosBridge),
  }

  // Pause the system
  await pauseBridgeBank(contracts.bridgeBank, signers);

  // Upgrade BridgeBank
  contracts.upgradedBridgeBank = await upgradeBridgeBank(contracts.bridgeBank, signers);

  // Upgrade CosmosBridge
  contracts.upgradedCosmosBridge = await upgradeCosmosBridge(contracts.cosmosBridge, signers);

  // Fetch values after the upgrade
  const afterStorageSlots: StorageSlots = {
    ...await setBridgeBankStorageSlots(addresses, contracts.upgradedBridgeBank),
    ...await setCosmosBridgeStorageSlots(addresses, contracts.upgradedCosmosBridge),
  }

  // Compare slots before and after the upgrade
  checkStorageSlots(beforeStorageSlots, afterStorageSlots);

  // Resume the system
  await resumeBridgeBank(contracts.bridgeBank, signers);

  // Clean up temporary files
  // cleanup();

  print("highlight", "~~~ DONE! 👏 Everything worked as expected. ~~~");

  if (!USE_FORKING) {
    print("h_green", `Peggy 1.0 has been upgraded to Peggy 2.0 on Ethereum`);
  }
}

async function setBridgeBankStorageSlots(addresses: StateAddresses, contract: BridgeBank) : Promise<BridgeBankStorageSlots> {
  return {
    pauser1: await contract.pausers(addresses.pauser1),
    pauser2: await contract.pausers(addresses.pauser2),
    owner: await contract.owner(),
    nonce: await contract.lockBurnNonce(),
    cosmosWhitelist: await contract.getCosmosTokenInWhiteList(
      addresses.cosmosWhitelistedtoken
    ),
  }
}

async function setCosmosBridgeStorageSlots(addresses: StateAddresses, contract: CosmosBridge): Promise<CosmosBridgeStorageSlots> {
  return {
    bridgeBank: await contract.bridgeBank(),
    consensusThreshold: await contract.consensusThreshold(),
    currentValsetVersion: await contract.currentValsetVersion(),
    hasBridgeBank: await contract.hasBridgeBank(),
    totalPower: await contract.totalPower(),
    validatorCount: await contract.validatorCount(),
    validator1IsValidator: await contract.isActiveValidator(addresses.validator1),
    validator2IsValidator: await contract.isActiveValidator(addresses.validator2),
    validator3IsValidator: await contract.isActiveValidator(addresses.validator3),
    validator4IsValidator: await contract.isActiveValidator(addresses.validator4),
    validator1Power: await contract.getValidatorPower(addresses.validator1),
    validator2Power: await contract.getValidatorPower(addresses.validator2),
    validator3Power: await contract.getValidatorPower(addresses.validator3),
    validator4Power: await contract.getValidatorPower(addresses.validator4),
    operator: await contract.operator(),
  }
}

function checkStorageSlots(before: StorageSlots, after: StorageSlots) {
  print("yellow", "🎯 Checking storage layout");

  const keys = Object.keys(before) as StorageSlotsKeys[];

  keys.forEach((key) => {
    testMatch(before[key], after[key], key);
  });
}

function testMatch(before: unknown, after: unknown, slotName: string) {
  if (String(before) === String(after)) {
    print("green", `✅ ${slotName} slot is safe`);
  } else {
    throw new Error(`💥 CRITICAL: ${slotName} Mismatch! From ${before} to ${after}`);
  }
}

async function connectToContracts() : Promise<StateContracts> {
  print("yellow", `🕑 Connecting to contracts...`);
  
  // Create an instance of BridgeBank from the deployed code
  const { contract: bridgeBank } = await support.getDeployedContract<BridgeBank>(
    DEPLOYMENT_NAME,
    "BridgeBank",
    CHAIN_ID
  );

  // Get the cosmosBridgeAddress and Blocklist Address
  const cosmosBridgeAddress = await bridgeBank.cosmosBridge();
  
  // Create an instance of CosmosBridge
  const { contract: cosmosBridge } = await support.getDeployedContract<CosmosBridge>(
    DEPLOYMENT_NAME,
    "CosmosBridge",
    CHAIN_ID
  );
  
  // Create an instance of the Blocklist
  const blocklistFactory = await ethers.getContractFactory("Blocklist");
  
  let blocklist: Blocklist;
  
  if (BLOCKLIST_ADDRESS !== "") {
    blocklist = await blocklistFactory.attach(BLOCKLIST_ADDRESS);
  } else {
    const blocklistAddress = await bridgeBank.blocklist();
    blocklist = await blocklistFactory.attach(blocklistAddress);
  }
  print("green", `✅ Contracts connected`);
  
  return {
    blocklist: blocklist,
    bridgeBank: bridgeBank,
    cosmosBridge: cosmosBridge
  }
}

async function setupAccounts(addresses: StateAddresses, contracts: StateContracts) : Promise<StateSigners> {
  const operatorAddress = await contracts.bridgeBank.operator();

  let admin: SignerWithAddress;
  let operator: SignerWithAddress;
  let pauser: SignerWithAddress;

  // If we're forking, we want to impersonate the owner account
    if(IMPERSONATE_ACCOUNTS) {
      print("magenta", "MAINNET FORKING :: IMPERSONATE ACCOUNT");
  
      admin = await support.impersonateAccount(
        support.PROXY_ADMIN_ADDRESS,
        "10000000000000000000",
        "Proxy Admin"
      );
  
      operator = await support.impersonateAccount(
        operatorAddress,
        "10000000000000000000",
        "Operator"
      );
  
      pauser = await support.impersonateAccount(
        addresses.pauser1,
        "10000000000000000000",
        "Pauser"
      );
    } else {
      print("magenta", "MAINNET FORKING :: NOT IMPERSONATING ACCOUNT");
    // If not, we want the connected accounts
    const accounts = await ethers.getSigners();
    admin = accounts[1];
    operator = accounts[2];
    pauser = accounts[3];
  }

  const hasCorrectAdmin =
    admin.address.toLowerCase() === support.PROXY_ADMIN_ADDRESS.toLowerCase();

  const hasCorrectOperator =
    operator.address.toLowerCase() === operatorAddress.toLowerCase();

  const hasCorrectPauser =
    pauser.address.toLowerCase() === addresses.pauser1.toLowerCase() ||
    pauser.address.toLowerCase() === addresses.pauser2.toLowerCase();

  if (!hasCorrectAdmin) {
    throw new Error(
      `The first Private Key is not the PROXY ADMIN's private key. Please use the Private Key that corresponds to the address ${support.PROXY_ADMIN_ADDRESS}`
    );
  }

  if (!hasCorrectOperator) {
    throw new Error(
      `The second Private Key is not the BridgeBank OPERATOR's private key. Please use the Private Key that corresponds to the address ${operatorAddress}`
    );
  }

  if (!hasCorrectPauser) {
    throw new Error(
      `The third Private Key is not the BridgeBank PAUSER's private key. Please use the Private Key that corresponds to the address ${addresses.pauser1} or ${addresses.pauser2}`
    );
  }

  const adminColor = hasCorrectAdmin ? "white" : "red";
  const operatorColor = hasCorrectOperator ? "white" : "red";
  const pauserColor = hasCorrectPauser ? "white" : "red";

  print(adminColor, `🤵 ProxyAdmin: ${admin.address}`);
  print(operatorColor, `🤵 Operator: ${operator.address}`);
  print(pauserColor, `🤵 Pauser: ${pauser.address}`);

  return {
    admin,
    operator,
    pauser
  }
}

async function pauseBridgeBank(bridgeBank: BridgeBank, signers: StateSigners) {
  print(
    "yellow",
    `🕑 Pausing the system before the upgrade. Please wait, this may take a while...`
  );
  await bridgeBank.connect(signers.pauser).pause();
  print("green", `✅ System is paused`);
}

async function resumeBridgeBank(bridgeBank: BridgeBank, signers: StateSigners) {
  print("yellow", `🕑 Unpausing the system. Please wait, this may take a while...`);
  await bridgeBank.connect(signers.pauser).unpause();
  print("green", `✅ System has been resumed`);
}

async function upgradeBridgeBank(bridgeBank: BridgeBank, signers: StateSigners): Promise<BridgeBank> {
  print("yellow", `🕑 Upgrading BridgeBank. Please wait, this may take a while...`);
  const newBridgeBankFactory = await hardhat.ethers.getContractFactory("BridgeBank");
  const bridgeBankUpgrade = await hardhat.upgrades.upgradeProxy(
    bridgeBank,
    newBridgeBankFactory.connect(signers.admin),
    // TODO: Why is this necessary? Our logic contracts do not have delegatecall in them.
    { unsafeAllow: ["delegatecall"] }
  ) as BridgeBank;
  await bridgeBankUpgrade.deployed();
  print("green", `✅ BridgeBank Upgraded`);
  return bridgeBankUpgrade;
}

async function upgradeCosmosBridge(cosmosBridge: CosmosBridge, signers: StateSigners): Promise<CosmosBridge> {
  print("yellow", `🕑 Upgrading CosmosBridge. Please wait, this may take a while...`);
  const newCosmosBridgeFactory = await hardhat.ethers.getContractFactory("CosmosBridge");
  const cosmosBridgeUpgrade = await hardhat.upgrades.upgradeProxy(
    cosmosBridge,
    newCosmosBridgeFactory.connect(signers.admin),
    { unsafeAllow: ["delegatecall"] }
  ) as CosmosBridge;
  await cosmosBridgeUpgrade.deployed();
  print("green", `✅ CosmosBridge Upgraded`);
  return cosmosBridgeUpgrade;
}

// Copy the manifest to the right place (where Hardhat wants it)
function copyManifest(injectChanges: boolean) {
  print("cyan", `👀 Fetching the correct manifest`);

  if (!injectChanges) {
    // just copy the file over to the correct directory
    fs.copySync(
      `./deployments/sifchain-1/.openzeppelin/mainnet.json`,
      `./.openzeppelin/mainnet.json`
    );
  } else {
    // will write the file into the correct directory at the end
    injectStorageChanges();
  }
}

// https://forum.openzeppelin.com/t/storage-layout-upgrade-with-hardhat-upgrades/14567
// All changes made here affect only deprecated variables;
// The injection is done so that OZ's lib doesn't complain about type changes;
// The specific changes we made are safe;
function injectStorageChanges() {
  print("cyan", "🕵  Injecting changes into manifest");

  // Fetch the deployed manifest
  const currentManifest = fs.readFileSync(
    "./deployments/sifchain-1/.openzeppelin/mainnet.json",
    "utf8"
  );

  // Parse the deployed manifest
  const parsedManifest = JSON.parse(currentManifest);

  // Change variable types
  const modManifest = support.replaceTypesInManifest({
    parsedManifest,
    originalType: "t_string_memory",
    newType: "t_string_memory_ptr",
  });

  // Write to file
  fs.writeFileSync("./.openzeppelin/mainnet.json", JSON.stringify(modManifest));
}

// Delete temporary files (the copied manifest)
function cleanup() {
  print("cyan", `🧹 Cleaning up temporary files`);

  fs.unlinkSync(`./.openzeppelin/mainnet.json`);
}

main()
  .catch((error) => {
    print("h_red", error.stack);
  })
  .finally(() => process.exit(0));
