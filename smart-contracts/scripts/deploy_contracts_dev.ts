/**
 * !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
 * THIS SCRIPT IS FOR DEPLOYING CONTRACTS IN TESTNETS AND LOCAL GETH INSTANCES ONLY
 * DO NOT USE IN PRODUCTION, ALL PRODUCTION KEYS NEED TO SUPPORT HARDWARE WALLETS AND
 * GNOSIS CONTRACTS....
 * !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
 **/

import * as dotenv from "dotenv"
import hardhat, { ethers, upgrades } from "hardhat";
import { SignerWithAddress } from "@nomiclabs/hardhat-ethers/signers"
import { BridgeBank, CosmosBridge } from "../build";
import {print} from "./helpers/utils"
import fs from "fs-extra";

export interface DeployedContractAddresses {
  blocklist: string
  cosmosBridge: string
  bridgeBank: string
  bridgeRegistry: string
  rowanContract: string
}

export interface SifchainAccounts {
  readonly operatorAccount: SignerWithAddress,
  readonly ownerAccount: SignerWithAddress,
  readonly pauserAccount: SignerWithAddress,
  readonly validatorAccounts: SignerWithAddress[],
  readonly availableAccounts: SignerWithAddress[]
};

async function hreToSifchainAccountsAsync(): Promise<SifchainAccounts> {
  const accounts = await hardhat.ethers.getSigners()
  // Keep this synched with run_env.py
  const [operatorAccount, ownerAccount, pauserAccount, validator1Account, ...extraAccounts] =
    accounts
  return {
    operatorAccount,
    ownerAccount,
    pauserAccount,
    validatorAccounts: [validator1Account],
    availableAccounts: extraAccounts
  }
}

const NETWORK_DESCRIPTOR = Number(process.env.NETWORK_DESCRIPTOR) || 9999;

// Delete temporary files (the copied manifest)
function cleanup() {
  print("cyan", `🧹 Cleaning up temporary files`);

  fs.removeSync(`./.openzeppelin/unknown-${NETWORK_DESCRIPTOR}.json`);
}


async function main() : Promise<DeployedContractAddresses> {
  print("warn", "THIS IS A DEVELOPMENT ONLY SCRIPT NEVER USE IN PRODUCTION");
  cleanup();
  print("white", "fetching accounts");
  const accounts = await hreToSifchainAccountsAsync();
  print("success", `Accounts Fetched: ${JSON.stringify(accounts)}`);
  const cosmosBridgeFactory = await ethers.getContractFactory("CosmosBridge");
  const validatorPowers = accounts.validatorAccounts.map(() => 100);
  const validatorAccounts = accounts.validatorAccounts.map(acc => acc.address);
  const threshold = validatorPowers.reduce((acc, x) => acc + x);
  print("white", "Deploying Cosmos Bridge contract");
  const cosmosBridge = (await upgrades.deployProxy(cosmosBridgeFactory, [
    accounts.operatorAccount.address, // _operator
    threshold, // _consensusThreshold
    validatorAccounts, // _initValidators
    validatorPowers, // _initPowers
    NETWORK_DESCRIPTOR
  ]) as CosmosBridge);
  print("success",`cosmosBridge deployed at address ${cosmosBridge.address}`);

  print("white", "deploying blocklist contract");
  const blocklistFactory = await ethers.getContractFactory("Blocklist");
  const blocklist = await blocklistFactory.connect(accounts.operatorAccount).deploy();
  print("success", `blocklist deployed successfully at address: ${blocklist.address}`);
  
  print("white", "Setting up ERowan ERC20 bridge token contract");
  const rowanFactory = await ethers.getContractFactory("BridgeToken");
  const rowan = await rowanFactory.deploy(
    "erowan",
    "erowan",
    18,
    "rowan"
  );
  print("success", `ERowan BridgeToken setup at address ${rowan.address}`);

  print("white", "Deploying and setting up bridgebank contract");
  const bridgeBankFactory = await ethers.getContractFactory("BridgeBank");
  const bridgeBank = (await upgrades.deployProxy(bridgeBankFactory, [
    accounts.operatorAccount.address, // _operator
    cosmosBridge.address, // _cosmosBridgeAddress
    accounts.ownerAccount.address, // _owner
    accounts.pauserAccount.address, // _pauser
    NETWORK_DESCRIPTOR,
    rowan.address
  ], {
    // Required because openZepplin Address library has a function that uses delegatecall 
    // delegate call is never used by our code and this library function is unused
    unsafeAllow: ["delegatecall"],
    initializer: "initialize(address,address,address,address,int32,address)"

  })) as BridgeBank;
  print("success", `Bridgebank deployed at address: ${bridgeBank.address}, must now finish setting up`);

  // Bridgebank must immediately call reinitialize
  await bridgeBank.connect(accounts.operatorAccount).reinitialize(
    accounts.operatorAccount.address, 
    cosmosBridge.address, 
    accounts.ownerAccount.address,
    accounts.pauserAccount.address,
    NETWORK_DESCRIPTOR,
    rowan.address
  );

  await bridgeBank.connect(accounts.operatorAccount).setBlocklist(blocklist.address);
  print("success", "Bridgebank setup successfully");

  print("white", "Setting the bridgebank address on CosmosBridge");
  await cosmosBridge.connect(accounts.operatorAccount).setBridgeBank(
    bridgeBank.address
  );
  print("success", "Successfully set BridgeBank address in Cosmos Bridge");

  print("white", "Setting up bridge registry");
  const bridgeRegistryFactory = await ethers.getContractFactory("BridgeRegistry");
  const bridgeRegistry = await upgrades.deployProxy(bridgeRegistryFactory, [
    cosmosBridge.address,
    bridgeBank.address
  ]);
  print("success",`BridgeRegistry setup at address: ${bridgeRegistry.address}`);

  // We must give bridgebank authority over rowan and revoke are admin rights over it
  print("white", "Attempting to grant BridgeBank Admin and Minting roles to Rowan");
  const rowanAdminRole = await rowan.DEFAULT_ADMIN_ROLE();
  const rowanMinterRole = await rowan.MINTER_ROLE();
  // We do these sequentially so that the nonces increment properly
  await rowan.grantRole(rowanAdminRole, bridgeBank.address),
  await rowan.grantRole(rowanMinterRole, bridgeBank.address)
  print("success", "Bridgebank now has Admin and Minting roles over Rowan");
  print("white", "Attempting to revoke deployer addresses Admin and Minting Roles");
  const rowanDeployer = await rowan.signer.getAddress();
  await rowan.renounceRole(rowanAdminRole, rowanDeployer),
  await rowan.renounceRole(rowanMinterRole, rowanDeployer)
  print("success", "Admin and Minter roles have been revoked from deployer");

  print("white", "Add Rowan to the CosmosWhiteList on BridgeBank");
  await bridgeBank.connect(accounts.ownerAccount).addExistingBridgeToken(rowan.address);
  print("success", "Rowan successfully added to CosmosWhiteList on BridgeBank");


  return {
    blocklist: blocklist.address,
    cosmosBridge: cosmosBridge.address,
    bridgeBank: bridgeBank.address,
    bridgeRegistry: bridgeRegistry.address,
    rowanContract: rowan.address
  }
}

print("magenta", "Attempting to deploy Development contracts as requested");
main()
  .then((result) => {
    print("bigSuccess", "All contracts deployed successfully, standby for JSON of addresses");
    console.log("\n\n\n");
    console.log(JSON.stringify(result))
    process.exit(0);
  })
  .catch((error) => {
    print("error", `Something has gone wrong with contract deployment: ${error}`)
    process.exit(1)
  })
