/* TODO
  x Connect Metamask https://github.com/NodeFactoryIo/dappeteer/blob/master/src/index.ts#L57
  x Different targets, local, sp, etc for deterministic addresses
  x xvfb server for remote test run
  - Peg, unpeg, happy path
  - what's wrong w workflow big.js? (setup at ../ level ? )
  - Cleanup


  TO RUN:
  1. yarn stack
  2. in second terminal: cd e2e && yarn test
*/
require("@babel/polyfill");
const path = require("path");
const fs = require("fs");
const { chromium } = require("playwright");

// configs
const { DEX_TARGET, MM_CONFIG, KEPLR_CONFIG } = require("./config.js");
const keplrConfig = require("../core/src/config.localnet.json");

// extension
const {
  MetaMask,
  connectMmAccount,
  confirmTransaction,
} = require("./metamask.js");
const { importKeplrAccount, connectKeplrAccount } = require("./keplr");

// services
const { getSifchainBalances } = require("./sifchain.js");
const { getEthBalance, advanceEthBlocks } = require("./ethereum.js");

const { extractFile } = require("./utils");

let browserContext;
let dexPage;

function sleep(ms) {
  return new Promise((resolve) => setTimeout(resolve, ms));
}

describe("connect to page", () => {
  beforeAll(async () => {
    // extract extension zips
    await extractExtensionPackages();
    const pathToKeplrExtension = path.join(__dirname, KEPLR_CONFIG.path);
    const pathToMmExtension = path.join(__dirname, MM_CONFIG.path);
    const userDataDir = path.join(__dirname, "./playwright");
    // need to rm userDataDir or else will store extension state
    if (fs.existsSync(userDataDir)) {
      fs.rmdirSync(userDataDir, { recursive: true });
    }

    browserContext = await chromium.launchPersistentContext(userDataDir, {
      // headless required with extensions. xvfb used for ci/cd
      headless: false,
      args: [
        `--disable-extensions-except=${pathToKeplrExtension},${pathToMmExtension}`,
        `--load-extension=${pathToKeplrExtension},${pathToMmExtension}`,
      ],
    });

    // setup metamask
    const MM = new MetaMask(browserContext, MM_CONFIG);
    await MM.setup(browserContext);

    // setup keplr account
    const keplrPage = await browserContext.newPage();
    await keplrPage.goto(
      "chrome-extension://dmkamcknogkgcdfhhbddcghachkejeap/popup.html#/register",
    );
    await importKeplrAccount(keplrPage, KEPLR_CONFIG.options);
    // goto dex page
    dexPage = await browserContext.newPage();
    await dexPage.goto(DEX_TARGET, { waitUntil: "domcontentloaded" });
  });

  afterAll(async () => {
    browserContext.close();
  });

  it("connect to keplr, check balance", async () => {
    const cEthBalance = await getSifchainBalances(
      keplrConfig.sifApiUrl,
      KEPLR_CONFIG.options.address,
      "ceth",
    );

    await connectKeplrAccount(dexPage, browserContext);
    await dexPage.waitForTimeout(1000); // todo capture out extension page close event
    expect(
      (await dexPage.innerText("[data-handle='ceth-row-amount']")).trim(),
    ).toBe(Number(cEthBalance).toFixed(6));
  });

  it("connects to metamask, check balance", async () => {
    const mmEthBalance = await getEthBalance(MM_CONFIG.options.address);
    // connect wallet
    await connectMmAccount(dexPage, browserContext, MM_CONFIG.id);
    await dexPage.waitForTimeout(1000); // todo capture out extension page close event
    // click external tokens tab
    await dexPage.click("text=External Tokens");
    // expect
    expect(await dexPage.innerText("[data-handle='eth-row-amount']")).toBe(
      Number(mmEthBalance).toFixed(6),
    );
  });

  it("pegs", async () => {
    // assumes wallets connected
    const mmEthBalance = await getEthBalance(MM_CONFIG.options.address);
    const cEthBalance = await getSifchainBalances(
      keplrConfig.sifApiUrl,
      KEPLR_CONFIG.options.address,
      "ceth",
    );

    const pegAmount = "1";

    await dexPage.click("[data-handle='external-tab']");
    await dexPage.click("[data-handle='peg-eth']");
    await dexPage.click('input[type="number"]');
    await dexPage.fill('input[type="number"]', pegAmount);
    await dexPage.click('button:has-text("Peg")');
    await dexPage.click('button:has-text("Confirm Peg")');

    await confirmTransaction(dexPage, browserContext, pegAmount, MM_CONFIG.id);
    // move chain forward
    await advanceEthBlocks(50); // NOTE: NOT ASYNC :(
    await sleep(20000);

    const rowAmount = await dexPage.innerText(
      "[data-handle='ceth-row-amount']",
    );

    const expected = (Number(cEthBalance) + Number(pegAmount)).toFixed(6);

    expect(rowAmount.trim()).toBe(expected);
  });

  it("swaps", async () => {});
});

async function extractExtensionPackages() {
  await extractFile(`downloads/${KEPLR_CONFIG.id}.zip`, "./extensions");
  await extractFile(`downloads/${MM_CONFIG.id}.zip`, "./extensions");
  return;
}
