require("dotenv").config();

const hardhat = require("hardhat");
const fs = require("fs-extra");

const support = require("./helpers/forkingSupport");
const { print } = require("./helpers/utils");
const toInject = require("../data/injector_upgrade_blocklist.json");

// If there is no DEPLOYMENT_NAME env var, we'll use the mainnet deployment
const DEPLOYMENT_NAME = process.env.DEPLOYMENT_NAME || "sifchain-1";

// If there is no FORKING_CHAIN_ID env var, we'll use the mainnet id
const CHAIN_ID = process.env.FORKING_CHAIN_ID || 1;

// We want to verify whther old pausers are still there after an upgrade
const PAUSER = "0x627306090abaB3A6e1400e9345bC60c78a8BEf57";

async function main() {
  print("highlight", "~~~ TEST BRIDGEBANK UPGRADE ~~~");

  // Make sure we're forking
  support.enforceForking();

  // Fetch the manifest and inject the new variables
  copyManifest(true);

  // Create an instance of BridgeBank from the deployed code
  const { instance: bridgeBank } = await support.getDeployedContract(
    DEPLOYMENT_NAME,
    "BridgeBank",
    CHAIN_ID
  );

  // Fetch and log the operator
  const operator_before = await bridgeBank.operator();
  print("white", `🤵 Operator: ${operator_before}`);

  // Fetch current values from the deployed contract
  const pauser_before = await bridgeBank.pausers(PAUSER);
  const owner_before = await bridgeBank.owner();
  const nonce_before = await bridgeBank.lockBurnNonce();

  // Impersonate the admin account
  const admin = await support.impersonateAccount(
    support.PROXY_ADMIN_ADDRESS,
    "10000000000000000000",
    "Proxy Admin"
  );

  // Upgrade BridgeBank
  const newBridgeBankFactory = await hardhat.ethers.getContractFactory(
    "BridgeBank"
  );
  const upgradedBridgeBank = await hardhat.upgrades.upgradeProxy(
    bridgeBank,
    newBridgeBankFactory.connect(admin)
  );
  await upgradedBridgeBank.deployed();

  // Fetch values after the upgrade
  const pauser_after = await upgradedBridgeBank.pausers(PAUSER);
  const owner_after = await upgradedBridgeBank.owner();
  const nonce_after = await upgradedBridgeBank.lockBurnNonce();

  // Compare values before and after the upgrade
  testMatch(pauser_before, pauser_after, "Pauser");
  testMatch(owner_before, owner_after, "Owner");
  testMatch(nonce_before.toString(), nonce_after.toString(), "LockBurnNonce");

  // Send a prophecy claim to see it fail

  // Set the blocklist and send a prophecy claim to see it go through

  // Block the sender's address and send a prophecy claim to see it fail

  // Clean up temporary files
  cleanup();

  print("highlight", "~~~ DONE! 👏 Everything worked as expected. ~~~");
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

  // Inject the new variable and change the gap
  toInject.parsedManifest = parsedManifest;
  const modManifest = support.injectInManifest(toInject);

  // Write to file
  fs.writeFileSync("./.openzeppelin/mainnet.json", JSON.stringify(modManifest));
}

function testMatch(before, after, slotName) {
  if (before === after) {
    print("green", `✅ ${slotName} slot is safe`);
  } else {
    throw new Error(
      `💥 CRITICAL: ${slotName} Mismatch! From ${before} to ${after}`
    );
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
