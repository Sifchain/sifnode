/**
 * This script will test the Blocklist upgrade off of a Mainnet fork
 * It will try many different scenarios regarding the Blocklist integration.
 * If anything blows up, you'll see a red error on your shell/console.
 *
 * It's a good idea to run this script before actually upgrading the BridgeBank in production.
 *
 * Usage: run the following command
 * `USE_FORKING=1 npx hardhat run scripts/upgrades/2021-11-blocklist/test.js`
 */

require("dotenv").config();

const hardhat = require("hardhat");
const fs = require("fs-extra");
const Web3 = require("web3");
const web3 = new Web3();

const support = require("../../helpers/forkingSupport");
const { print } = require("../../helpers/utils");

// If there is no DEPLOYMENT_NAME env var, we'll use the mainnet deployment
const DEPLOYMENT_NAME = process.env.DEPLOYMENT_NAME || "sifchain-1";

// If there is no FORKING_CHAIN_ID env var, we'll use the mainnet id
const CHAIN_ID = process.env.FORKING_CHAIN_ID || 1;

// Will estimate gas and multiply the result by this value to use as MaxFeePerGas
const GAS_PRICE_BUFFER = 5;

const state = {
  addresses: {
    validator1: "0x0D7dEF5C00a8B6ddc58A0255f0a94cc739C6d0B5",
    validator2: "0x9B4002670C210A3b64e13807250BE62B8dEae201",
    validator3: "0xbF45BFc92ebD305d4C0baf8395c4299bdFCE9EA2",
    validator4: "0xeB29C7016eDd2D6B413fceE4C51474FED058005a",
    operator: "",
    user1: "0xfc854524613dA7244417908d199857754189633c",
    user2: "0xb6fa1F5304aa0a17E5B85088e720b0e39dD1b233",
    user3: "0x6F165B30ee4bFc9565E977Ae252E4110624ab147",
    sifRecipient: web3.utils.utf8ToHex("sif1nx650s8q9w28f2g3t9ztxyg48ugldptuwzpace"),
    pauser: "0xc0a586fb260b2c14098a9d95b75f56f13cad2dd9",
    whitelistedToken: "0x5a98fcbea516cf06857215779fd812ca3bef1b32",
    cosmosWhitelistedtoken: "0x8D983cb9388EaC77af0474fA441C4815500Cb7BB",
  },
  signers: {
    admin: null,
    operator: null,
    user1: null,
    validator1: null,
    validator2: null,
    validator3: null,
    pauser: null,
  },
  contracts: {
    bridgeBank: null,
    cosmosBridge: null,
    blocklist: null,
    bridgeToken: null,
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
    after: {
      pauser: "",
      owner: "",
      nonce: "",
      whitelist: "",
      cosmosWhitelist: "",
      lockedFunds: "",
    },
  },
  tokenBalance: 10000,
  amount: 1000,
  maxFeePerGas: 1e12, // 1000 GWEI as default
};

async function main() {
  print("highlight", "~~~ TEST BRIDGEBANK UPGRADE ~~~");

  // Make sure we're forking
  support.enforceForking();

  // Fetch the manifest
  copyManifest();

  // Estimate gasPrice:
  state.maxFeePerGas = await estimateGasPrice();

  // Deploy or connect to each contract
  await deployContracts();

  // Impersonate accounts
  await impersonateAccounts();

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

  // Resume the system
  await resumeBridgeBank();

  /** Below this point we have tests only; the upgrade is done */
  // Setup the BridgeToken (register in BridgeBank, mint and set allowance)
  await setupBridgeToken();

  // Try to lock tokens to see it go through (because BridgeBank doesn't know the Blocklist yet)
  await lock({ expectedError: null });

  // Set the Blocklist in BridgeBank
  print("yellow", `🕑 Registering the Blocklist in BridgeBank...`);
  await state.contracts.upgradedBridgeBank
    .connect(state.signers.operator)
    .setBlocklist(state.contracts.blocklist.address, { maxFeePerGas: state.maxFeePerGas });
  print("green", `✅ Blocklist registered in BridgeBank`);

  // Try to lock tokens to see it go through
  await lock({ expectedError: null });

  // Block the sender's address
  print("yellow", `🕑 Blocklisting user1...`);
  await state.contracts.blocklist.addToBlocklist(state.addresses.user1);
  print("green", `✅ User1 blocklisted`);

  // Try to lock tokens to see it fail
  await lock({ expectedError: "Address is blocklisted" });

  // UNblock the sender's address
  print("yellow", `🕑 Removing user1 from the blocklist...`);
  await state.contracts.blocklist.removeFromBlocklist(state.addresses.user1);
  print("green", `✅ User1 removed from the blocklist`);

  // Try to lock tokens to see it go through
  await lock({ expectedError: null });

  // Send prohpecyClaim to see it go through 3 times
  await newProphecyClaim({
    signer: state.signers.validator1,
    nonce: 1,
    expectedError: null,
  });
  await newProphecyClaim({
    signer: state.signers.validator2,
    nonce: 1,
    expectedError: null,
  });
  await newProphecyClaim({
    signer: state.signers.validator3,
    nonce: 1,
    expectedError: null,
  });

  // Block user1 again
  print("yellow", `🕑 Blocklisting user1...`);
  await state.contracts.blocklist.addToBlocklist(state.addresses.user1);
  print("green", `✅ User1 blocklisted`);

  // Send prohpecyClaim to see it fail on the third time
  await newProphecyClaim({
    signer: state.signers.validator1,
    nonce: 2,
  });
  await newProphecyClaim({
    signer: state.signers.validator2,
    nonce: 2,
  });
  await newProphecyClaim({
    signer: state.signers.validator3,
    nonce: 2,
    expectedError: "Address is blocklisted",
  });

  // Clean up temporary files
  cleanup();

  print(
    "highlight",
    "~~~ DONE! 👏 Everything worked as expected. It's safe to execute the script in production. ~~~"
  );
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

async function impersonateAccounts() {
  // Fetch and log the operator
  state.addresses.operator = await state.contracts.bridgeBank.operator();
  print("white", `🤵 Operator: ${state.addresses.operator}`);

  state.signers.admin = await support.impersonateAccount(
    support.PROXY_ADMIN_ADDRESS,
    "10000000000000000000",
    "Proxy Admin"
  );

  state.signers.operator = await support.impersonateAccount(
    state.addresses.operator,
    "10000000000000000000",
    "Operator"
  );

  state.signers.pauser = await support.impersonateAccount(
    state.addresses.pauser,
    "10000000000000000000",
    "Pauser"
  );

  state.signers.user1 = await support.impersonateAccount(
    state.addresses.user1,
    "10000000000000000000",
    "User1"
  );

  state.signers.validator1 = await support.impersonateAccount(
    state.addresses.validator1,
    "10000000000000000000",
    "Validator 1"
  );

  state.signers.validator2 = await support.impersonateAccount(
    state.addresses.validator2,
    "10000000000000000000",
    "Validator 2"
  );

  state.signers.validator3 = await support.impersonateAccount(
    state.addresses.validator3,
    "10000000000000000000",
    "Validator 3"
  );
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

async function deployContracts() {
  print("yellow", `🕑 Deploying contracts...`);
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

  // Deploy the Blocklist
  const blocklistFactory = await hardhat.ethers.getContractFactory("Blocklist");
  const blocklist = await blocklistFactory.deploy();
  await blocklist.deployed();
  state.contracts.blocklist = blocklist;

  // Deploy the BridgeToken
  const bridgeTokenFactory = await hardhat.ethers.getContractFactory("BridgeToken");
  const token = await bridgeTokenFactory.deploy("TEST");
  await token.deployed();
  state.contracts.bridgeToken = token;

  print("green", `✅ Contracts deployed`);
}

async function upgradeBridgeBank() {
  print("yellow", `🕑 Upgrading BridgeBank...`);
  const newBridgeBankFactory = await hardhat.ethers.getContractFactory("BridgeBank");
  state.contracts.upgradedBridgeBank = await hardhat.upgrades.upgradeProxy(
    state.contracts.bridgeBank,
    newBridgeBankFactory.connect(state.signers.admin),
    { maxFeePerGas: state.maxFeePerGas }
  );
  await state.contracts.upgradedBridgeBank.deployed();
  print("green", `✅ BridgeBank Upgraded`);
}

async function setupBridgeToken() {
  // Add it to the whitelist (only the OPERATOR can do that)
  print("yellow", `🕑 Adding the token to the whitelist...`);
  await state.contracts.bridgeBank
    .connect(state.signers.operator)
    .updateEthWhiteList(state.contracts.bridgeToken.address, true);
  print("green", `✅ Token added to the whitelist`);

  // Load user account with ERC20 tokens
  print("yellow", `🕑 Minting tokens to user1...`);
  await state.contracts.bridgeToken.mint(state.addresses.user1, state.tokenBalance);
  print("green", `✅ Tokens minted to user1`);

  // Approve tokens to contract
  print("yellow", `🕑 Approving BridgeBank to spend BridgeTokens...`);
  await state.contracts.bridgeToken
    .connect(state.signers.user1)
    .approve(state.contracts.upgradedBridgeBank.address, state.tokenBalance);
  print("green", `✅ BridgeBank approved to spend BridgeTokens`);
}

async function lock({ expectedError }) {
  print("yellow", `🕑 Trying to lock tokens...`);

  let errorMessage;
  try {
    await state.contracts.upgradedBridgeBank
      .connect(state.signers.user1)
      .lock(state.addresses.sifRecipient, state.contracts.bridgeToken.address, state.amount, {
        value: 0,
      });
  } catch (e) {
    errorMessage = e.message;
  }

  treatExpectedError({ functionName: "lock", expectedError, errorMessage });
}

async function newProphecyClaim({ signer, nonce, expectedError }) {
  print("yellow", `🕑 Sending new ProphecyClaim...`);

  let errorMessage;
  try {
    await state.contracts.cosmosBridge.connect(signer).newProphecyClaim(
      2, // LOCK TYPE
      state.addresses.sifRecipient,
      nonce,
      state.addresses.user1,
      "TEST",
      state.amount
    );
  } catch (e) {
    errorMessage = e.message;
  }

  treatExpectedError({
    functionName: "newProphecyClaim",
    expectedError,
    errorMessage,
  });
}

function treatExpectedError({ functionName, expectedError, errorMessage }) {
  if (!expectedError && !errorMessage) {
    print("green", `✅ ${functionName}() went through as expected`);
    return;
  }

  if (expectedError && errorMessage) {
    if (errorMessage.indexOf(expectedError) !== -1) {
      print("green", `✅ ${functionName}() failed as expected`);
    } else {
      throw new Error(
        `💥 CRITICAL: ${functionName}() should have failed with '${expectedError}', but failed with '${errorMessage}'`
      );
    }
    return;
  }

  if (!expectedError && errorMessage) {
    throw new Error(errorMessage);
  }

  if (expectedError && !errorMessage) {
    throw new Error(
      `💥 CRITICAL: ${functionName}() should have failed with '${expectedError}', but it went through normally`
    );
  }

  print(
    "highlight",
    "OOPS: Shouldn't have gotten here! Please review the flow, something is wrong"
  );
}

// Copy the manifest to the right place (where Hardhat wants it)
function copyManifest() {
  print("cyan", `👀 Fetching the correct manifest`);

  fs.copySync(
    `./deployments/sifchain-1/.openzeppelin/mainnet.json`,
    `./.openzeppelin/mainnet.json`
  );
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

  const storage = state.storageSlots;
  testMatch(storage.before.pauser, storage.after.pauser, "Pauser");
  testMatch(storage.before.owner, storage.after.owner, "Owner");
  testMatch(storage.before.nonce, storage.after.nonce, "Nonce");
  testMatch(storage.before.whitelist, storage.after.whitelist, "Whitelist");
  testMatch(storage.before.cosmosWhitelist, storage.after.cosmosWhitelist, "CosmosWhitelist");
  testMatch(storage.before.lockedFunds, storage.after.lockedFunds, "LockedFunds");
}

function testMatch(before, after, slotName) {
  if (String(before) === String(after)) {
    print("green", `✅ ${slotName} slot is safe`);
  } else {
    throw new Error(`💥 CRITICAL: ${slotName} Mismatch! From ${before} to ${after}`);
  }
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
