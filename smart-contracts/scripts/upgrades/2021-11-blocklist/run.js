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

// Helps converting an address to a checksum address
const addr = Web3.utils.toChecksumAddress;

// Defaults to the Ethereum Mainnet address
const BLOCKLIST_ADDRESS =
  process.env.BLOCKLIST_ADDRESS || addr("0x9C8a2011cCb697D7EDe3c94f9FBa5686a04DeACB");

// If there is no DEPLOYMENT_NAME env var, we'll use the mainnet deployment
const DEPLOYMENT_NAME = process.env.DEPLOYMENT_NAME || "sifchain-1";

// If there is no FORKING_CHAIN_ID env var, we'll use the mainnet id
const CHAIN_ID = process.env.CHAIN_ID || 1;

// Are we running this script in test mode?
const USE_FORKING = !!process.env.USE_FORKING;

// Will estimate gas and multiply the result by this value to use as MaxFeePerGas
const GAS_PRICE_BUFFER = 5;

const state = {
  addresses: {
    pauser1: addr("c0a586fb260b2c14098a9d95b75f56f13cad2dd9"),
    pauser2: addr("0x9910ade709043d8b9ed2a31fdfcbfb6538f9a397"),
    whitelistedToken: addr("0x5a98fcbea516cf06857215779fd812ca3bef1b32"),
    cosmosWhitelistedtoken: addr("0x8D983cb9388EaC77af0474fA441C4815500Cb7BB"),
    knownBlocklistedAddress: addr("0x8576acc5c05d6ce88f4e49bf65bdf0c62f91353c"),
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
  },
  storageSlots: {
    before: {
      pauser1: "",
      pauser2: "",
      owner: "",
      nonce: "",
      whitelist: "",
      cosmosWhitelist: "",
      lockedFunds: "",
    },
    after: {},
  },
  maxFeePerGas: 1e12,
};

async function main() {
  print("highlight", "~~~ UPGRADE BRIDGEBANK: BLOCKLIST ~~~");

  // Fetch the manifest and inject the new variables
  copyManifest();

  // Estimate gasPrice:
  state.maxFeePerGas = await estimateGasPrice();

  // Connect to each contract
  await connectToContracts();

  // Get signers to send transactions
  await setupAccounts();

  // Pause the system
  await pauseBridgeBank();

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

  // Double-check that the blocklist is registered
  await checkBlocklist();

  // Resume the system
  await resumeBridgeBank();

  // Clean up temporary files
  cleanup();

  print("highlight", "~~~ DONE! 👏 Everything worked as expected. ~~~");

  if (!USE_FORKING) {
    print(
      "h_green",
      `BridgeBank is upgraded IN PRODUCTION and the Blocklist is correctly registered.`
    );
  }
}

async function estimateGasPrice() {
  console.log("Estimating ideal Gas price, please wait...");

  let finalGasPrice;
  try {
    const gasPrice = await ethers.provider.getGasPrice();
    finalGasPrice = Math.round(gasPrice.toNumber() * GAS_PRICE_BUFFER);
  } catch (e) {
    finalGasPrice = state.maxFeePerGas;
  }

  console.log(`Using ideal Gas price: ${ethers.utils.formatUnits(finalGasPrice, "gwei")} GWEI`);

  return finalGasPrice;
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

  await accounts_sanityCheck();
}

async function accounts_sanityCheck() {
  const operatorAddress = await state.contracts.bridgeBank.operator();

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

  const isPaused = await state.contracts.bridgeBank.paused();

  if (!isPaused) {
    await state.contracts.bridgeBank
      .connect(state.signers.pauser)
      .pause({ maxFeePerGas: state.maxFeePerGas });
  }

  print("green", `✅ System is paused`);
}

async function resumeBridgeBank() {
  print("yellow", `🕑 Unpausing the system. Please wait, this may take a while...`);

  const isPaused = await state.contracts.bridgeBank.paused();

  if (isPaused) {
    await state.contracts.bridgeBank
      .connect(state.signers.pauser)
      .unpause({ maxFeePerGas: state.maxFeePerGas });
  }

  print("green", `✅ System has been resumed`);
}

async function setStorageSlots(beforeUpgrade = true) {
  const contract = beforeUpgrade ? state.contracts.bridgeBank : state.contracts.upgradedBridgeBank;
  const prefix = beforeUpgrade ? "before" : "after";

  state.storageSlots[prefix].pauser1 = await contract.pausers(state.addresses.pauser1);
  state.storageSlots[prefix].pauser2 = await contract.pausers(state.addresses.pauser2);
  state.storageSlots[prefix].owner = await contract.owner();
  state.storageSlots[prefix].nonce = await contract.lockBurnNonce();
  state.storageSlots[prefix].whitelist = await contract.getTokenInEthWhiteList(
    state.addresses.whitelistedToken
  );
  state.storageSlots[prefix].cosmosWhitelist = await contract.getCosmosTokenInWhiteList(
    state.addresses.cosmosWhitelistedtoken
  );
  state.storageSlots[prefix].lockedFunds = await contract.getLockedFunds(
    state.addresses.whitelistedToken
  );
}

function checkStorageSlots() {
  print("yellow", "🎯 Checking storage layout");

  const keys = Object.keys(state.storageSlots.before);
  keys.forEach((key) => {
    testMatch(key);
  });
}

function testMatch(key) {
  if (String(state.storageSlots.before[key]) === String(state.storageSlots.after[key])) {
    print("green", `✅ ${key} slot is safe`);
  } else {
    throw new Error(
      `💥 CRITICAL: ${key} Mismatch! From ${state.storageSlots.before[key]} to ${state.storageSlots.after[key]}`
    );
  }
}

async function upgradeBridgeBank() {
  print("yellow", `🕑 Upgrading BridgeBank. Please wait, this may take a while...`);
  const newBridgeBankFactory = await hardhat.ethers.getContractFactory("BridgeBank");
  state.contracts.upgradedBridgeBank = await hardhat.upgrades.upgradeProxy(
    state.contracts.bridgeBank,
    newBridgeBankFactory.connect(state.signers.admin),
    { maxFeePerGas: state.maxFeePerGas }
  );
  await state.contracts.upgradedBridgeBank.deployed();
  print("green", `✅ BridgeBank Upgraded`);
}

// Copy the manifest to the right place (where Hardhat wants it)
function copyManifest(injectChanges) {
  print("cyan", `👀 Fetching the correct manifest`);

  fs.copySync(
    `./deployments/sifchain-1/.openzeppelin/mainnet.json`,
    `./.openzeppelin/mainnet.json`
  );
}

async function registerBlocklist() {
  print("yellow", "🕑 Registering the Blocklist in BridgeBank. Please wait...");
  await state.contracts.upgradedBridgeBank
    .connect(state.signers.operator)
    .setBlocklist(BLOCKLIST_ADDRESS, { maxFeePerGas: state.maxFeePerGas });
  print("green", "✅ Blocklist registered in BridgeBank");
}

async function checkBlocklist() {
  print("yellow", "🕑 Double-checking the blocklist. Please wait...");

  const hasBlocklist = await state.contracts.upgradedBridgeBank.hasBlocklist();
  if (!hasBlocklist) {
    throw new Error(
      "💥 CRITICAL: the blocklist is NOT registered in BridgeBank. Something went wrong."
    );
  }

  const knownBlocklistedAddressIsBlocked = await state.contracts.blocklist.isBlocklisted(
    state.addresses.knownBlocklistedAddress
  );
  if (!knownBlocklistedAddressIsBlocked) {
    console.log(
      `state.addresses.knownBlocklistedAddress: ${state.addresses.knownBlocklistedAddress}`
    );
    console.log(`knownBlocklistedAddressIsBlocked: ${knownBlocklistedAddressIsBlocked}`);

    throw new Error(
      "💥 CRITICAL: cannot find known blocklisted address in the blocklist. Something went wrong."
    );
  }

  print("green", "✅ The Blocklist is accessible and correctly set");
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
