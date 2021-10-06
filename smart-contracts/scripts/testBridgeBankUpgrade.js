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
  const operator_bb = await bridgeBank.operator();
  print("cyan", `Operator: ${operator_bb}`);

  // Impersonate the admin account
  const admin = await support.impersonateAccount(
    support.PROXY_ADMIN_ADDRESS,
    "10000000000000000000"
  );

  // Upgrade BridgeBank
  const newBridgeBankFactory = await hardhat.ethers.getContractFactory(
    "BridgeBank"
  );
  await hardhat.upgrades.upgradeProxy(
    bridgeBank,
    newBridgeBankFactory.connect(admin)
  );

  // Clean up temporary files
  cleanup();

  print("highlight", "~~~ DONE! ~~~");
}

// Copy the manifest to the right place (where Hardhat wants it)
function copyManifest(injectChanges) {
  print("cyan", `Fetching the correct manifest`);

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

// Delete temporary files (the copied manifest)
function cleanup() {
  print("cyan", `Cleaning up temporary files`);

  fs.unlinkSync(`./.openzeppelin/mainnet.json`);
}

main()
  .catch((error) => {
    print("h_red", error.stack);
  })
  .finally(() => process.exit(0));
