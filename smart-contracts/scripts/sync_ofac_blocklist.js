require("dotenv").config();

const support = require("./helpers/forkingSupport");
const { print } = require("./helpers/utils");
const parser = require("./helpers/ofacParser");

// Defaults to the Ethereum Mainnet address
const BLOCKLIST_ADDRESS =
  process.env.BLOCKLIST_ADDRESS || "0x2a9714AA672419cd189Dd8ac06Fc61478573Ba2C";

const USE_FORKING = !!process.env.USE_FORKING;

const state = {
  ofac: [],
  evm: [],
  toAdd: [],
  toRemove: [],
  blocklistInstance: null,
};

async function main() {
  print("highlight", "~~~ SYNC OFAC BLOCKLIST ~~~");

  // Fetches lists, compares them and figures out what has to be added or removed
  await setupState();

  // If we're forking, we want to impersonate the owner account
  if (USE_FORKING) await setupForking();

  // Add addresses to the blocklist
  await addToBlocklist();
  print("cyan", `----`);

  // Remove addresses from the blocklist
  await removeFromBlocklist();
  print("cyan", `----`);

  // Print success
  print("h_green", "Our EVM blocklist is synced with OFAC's blocklist");
  print("highlight", "~~~ DONE ~~~");
}

async function setupState() {
  // Set the deployed blocklist instance
  state.blocklistInstance = await support.getContractAt(
    "Blocklist",
    BLOCKLIST_ADDRESS
  );

  // Set the OFAC list
  state.ofac = await parser.getList();
  print("cyan", `OFAC LIST: ${state.ofac}`);
  print("cyan", `----`);

  // Set the EVM list
  print("yellow", "Fetching EVM blocklist...");
  state.evm = await state.blocklistInstance.getFullList();
  print("cyan", `EVM LIST : ${state.evm}`);
  print("cyan", `----`);

  // Find out what the diff betweeen lists is
  print("yellow", "Calculating Diff...");

  // Addresses that must be added don't exist on evm, but exist on ofac
  state.toAdd = state.ofac.filter((address) => !state.evm.includes(address));
  print("cyan", `Will add: ${state.toAdd}`);

  // Addresses that must be removed exist on evm, but don't exist on ofac
  state.toRemove = state.evm.filter((address) => !state.ofac.includes(address));
  print("cyan", `Will remove: ${state.toRemove}`);
  print("cyan", "----");
}

async function setupForking() {
  print("magenta", "MAINNET FORKING :: IMPERSONATE ACCOUNT");
  // Fetch the current owner of the blocklist
  const ownerAddress = await state.blocklistInstance.owner();

  // Impersonate the blocklist owner
  const owner = await support.impersonateAccount(
    ownerAddress,
    "10000000000000000000"
  );

  // Set the owner as the caller for blocklist functions
  state.blocklistInstance = state.blocklistInstance.connect(owner);
  print("cyan", "----");
}

async function addToBlocklist() {
  if (state.toAdd.length === 0) {
    print("yellow", "The are no new addresses to add to the blocklist");
    return;
  }

  print("yellow", "Adding addresses to the blocklist. Please wait...");

  let tx;
  if (state.toAdd.length === 1) {
    tx = await state.blocklistInstance
      .addToBlocklist(state.toAdd[0])
      .catch((e) => {
        throw e;
      });
  } else {
    // there are many addresses to add
    tx = await state.blocklistInstance
      .batchAddToBlocklist(state.toAdd)
      .catch((e) => {
        throw e;
      });
  }

  print("cyan", `Added ${state.toAdd} to the blocklist.`);
  print("h_green", `TX Hash: ${tx.hash}`);
}

async function removeFromBlocklist() {
  if (state.toRemove.length === 0) {
    print("yellow", "The are no addresses to remove from the blocklist");
    return;
  }

  print("yellow", "Removing addresses from the blocklist. Please wait...");

  let tx;
  if (state.toRemove.length === 1) {
    tx = await state.blocklistInstance
      .removeFromBlocklist(state.toRemove[0])
      .catch((e) => {
        throw e;
      });
  } else {
    // there are many addresses to remove
    tx = await state.blocklistInstance
      .batchRemoveFromBlocklist(state.toRemove)
      .catch((e) => {
        throw e;
      });
  }

  print("cyan", `Removed ${state.toRemove} from the blocklist.`);
  print("h_green", `TX Hash: ${tx.hash}`);
}

function treatCommonErrors(e) {
  if (e.message.indexOf("getFullList") !== -1) {
    print(
      "h_red",
      "Error: cannot execute functions on the blocklist contract. Are you sure you have the right address in your .env variables?"
    );
  } else if (e.message.indexOf("Unsupported method") !== -1) {
    print(
      "h_red",
      "Error: if you are NOT trying to test this with a mainnet fork, please remove the variable USE_FORKING from your .env"
    );
  } else if (e.message.indexOf("insufficient funds") !== -1) {
    print(
      "h_red",
      "Error: insufficient funds. If you are using the correct private key, please refill your account with EVM native coins."
    );
  } else if (e.message.indexOf("caller is not the owner") !== -1) {
    print(
      "h_red",
      "Error: caller is not the owner. Either you have the wrong private key set in your .env, or you should add USE_FORKING=1 to your .env if you want to test the script."
    );
  } else {
    console.error({ e });
  }
}

main()
  .catch((error) => {
    treatCommonErrors(error);
  })
  .finally(() => process.exit(0));
