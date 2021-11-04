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

// Helps converting an address to a checksum address
const addr = Web3.utils.toChecksumAddress;

const state = {
  addresses: {
    validator1: addr("0x0D7dEF5C00a8B6ddc58A0255f0a94cc739C6d0B5"),
    validator2: addr("0x9B4002670C210A3b64e13807250BE62B8dEae201"),
    validator3: addr("0xbF45BFc92ebD305d4C0baf8395c4299bdFCE9EA2"),
    validator4: addr("0xeB29C7016eDd2D6B413fceE4C51474FED058005a"),
    operator: addr("0x2dc894e2e87fb728b2520ce8418983a834357824"),
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
      bridgeBank: "", // addr("0xb5f54ac4466f5ce7e0d8a5cb9fe7b8c0f35b7ba8"),
      consensusThreshold: "", // 75,
      currentValsetVersion: "", // 1,
      hasBridgeBank: "", // true,
      totalPower: "", // 100,
      validatorCount: "", // 4,
      validator1IsValidator: "", // true,
      validator2IsValidator: "", // true,
      validator3IsValidator: "", // true,
      validator4IsValidator: "", // true,
      validator1Power: "", // 25,
      validator2Power: "", // 25,
      validator3Power: "", // 25,
      validator4Power: "", // 25,
      operator: "", // addr("0x2dc894e2e87fb728b2520ce8418983a834357824"),
    },
    after: {
      bridgeBank: "", // addr("0xb5f54ac4466f5ce7e0d8a5cb9fe7b8c0f35b7ba8"),
      consensusThreshold: "", // 75,
      currentValsetVersion: "", // 1,
      hasBridgeBank: "", // true,
      totalPower: "", // 100,
      validatorCount: "", // 4,
      validator1IsValidator: "", // true,
      validator2IsValidator: "", // true,
      validator3IsValidator: "", // true,
      validator4IsValidator: "", // true,
      validator1Power: "", // 25,
      validator2Power: "", // 25,
      validator3Power: "", // 25,
      validator4Power: "", // 25,
      operator: "", // addr("0x2dc894e2e87fb728b2520ce8418983a834357824"),
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
  await setStorageSlots();

  // Upgrade CosmosBridge
  await upgradeCosmosBridge();

  // Fetch values after the upgrade
  await setStorageSlots(false);

  // Compare slots before and after the upgrade
  checkStorageSlots();

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
  const contract = beforeUpgrade
    ? state.contracts.cosmosBridge
    : state.contracts.upgradedCosmosBridge;
  const prefix = beforeUpgrade ? "before" : "after";

  state.storageSlots[prefix].bridgeBank = await contract.bridgeBank();
  state.storageSlots[prefix].consensusThreshold = await contract.consensusThreshold();
  state.storageSlots[prefix].currentValsetVersion = await contract.currentValsetVersion();
  state.storageSlots[prefix].hasBridgeBank = await contract.hasBridgeBank();
  state.storageSlots[prefix].totalPower = await contract.totalPower();
  state.storageSlots[prefix].validatorCount = await contract.validatorCount();

  state.storageSlots[prefix].validator1IsValidator = await contract.isActiveValidator(
    state.addresses.validator1
  );
  state.storageSlots[prefix].validator2IsValidator = await contract.isActiveValidator(
    state.addresses.validator2
  );
  state.storageSlots[prefix].validator3IsValidator = await contract.isActiveValidator(
    state.addresses.validator3
  );
  state.storageSlots[prefix].validator4IsValidator = await contract.isActiveValidator(
    state.addresses.validator4
  );

  state.storageSlots[prefix].validator1Power = await contract.getValidatorPower(
    state.addresses.validator1
  );
  state.storageSlots[prefix].validator2Power = await contract.getValidatorPower(
    state.addresses.validator2
  );
  state.storageSlots[prefix].validator3Power = await contract.getValidatorPower(
    state.addresses.validator3
  );
  state.storageSlots[prefix].validator4Power = await contract.getValidatorPower(
    state.addresses.validator4
  );

  state.storageSlots[prefix].operator = await contract.operator();
}

function checkStorageSlots() {
  print("yellow", "🎯 Checking storage layout");

  const keys = Object.keys(state.storageSlots.before);

  keys.forEach((key) => {
    testMatch(state.storageSlots.before[key], state.storageSlots.after[key], key);
  });
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
        `The second Private Key is not the CosmosBank OPERATOR's private key. Please use the Private Key that corresponds to the address ${operatorAddress}`
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
  print("green", `✅ CosmosBridge Upgraded`);
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
