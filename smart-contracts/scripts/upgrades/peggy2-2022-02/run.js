/**
 * This script will upgrade BridgeBank and CosmosBridge in production and connect it to the Blocklist.
 * For instructions, please consult scripts/upgrades/peggy2-2022-02/runbook.md
 */
require("dotenv").config();

const hardhat = require("hardhat");
const fs = require("fs-extra");
const Web3 = require("web3");

const support = require("../../helpers/forkingSupport");
const { print } = require("../../helpers/utils");

// Helps converting an address to a checksum address
const addr = Web3.utils.toChecksumAddress;

// Defaults to the Ethereum Mainnet address
const BLOCKLIST_ADDRESS =
  process.env.BLOCKLIST_ADDRESS || addr("0x1FBeF5a068bFCC4CB1Fae9039EA716EAaaDaeA82");

// If there is no DEPLOYMENT_NAME env var, we'll use the mainnet deployment
const DEPLOYMENT_NAME = process.env.DEPLOYMENT_NAME || "sifchain-1";

// If there is no FORKING_CHAIN_ID env var, we'll use the mainnet id
const CHAIN_ID = process.env.CHAIN_ID || 1;

// Are we running this script in test mode?
const USE_FORKING = !!process.env.USE_FORKING;

const state = {
  addresses: {
    pauser1: addr("c0a586fb260b2c14098a9d95b75f56f13cad2dd9"),
    pauser2: addr("0x9910ade709043d8b9ed2a31fdfcbfb6538f9a397"),
    validator1: addr("0x0D7dEF5C00a8B6ddc58A0255f0a94cc739C6d0B5"),
    validator2: addr("0x9B4002670C210A3b64e13807250BE62B8dEae201"),
    validator3: addr("0xbF45BFc92ebD305d4C0baf8395c4299bdFCE9EA2"),
    validator4: addr("0xeB29C7016eDd2D6B413fceE4C51474FED058005a"),
    cosmosWhitelistedtoken: addr("0x8D983cb9388EaC77af0474fA441C4815500Cb7BB"),
  },
  signers: {
    admin: null,
    operator: null,
    pauser: null,
  },
  contracts: {
    bridgeBank: null,
    cosmosBridge: null,
    blocklist: null,
    upgradedBridgeBank: null,
    upgradedCosmosBridge: null,
  },
  storageSlots: {
    before: {
      pauser1: "",
      pauser2: "",
      owner: "",
      nonce: "",
      cosmosWhitelist: "",
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
      // TODO: once peggy 1 is released with the blocklist, check if the blocklist storage slot is safe too
    },
    after: {},
  },
};

async function main() {
  print("highlight", "~~~ UPGRADE FROM PEGGY 1.0 TO PEGGY 2.0 ~~~");

  // We heeded the warnings and made sure the upgrade will work
  hardhat.upgrades.silenceWarnings();

  // Fetch the manifest and inject the new variables
  copyManifest(true);

  // Connect to each contract
  await connectToContracts();

  // Get signers to send transactions
  await setupAccounts();

  // Fetch current values from the deployed contract
  await setBridgeBankStorageSlots();
  await setCosmosBridgeStorageSlots();

  // Pause the system
  await pauseBridgeBank();

  // Upgrade BridgeBank
  await upgradeBridgeBank();

  // Upgrade CosmosBridge
  await upgradeCosmosBridge();

  // Fetch values after the upgrade
  await setBridgeBankStorageSlots(false);
  await setCosmosBridgeStorageSlots(false);

  // Compare slots before and after the upgrade
  checkStorageSlots();

  // Resume the system
  await resumeBridgeBank();

  // Clean up temporary files
  cleanup();

  print("highlight", "~~~ DONE! 👏 Everything worked as expected. ~~~");

  if (!USE_FORKING) {
    print("h_green", `Peggy 1.0 has been upgraded to Peggy 2.0 in ${currentEnv}`);
  }
}

async function setBridgeBankStorageSlots(beforeUpgrade = true) {
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

async function setCosmosBridgeStorageSlots(beforeUpgrade = true) {
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
  const { instance: bridgeBank } = await support.getDeployedContract(
    DEPLOYMENT_NAME,
    "BridgeBank",
    CHAIN_ID
  );
  state.contracts.bridgeBank = bridgeBank;

  // Create an instance of BridgeBank from the deployed code
  const { instance: cosmosBridge } = await support.getDeployedContract(
    DEPLOYMENT_NAME,
    "CosmosBridge",
    CHAIN_ID
  );
  state.contracts.cosmosBridge = cosmosBridge;

  // Connect to the Blocklist
  state.contracts.blocklist = await support.getContractAt("Blocklist", BLOCKLIST_ADDRESS);

  print("green", `✅ Contracts connected`);
}

async function setupAccounts() {
  const operatorAddress = await state.contracts.bridgeBank.operator();

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

    state.signers.pauser = await support.impersonateAccount(
      state.addresses.pauser1,
      "10000000000000000000",
      "Pauser"
    );
  } else {
    // If not, we want the connected accounts
    const accounts = await ethers.getSigners();
    state.signers.admin = accounts[0];
    state.signers.operator = accounts[1];
    state.signers.pauser = accounts[2];
  }

  const hasCorrectAdmin =
    state.signers.admin.address.toLowerCase() === support.PROXY_ADMIN_ADDRESS.toLowerCase();

  const hasCorrectOperator =
    state.signers.operator.address.toLowerCase() === operatorAddress.toLowerCase();

  const hasCorrectPauser =
    state.signers.pauser.address.toLowerCase() === state.addresses.pauser1.toLowerCase() ||
    state.signers.pauser.address.toLowerCase() === state.addresses.pauser2.toLowerCase();

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
      `The third Private Key is not the BridgeBank PAUSER's private key. Please use the Private Key that corresponds to the address ${state.addresses.pauser1} or ${state.addresses.pauser2}`
    );
  }

  const adminColor = hasCorrectAdmin ? "white" : "red";
  const operatorColor = hasCorrectOperator ? "white" : "red";
  const pauserColor = hasCorrectPauser ? "white" : "red";

  print(adminColor, `🤵 ProxyAdmin: ${state.signers.admin.address}`);
  print(operatorColor, `🤵 Operator: ${state.signers.operator.address}`);
  print(pauserColor, `🤵 Pauser: ${state.signers.pauser.address}`);
}

async function pauseBridgeBank() {
  print(
    "yellow",
    `🕑 Pausing the system before the upgrade. Please wait, this may take a while...`
  );
  await state.contracts.bridgeBank.connect(state.signers.pauser).pause();
  print("green", `✅ System is paused`);
}

async function resumeBridgeBank() {
  print("yellow", `🕑 Unpausing the system. Please wait, this may take a while...`);
  await state.contracts.bridgeBank.connect(state.signers.pauser).unpause();
  print("green", `✅ System has been resumed`);
}

async function upgradeBridgeBank() {
  print("yellow", `🕑 Upgrading BridgeBank. Please wait, this may take a while...`);
  const newBridgeBankFactory = await hardhat.ethers.getContractFactory("BridgeBank");
  state.contracts.upgradedBridgeBank = await hardhat.upgrades.upgradeProxy(
    state.contracts.bridgeBank,
    newBridgeBankFactory.connect(state.signers.admin),
    { unsafeAllow: ["delegatecall"] }
  );
  await state.contracts.upgradedBridgeBank.deployed();
  print("green", `✅ BridgeBank Upgraded`);
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
