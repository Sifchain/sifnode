/* TODO
  - Connect Metamask https://github.com/NodeFactoryIo/dappeteer/blob/master/src/index.ts#L57
  x Different targets, local, sp, etc for deterministic addresses
  x xvfb server for remote test run
  - Peg, unpeg, happy path
  - Cleanup


  TO RUN:
  1. yarn stack
  2. in second terminal: cd e2e && yarn test
*/
import "@babel/polyfill";
const path = require("path");
const fs = require("fs");
const { chromium } = require("playwright");

const { importKeplrAccount, connectKeplrAccount } = require("./keplr");
const keplrConfig = require("../core/src/config.localnet.json");

const { getSifchainBalances } = require("./sifchain.js");
const { getEthBalance } = require("./ethereum.js")

const { extractFile } = require("./utils");
const { MetaMask, connectMmAccount } = require("./metamask.js")

const DEX_TARGET = "localhost:8080";

const KEPLR_CONFIG = {
  id: "dmkamcknogkgcdfhhbddcghachkejeap",
  ver: "0.8.1_0",
  get path() {
    return `./extensions/${this.id}/${this.ver}`
  },
  options: {
    address: "sif1m625hcmnkc84cgmef6upzzyfu6mxd4jkpnfwwl",
    name: "juniper",
    mnemonic:
      "clump genre baby drum canvas uncover firm liberty verb moment access draft erupt fog alter gadget elder elephant divide biology choice sentence oppose avoid",
  }
}

const MM_CONFIG = {
  id: "nkbihfbeogaeaoehlefnkodbefgpgknn",
  ver: "9.1.1_0",
  get path() {
    return `./extensions/${this.id}/${this.ver}`
  },
  network: {
    name: "mm-e2e",
    port: "7545",
    chainId: "1337"
  },
  options: {
    address: "0x627306090abaB3A6e1400e9345bC60c78a8BEf57",
    mnemonic: "candy maple cake sugar pudding cream honey rich smooth crumble sweet treat",
    password: "coolguy21"
  }
}

let browserContext;
let dexPage;

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
    const MM = new MetaMask(browserContext, MM_CONFIG)
    await MM.setup(browserContext)

    // setup keplr account
    // const keplrPage = await browserContext.newPage();
    // await keplrPage.goto(
    //   "chrome-extension://dmkamcknogkgcdfhhbddcghachkejeap/popup.html#/register",
    // );
    // await importKeplrAccount(keplrPage, KEPLR_CONFIG.options);
    

    dexPage = await browserContext.newPage();
    await dexPage.goto(DEX_TARGET, { waitUntil: "domcontentloaded" });
  });

  afterAll(async () => {
    browserContext.close();
  });

  it.skip("connect to keplr, check balance", async () => {
    const cEthBalance = await getSifchainBalances(
      keplrConfig.sifApiUrl,
      KEPLR_CONFIG.options.address,
      "ceth",
    ); //"100.000000"; // Fetch balance

    await connectKeplrAccount(dexPage, browserContext);
    await dexPage.waitForTimeout(1000); // this is only necessary bc popup
    await dexPage.pause();
    expect(
      await dexPage.innerText("[data-handle='ceth-row-amount']",
      ),
    ).toBe(cEthBalance);
  });

  it("connects to metamask, check balance", async () => {
    const mmEthBalance = await getEthBalance(MM_CONFIG.options.address)
    await dexPage.waitForTimeout(1000); // this is only necessary bc popup
    // connect wallet
    await connectMmAccount(dexPage, browserContext)
    // see https://github.com/NodeFactoryIo/dappeteer/blob/master/src/index.ts#L57

  })
  it("pegs", async () => {})

});

async function extractExtensionPackages() {
  await extractFile(`downloads/${KEPLR_CONFIG.id}.zip`, "./extensions");
  await extractFile(`downloads/${MM_CONFIG.id}.zip`, "./extensions");
  return
} 