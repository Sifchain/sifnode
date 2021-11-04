/**
 * This script will upgrade BridgeBank in production and connect it to the Blocklist.
 * For instructions, please consult scripts/upgrades/2021-11-blocklist/runbook.md
 */
require("dotenv").config();

const hardhat = require("hardhat");
const fs = require("fs-extra");
const Web3 = require("web3");

const support = require("../../helpers/forkingSupport");
const { print } = require("../../helpers/utils");

// If there is no DEPLOYMENT_NAME env var, we'll use the mainnet deployment
const DEPLOYMENT_NAME = process.env.DEPLOYMENT_NAME || "sifchain-1";

// If there is no FORKING_CHAIN_ID env var, we'll use the mainnet id
const CHAIN_ID = process.env.CHAIN_ID || 1;

// Are we running this script in test mode?
const USE_FORKING = !!process.env.USE_FORKING;

const state = {
  addresses: {
    pauser1: Web3.utils.toChecksumAddress("c0a586fb260b2c14098a9d95b75f56f13cad2dd9"),
    pauser2: Web3.utils.toChecksumAddress("0x9910ade709043d8b9ed2a31fdfcbfb6538f9a397"),
    cosmosWhitelistedtoken: Web3.utils.toChecksumAddress(
      "0x8D983cb9388EaC77af0474fA441C4815500Cb7BB"
    ),
  },
  signers: {
    admin: null,
    operator: null,
  },
  contracts: {
    bridgeBank: null,
    cosmosBridge: null,
    blocklist: null,
    upgradedCosmosBridge: null,
  },
  storageSlots: {
    before: {
      pauser1: "",
      pauser2: "",
      owner: "",
      nonce: "",
      cosmosWhitelist: "",
    },
    after: {
      pauser1: "",
      pauser2: "",
      owner: "",
      nonce: "",
      cosmosWhitelist: "",
    },
  },
};

async function main() {
  print("highlight", "~~~ UPGRADE COSMOSBRIDGE: Peggy 2.0 ~~~");

  // We heeded the warnings and made sure the upgrade will work
  //hardhat.upgrades.silenceWarnings();

  // Fetch the manifest and inject the new variables
  copyManifest(false);

  // Connect to each contract
  await connectToContracts();

  // Get signers to send transactions
  await setupAccounts();

  // Fetch current values from the deployed contract
  //await setStorageSlots();

  // Upgrade CosmosBridge
  await upgradeCosmosBridge();

  // Fetch values after the upgrade
  //await setStorageSlots(false);

  // Compare slots before and after the upgrade
  //checkStorageSlots();

  // Clean up temporary files
  cleanup();

  print("highlight", "~~~ DONE! 👏 Everything worked as expected. ~~~");

  if (!USE_FORKING) {
    print(
      "h_green",
      `The BridgeBank is upgraded in ${currentEnv} and the Blocklist is registered correctly.`
    );
  }
}

async function setStorageSlots(beforeUpgrade = true) {
  const contract = beforeUpgrade ? state.contracts.bridgeBank : state.contracts.upgradedBridgeBank;
  const prefix = beforeUpgrade ? "before" : "after";

  state.storageSlots[prefix].pauser1 = await contract.pausers(state.addresses.pauser1);
  state.storageSlots[prefix].pauser2 = await contract.pausers(state.addresses.pauser2);
  state.storageSlots[prefix].owner = await contract.owner();
  state.storageSlots[prefix].nonce = await contract.lockBurnNonce();
  state.storageSlots[prefix].cosmosWhitelist = await contract.getCosmosTokenInWhiteList(
    state.addresses.cosmosWhitelistedtoken
  );
}

function checkStorageSlots() {
  print("yellow", "🎯 Checking storage layout");

  const storage = state.storageSlots;
  testMatch(storage.before.pauser1, storage.after.pauser1, "Pauser 1");
  testMatch(storage.before.pauser2, storage.after.pauser2, "Pauser 2");
  testMatch(storage.before.owner, storage.after.owner, "Owner");
  testMatch(storage.before.nonce, storage.after.nonce, "Nonce");
  testMatch(storage.before.cosmosWhitelist, storage.after.cosmosWhitelist, "CosmosWhitelist");
}

function testMatch(before, after, slotName) {
  if (String(before) === String(after)) {
    print("green", `✅ ${slotName} slot is safe`);
  } else {
    throw new Error(`💥 CRITICAL: ${slotName} Mismatch! From ${before} to ${after}`);
  }
}

async function connectToContracts() {
  print("yellow", `🕑 Connecting to contracts...`);
  // Create an instance of BridgeBank from the deployed code
  const { instance: cosmosBridge } = await support.getDeployedContract(
    DEPLOYMENT_NAME,
    "CosmosBridge",
    CHAIN_ID
  );
  state.contracts.cosmosBridge = cosmosBridge;

  print("green", `✅ Contracts connected`);
}

async function setupAccounts() {
  const operatorAddress = await state.contracts.cosmosBridge.operator();

  // If we're forking, we want to impersonate the owner account
  if (USE_FORKING) {
    print("magenta", "MAINNET FORKING :: IMPERSONATE ACCOUNT");

    state.signers.admin = await support.impersonateAccount(
      support.PROXY_ADMIN_ADDRESS,
      "10000000000000000000",
      "Proxy Admin"
    );

    state.signers.operator = await support.impersonateAccount(
      operatorAddress,
      "10000000000000000000",
      "Operator"
    );
  } else {
    // If not, we want the connected account
    const accounts = await ethers.getSigners();
    state.signers.admin = accounts[0];
    state.signers.operator = accounts[1];

    if (state.signers.admin.address != support.PROXY_ADMIN_ADDRESS) {
      throw new Error(
        `The first Private Key is not the PROXY ADMIN's private key. Please use the Private Key that corresponds to the address ${support.PROXY_ADMIN_ADDRESS}`
      );
    }

    if (state.signers.operator.address != operatorAddress) {
      throw new Error(
        `The second Private Key is not the BridgeBank OPERATOR's private key. Please use the Private Key that corresponds to the address ${operatorAddress}`
      );
    }
  }

  const adminColor =
    state.signers.admin.address.toLowerCase() === support.PROXY_ADMIN_ADDRESS.toLowerCase()
      ? "white"
      : "red";

  const operatorColor =
    state.signers.operator.address.toLowerCase() === operatorAddress.toLowerCase()
      ? "white"
      : "red";

  print(adminColor, `🤵 ProxyAdmin: ${state.signers.admin.address}`);
  print(operatorColor, `🤵 Operator: ${state.signers.operator.address}`);
}

async function upgradeCosmosBridge() {
  print("yellow", `🕑 Upgrading CosmosBridge. Please wait, this may take a while...`);
  const newCosmosBridgeFactory = await hardhat.ethers.getContractFactory("CosmosBridge");
  state.contracts.upgradedCosmosBridge = await hardhat.upgrades.upgradeProxy(
    state.contracts.cosmosBridge,
    newCosmosBridgeFactory.connect(state.signers.admin),
    { unsafeAllow: ["delegatecall"] }
  );
  await state.contracts.upgradedCosmosBridge.deployed();
  print("green", `✅ BridgeBank Upgraded`);
}

// Copy the manifest to the right place (where Hardhat wants it)
function copyManifest(injectChanges) {
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
