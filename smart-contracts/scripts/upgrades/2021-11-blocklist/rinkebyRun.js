/**
 * This script will upgrade BridgeBank on ROPSTEN and connect it to the Blocklist.
 * We won't use our usual Ropsten deployment, as we don't want to overwrite values
 * used in our testnet test flows.
 * Instead, we'll use an alternate list of contracts that were deployed specifically for this.
 *
 * The list is as follows (Rinkeby addresses):
 * CosmosBridge: 0x84A096672AA417e1afF2bAA9994247FF73347E55
 * BridgeBank: 0x5CAf4CB0693AD0e8f354A30D01CC20F9496988D4
 * Blocklist: 0xbB4fa6cC28f18Ae005998a5336dbA3bC49e3dc57
 *
 * There are 3 important roles:
 * The Proxy Admin can upgrade BridgeBank: 0xDDEe73fb5c91EDf22fe8293C72BE5Fca7cDbc872
 * The BridgeBank Operator can set the blocklist in BridgeBank: 0xAa13C6edb99Fe18Ca97DE8Cc3c2467a5DabFF998
 * The Pauser can pause and unpause the system: 0x2DAe2e893DB771D01b5CB24d1C26692d9b034D3C
 *
 * Before running this script, please set the following variables on your .env:
 * RINKEBY_URL=https://eth-rinkeby.alchemyapi.io/v2/your-alchemy-id
 * RINKEBY_PRIVATE_KEYS=XXXXXXXX,YYYYYYYY,ZZZZZZZZ
 * ACTIVE_PRIVATE_KEY=RINKEBY_PRIVATE_KEYS
 *
 * Where:
 *
 * ```
 * XXXXXXXX is the PROXY ADMIN private key
 * YYYYYYYY is BridgeBank's and CosmosBridge's OPERATOR private key
 * ZZZZZZZZ is the PAUSER's private key
 * ```
 *
 * They should be separated by a comma, and they have to be in that order (admin first, operator second, pauser third).
 * Please also make sure you changed `your-alchemy-id` for your actual Alchemy id in `RINKEBY_URL`.
 *
 * To run the script, use the following command:
 *
 * ```
 * npx hardhat run scripts/upgrades/2021-11-blocklist/rinkebyRun.js
 * ```
 *
 */
require("dotenv").config();

const hardhat = require("hardhat");
const fs = require("fs-extra");
const Web3 = require("web3");

const support = require("../../helpers/forkingSupport");
const { print } = require("../../helpers/utils");

// Helps converting an address to a checksum address
const addr = Web3.utils.toChecksumAddress;

// Defaults to the Rinkeby address
const BLOCKLIST_ADDRESS =
  process.env.BLOCKLIST_ADDRESS || addr("0xAE864F340043ba8d45dbDdFd589F3957A74Dc7FF");

const BRIDGEBANK_ADDRESS =
  process.env.BRIDGEBANK_ADDRESS || addr("0x5CAf4CB0693AD0e8f354A30D01CC20F9496988D4");

const COSMOSBRIDGE_ADDRESS =
  process.env.COSMOSBRIDGE_ADDRESS || addr("0x84A096672AA417e1afF2bAA9994247FF73347E55");

const PROXY_ADMIN_ADDRESS =
  process.env.PROXY_ADMIN_ADDRESS || addr("0xDDEe73fb5c91EDf22fe8293C72BE5Fca7cDbc872");

// Will estimate gas and multiply the result by this value to use as MaxFeePerGas
const GAS_PRICE_BUFFER = 5;

const state = {
  addresses: {
    pauser: addr("0x2DAe2e893DB771D01b5CB24d1C26692d9b034D3C"),
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
      pauser: "",
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
  const { instance: bridgeBank } = await getDeployedContract("BridgeBank", BRIDGEBANK_ADDRESS);
  state.contracts.bridgeBank = bridgeBank;

  // Create an instance of CosmosBridge from the deployed code
  const { instance: cosmosBridge } = await getDeployedContract(
    "CosmosBridge",
    COSMOSBRIDGE_ADDRESS
  );
  state.contracts.cosmosBridge = cosmosBridge;

  // Connect to the Blocklist
  state.contracts.blocklist = await support.getContractAt("Blocklist", BLOCKLIST_ADDRESS);
  print("green", `✅ Blocklist connected at ${state.contracts.blocklist.address}`);

  print("green", `✅ -- All contracts connected --`);
}

async function getDeployedContract(contractName, forceAddress) {
  chainId = 4;

  const filename = `scripts/upgrades/2021-11-blocklist/deployments-test/${contractName}.json`;
  const artifactContents = fs.readFileSync(filename, { encoding: "utf-8" });
  const parsed = JSON.parse(artifactContents);
  const ethersInterface = new ethers.utils.Interface(parsed.abi);

  const address = forceAddress || parsed.networks[chainId].address;
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

async function setupAccounts() {
  const accounts = await ethers.getSigners();
  state.signers.admin = accounts[0];
  state.signers.operator = accounts[1];
  state.signers.pauser = accounts[2];

  await accounts_sanityCheck();
}

async function accounts_sanityCheck() {
  const operatorAddress = await state.contracts.bridgeBank.operator();

  const hasCorrectAdmin =
    state.signers.admin.address.toLowerCase() === PROXY_ADMIN_ADDRESS.toLowerCase();

  const hasCorrectOperator =
    state.signers.operator.address.toLowerCase() === operatorAddress.toLowerCase();

  const hasCorrectPauser =
    state.signers.pauser.address.toLowerCase() === state.addresses.pauser.toLowerCase();

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
      `The third Private Key is not the BridgeBank PAUSER's private key. Please use the Private Key that corresponds to the address ${state.addresses.pauser}`
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

  state.storageSlots[prefix].pauser = await contract.pausers(state.addresses.pauser);
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
    `./scripts/upgrades/2021-11-blocklist/deployments-test/rinkeby.json`,
    `./.openzeppelin/rinkeby.json`
  );
}

async function registerBlocklist() {
  print("yellow", "🕑 Registering the Blocklist in BridgeBank. Please wait...");

  print(
    "magenta",
    `state.contracts.upgradedBridgeBank.address: ${state.contracts.upgradedBridgeBank.address}`
  );

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
