/**
 * This script will actually upgrade BridgeBank in production and connect it to the Blocklist.
 * It will also work with Ropsten, as long as you set
 */
require("dotenv").config();

const hardhat = require("hardhat");
const fs = require("fs-extra");
const Web3 = require("web3");
const web3 = new Web3();

const support = require("../../helpers/forkingSupport");
const { print } = require("../../helpers/utils");
const toInject_1 = require("./injector_upgrade_blocklist-1.json");
const toInject_2 = require("./injector_upgrade_blocklist-2.json");

// Defaults to the Ethereum Mainnet address
const BLOCKLIST_ADDRESS =
  process.env.BLOCKLIST_ADDRESS || "0x1FBeF5a068bFCC4CB1Fae9039EA716EAaaDaeA82";

// If there is no DEPLOYMENT_NAME env var, we'll use the mainnet deployment
const DEPLOYMENT_NAME = process.env.DEPLOYMENT_NAME || "sifchain-1";

// If there is no FORKING_CHAIN_ID env var, we'll use the mainnet id
const CHAIN_ID = process.env.CHAIN_ID || 1;

// Are we running this script in test mode?
const USE_FORKING = !!process.env.USE_FORKING;

const state = {
  addresses: {
    pauser: "0x627306090abaB3A6e1400e9345bC60c78a8BEf57", // TODO: get the correct pauser address (this is not it)
    whitelistedToken: "0x5a98fcbea516cf06857215779fd812ca3bef1b32",
  },
  signers: {
    admin: null,
    operator: null,
  },
  contracts: {
    bridgeBank: null,
    cosmosBridge: null,
    blocklist: null,
    upgradedBridgeBank: null,
  },
  storageSlots: {
    before: {
      pauser: "",
      owner: "",
      nonce: "",
      whitelist: "",
      lockedFunds: "",
    },
    after: {
      pauser: "",
      owner: "",
      nonce: "",
      whitelist: "",
      lockedFunds: "",
    },
  },
};

async function main() {
  print("highlight", "~~~ UPGRADE BRIDGEBANK: BLOCKLIST ~~~");

  // Fetch the manifest and inject the new variables
  copyManifest(true);

  // Connect to each contract
  await connectToContracts();

  // Get signers to send transactions
  await setupAccounts();

  // Fetch current values from the deployed contract
  await setStorageSlots();

  // Upgrade BridgeBank
  await upgradeBridgeBank();

  // Fetch values after the upgrade
  await setStorageSlots(false);

  // Compare slots before and after the upgrade
  checkStorageSlots();

  // Register the Blocklist in BridgeBank
  await registerBlocklist();

  // Populate the Blocklist // should this be run before the upgrade? I think it should, and it be done separetely
  //await syncBlocklist();

  // Save the deployment
  //await saveDeployment();

  // Clean up temporary files
  cleanup();

  print("highlight", "~~~ DONE! 👏 Everything worked as expected. ~~~");
}

async function setStorageSlots(beforeUpgrade = true) {
  const contract = beforeUpgrade ? state.contracts.bridgeBank : state.contracts.upgradedBridgeBank;
  const prefix = beforeUpgrade ? "before" : "after";

  state.storageSlots[prefix].pauser = await contract.pausers(state.addresses.pauser);
  state.storageSlots[prefix].owner = await contract.owner();
  state.storageSlots[prefix].nonce = await contract.lockBurnNonce();
  state.storageSlots[prefix].whitelist = await contract.getTokenInEthWhiteList(
    state.addresses.whitelistedToken
  );
  state.storageSlots[prefix].lockedFunds = await contract.getLockedFunds(
    state.addresses.whitelistedToken
  );
}

function checkStorageSlots() {
  print("yellow", "🎯 Checking storage layout");

  const storage = state.storageSlots;
  testMatch(storage.before.pauser, storage.after.pauser, "Pauser");
  testMatch(storage.before.owner, storage.after.owner, "Owner");
  testMatch(storage.before.nonce, storage.after.nonce, "Nonce");
  testMatch(storage.before.whitelist, storage.after.whitelist, "Whitelist");
  testMatch(storage.before.lockedFunds, storage.after.lockedFunds, "LockedFunds");
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

  // Create an instance of CosmosBridge from the deployed code
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
  } else {
    // If not, we want the connected account
    const accounts = await ethers.getSigners();
    state.signers.admin = accounts[0];
    state.signers.operator = accounts[1];

    if (state.signers.admin.address != support.PROXY_ADMIN_ADDRESS) {
      throw new Error("ACCOUNT_IS_NOT_PROXY_ADMIN");
    }

    if (state.signers.operator.address != operatorAddress) {
      throw new Error("ACCOUNT_IS_NOT_PROXY_ADMIN");
    }
  }

  print("white", `🤵 ProxyAdmin: ${state.signers.admin.address}`);
  print("white", `🤵 Operator: ${state.signers.operator.address}`);
}

async function upgradeBridgeBank() {
  print("yellow", `🕑 Upgrading BridgeBank. Please wait, this may take a while...`);
  const newBridgeBankFactory = await hardhat.ethers.getContractFactory("BridgeBank");
  state.contracts.upgradedBridgeBank = await hardhat.upgrades.upgradeProxy(
    state.contracts.bridgeBank,
    newBridgeBankFactory.connect(state.signers.admin)
  );
  await state.contracts.upgradedBridgeBank.deployed();
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
function injectStorageChanges() {
  print("cyan", "🕵  Injecting changes into manifest");

  // Fetch the deployed manifest
  const currentManifest = fs.readFileSync(
    "./deployments/sifchain-1/.openzeppelin/mainnet.json",
    "utf8"
  );

  // Parse the deployed manifest
  const parsedManifest = JSON.parse(currentManifest);

  // Inject the new variables and change the gap
  toInject_1.parsedManifest = parsedManifest;
  const modManifest_1 = support.injectInManifest(toInject_1);

  toInject_2.parsedManifest = modManifest_1;
  const modManifest_2 = support.injectInManifest(toInject_2);

  // Write to file
  fs.writeFileSync("./.openzeppelin/mainnet.json", JSON.stringify(modManifest_2));
}

async function registerBlocklist() {
  print("yellow", "🕑 Registering the Blocklist in BridgeBank. Please wait...");
  await state.contracts.upgradedBridgeBank
    .connect(state.signers.operator)
    .setBlocklist(BLOCKLIST_ADDRESS);
  print("green", "✅ Blocklist registered in BridgeBank");
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
