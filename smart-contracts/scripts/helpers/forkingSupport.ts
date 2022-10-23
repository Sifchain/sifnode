/**
 * Responsible for fetching deployment data and returning a valid ethers contract instance
 */
import fs from "fs";
import { ethers, network } from "hardhat";
import { print } from "./utils";

// By default, this will work with a mainnet fork,
// but it can also be used to fork Ropsten
const DEPLOYMENT_DIRECTORY = "deployments";
const DEFAULT_DEPLOYMENT_NAME = "sifchain-1";

// The address of the Proxy admin (used to impersonate the account that has permission to upgrade proxies)
const PROXY_ADMIN_ADDRESS = "0x7c6c6ea036e56efad829af5070c8fb59dc163d88";

/**
 * Figures out the correct details for a given contract that has already been deployed in production
 * @param {string} deploymentName
 * @param {string} contractName
 * @param {number} chainId
 * @returns An object containing the factory, the instance, its address and the first user found in the accounts list
 */
export async function getDeployedContract(deploymentName: string, contractName: string, chainId: number) {
  deploymentName = deploymentName ?? DEFAULT_DEPLOYMENT_NAME;
  contractName = contractName ?? "BridgeBank";
  chainId = chainId ?? 1;

  const filename = `${DEPLOYMENT_DIRECTORY}/${deploymentName}/${contractName}.json`;
  const artifactContents = fs.readFileSync(filename, { encoding: "utf-8" });
  const parsed = JSON.parse(artifactContents);
  const ethersInterface = new ethers.utils.Interface(parsed.abi);

  const address = parsed.networks[chainId].address;
  print("yellow", `🕑 Connecting to ${contractName} at ${address} on chain ${chainId}`);

  const accounts = await ethers.getSigners();
  const activeUser = accounts[0];

  const contract = new ethers.Contract(address, ethersInterface, activeUser);
  const instance = await contract.attach(address);

  print("green", `🌎 Connected to ${contractName} at ${address} on chain ${chainId}`);

  return {
    contract,
    instance,
    address,
    activeUser,
  };
}

/**
 * Use this function to impersonate accounts when forking
 * @param {string} address
 * @param {string} newBalance
 * @param {string} accountName A name that will appear in the logs to facilitate things
 * @returns An ethers SIGNER object
 */
export async function impersonateAccount(address: string, newBalance: string, accountName: string) {
  accountName = accountName ? ` (${accountName})` : "";

  print("magenta", `🔒 Impersonating account ${address}${accountName}`);

  await network.provider.request({
    method: "hardhat_impersonateAccount",
    params: [address],
  });

  if (newBalance) {
    await setNewEthBalance(address, newBalance);
  }

  print("magenta", `🔓 Account ${address}${accountName} successfully impersonated`);

  return ethers.getSigner(address);
}

/**
 * When impersonating an account, this function sets its balance
 * @param {string} address
 * @param {string | number} newBalance
 */
export async function setNewEthBalance(address: string, newBalance: string | number) {
  const newValue = `0x${newBalance.toString(16)}`;
  await ethers.provider.send("hardhat_setBalance", [address, newValue]);

  print("magenta", `💰 Balance of account ${address} set to ${newBalance}`);
}

/**
 * Throws an error if USE_FORKING is not set in .env
 */
export function enforceForking() {
  const forkingActive = !!process.env.USE_FORKING;
  if (!forkingActive) {
    throw new Error("❌ Forking is not active. Operation aborted.");
  }
}

/**
 * Returns an instance of the contract on the currently connected network
 * @dev Use this function to connect to a deployed contract
 * @dev THAT HAS THE SAME INTERFACE OF A CONTRACT IN THE CONTRACTS FOLDER
 * @dev It means that it won't work for outdated contracts (for that, please use the function getDeployedContract)
 * @param {string} contractName
 * @param {string} contractAddress
 * @returns An instance of the contract on the currently connected network
 */
export async function getContractAt(contractName: string, contractAddress: string) {
  const factory = await ethers.getContractFactory(contractName);
  return await factory.attach(contractAddress);
}

interface StorageObject {
  contract: string;
  label: string;
  type: string;
  src: string;
}

interface InjectInManifestParam{
  topContractMainnetAddress: string;
  parsedManifest: ContractManifestFile;
  contractName: string;
  previousLabel: string;
  newVarObject: StorageObject;
  previousGapSize: number;
  newGapSize: number;
  newTypeName: string;
  newTypeLabel: string;
}

interface TypesObject {
  [key: string]: {
    label: string
  }
}

interface ContractManifestFile {
  impls: {
    [key: string]: {
      address: string;
      layout: {
        storage: StorageObject[];
        types: TypesObject;
      }
    }
  }
}

/**
  * Injects a new variable in a gapped contract's manifest, so that you can upgrade it without errors
  * @param {string} topContractMainnetAddress Address of the top contract, such as BridgeBank or CosmosBridge (NOT the proxy)
  * @param {object} parsedManifest The manifest after a JSON.parse(manifestFile)
  * @param {string} contractName The name of the modified contract
  * @param {string} previousLabel Your new variable will be injected after this object (you have to manually find that out!)
  * @param {object} newVarObject The object that contains your new variable
  * @param {number} previousGapSize The gap size as it is in the currently deployed contract
  * @param {number} newGapSize The new gap size
  * @param {string} newTypeName The name of your new type, if any (this is optional)
  * @param {string} newTypeLabel The label of your new type, if any  (this is mandatory IF you passed `newTypeName`)
  * @returns {object} The modified manifest object (you can now stringify it and save it to a file)
  * 
  * Example:
   {
     topContractMainnetAddress: '0x714b49640c2a545b672e8bbd53cc8935725c6a14',
     parsedManifest,
     contractName: "EthereumWhiteList",
     previousLabel: "_ethereumTokenWhiteList",
     newVarObject: {
       contract: "EthereumWhiteList",
       label: "blocklist",
       type: "t_contract(IBlocklist)4736",
       src: "../project:/contracts/BridgeBank/EthereumWhitelist.sol:21",
     },
     previousGapSize: 100,
     newGapSize: 99,
     newTypeName: "t_contract(IBlocklist)4736",
     newTypeLabel: "contract IBlocklist",
   }
  */
export function injectInManifest({
  topContractMainnetAddress,
  parsedManifest,
  contractName,
  previousLabel,
  newVarObject,
  previousGapSize,
  newGapSize,
  newTypeName,
  newTypeLabel,
}: InjectInManifestParam) {
  // Make a copy of the manifest
  parsedManifest = { ...parsedManifest };

  // Find the correct implementation in the Manifest
  const impls = parsedManifest.impls;
  const implIndex = Object.keys(impls).findIndex((key) => {
    return impls[key].address.toLowerCase() === topContractMainnetAddress.toLowerCase();
  });
  const impl = impls[Object.keys(impls)[implIndex]];

  // Helpers
  const layout = impl.layout;
  const storage = layout.storage;
  const types = layout.types;
  const newStorage = [];

  // STORAGE
  // Find the slot where to inject the new var
  // @dev: this is not optimal, but we might want it as is to be able to deal with many new vars at once
  const storagePreviousItemIndex = storage.findIndex((elem) => {
    return elem.contract === contractName && elem.label === previousLabel;
  });

  // Populate the new storage up to the slot
  for (let i = 0; i < storagePreviousItemIndex + 1; i++) {
    newStorage.push(storage[i]);
  }

  // Push the new var to storage
  newStorage.push(newVarObject);

  // Finish populating the storage with what was already there
  for (let i = storagePreviousItemIndex + 1; i < storage.length; i++) {
    newStorage.push(storage[i]);
  }

  // GAP IN STORAGE:
  // Find the gap declaration
  const gapIndex = newStorage.findIndex((elem) => {
    return elem.contract === contractName && elem.label === "____gap";
  });

  // Replace the size of the gap
  newStorage[gapIndex]["type"] = newStorage[gapIndex]["type"].replace(String(previousGapSize), String(newGapSize));

  // GAP IN TYPES
  // In the Types object of the manifest, add a new gap with the new size
  types[`t_array(t_uint256)${newGapSize}_storage`] = {
    label: "uint256[${newGapSize}]",
  };

  // Delete the old gap
  delete types[`t_array(t_uint256)${previousGapSize}_storage`];

  // If there's a new type to add, add it to the types object
  if (newTypeName) {
    if (!newTypeLabel) throw new Error("MISSING_NEW_TYPE_LABEL");

    types[newTypeName] = {
      label: newTypeLabel,
    };
  }

  // Restructure the manifest
  layout.storage = newStorage;
  layout.types = types;
  impl.layout = layout;
  parsedManifest.impls[Object.keys(impls)[implIndex]] = impl;

  return parsedManifest;
}

export interface ReplaceTypesInManifestParams {
  parsedManifest: ContractManifestFile;
  originalType: string;
  newType: string;
}

/**
 * Replaces all instances of originalTypeName for newTypeName in a manifest
 * @dev we'll expect a parsedManifest here to maintain the pattern established in injectInManifest()
 * @param {object} parsedManifest
 * @param {string} originalType
 * @param {string} newType
 * @return {object} The updated manifest
 *
 * Example:
 * const modManifest = replaceTypesInManifest({
 *   parsedManifest,
 *   originalType: "t_string_memory",
 *   newType: "t_string_memory_ptr",
 * });
 */
export function replaceTypesInManifest({ parsedManifest, originalType, newType }: ReplaceTypesInManifestParams) {
  // Make a copy of the manifest
  parsedManifest = { ...parsedManifest };

  const stringManifest = JSON.stringify(parsedManifest);
  const replacedStringManifest = stringManifest.replace(new RegExp(originalType, "g"), newType);
  const reconstructedManifest = JSON.parse(replacedStringManifest);

  return reconstructedManifest;
}

export default {
  PROXY_ADMIN_ADDRESS,
  getDeployedContract,
  impersonateAccount,
  setNewEthBalance,
  enforceForking,
  getContractAt,
  injectInManifest,
  replaceTypesInManifest,
};
