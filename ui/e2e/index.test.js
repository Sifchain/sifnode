/* TODO
  - Connect Metamask
  - Different targets, local, sp, etc for deterministic addresses
  - xvfb server for remote test run
  - Peg, unpeg, happy path
  - Cleanup
*/

import "@babel/polyfill";
const path = require("path");
const fs = require("fs");
const { chromium } = require("playwright");

const { importKeplrAccount, connectKeplrAccount } = require("./keplr");
const { extractFile } = require("./utils");

const DEX_TARGET =
  "https://gateway.pinata.cloud/ipfs/QmV7XdW4VTiaiyDJuHcqoLj67foxLuoequtTax8f7kgJAf";

const KEPLR_EXT_ID = "dmkamcknogkgcdfhhbddcghachkejeap";
const KEPLR_VERSION = "0.8.1_0";
const KEPLR_PATH = `./extensions/${KEPLR_EXT_ID}/${KEPLR_VERSION}`;
const KEPLR_OPTIONS = {
  name: "juniper",
  mnemonic:
    "clump genre baby drum canvas uncover firm liberty verb moment access draft erupt fog alter gadget elder elephant divide biology choice sentence oppose avoid",
};
// const METAMASK_PATH = "./extensions/nkbihfbeogaeaoehlefnkodbefgpgknn/9.1.1_0";

let browserContext;

describe("connect to page", () => {
  beforeAll(async () => {
    // const pathToMetamaskExtension = path.join(__dirname, METAMASK_PATH);
    await extractFile(`downloads/${KEPLR_EXT_ID}.zip`, "./extensions");
    const pathToKeplrExtension = path.join(__dirname, KEPLR_PATH);
    const userDataDir = path.join(__dirname, "./playwright");
    // need to rm userDataDir or else will store extension state
    if (fs.existsSync(userDataDir)) {
      fs.rmdirSync(userDataDir, { recursive: true });
    }
    browserContext = await chromium.launchPersistentContext(userDataDir, {
      // may be able to run within remote xvfb server
      headless: false,
      args: [
        `--disable-extensions-except=${pathToKeplrExtension}`,
        `--load-extension=${pathToKeplrExtension}`,
      ],
    });
    const keplrPage = await browserContext.newPage();
    await keplrPage.goto(
      "chrome-extension://dmkamcknogkgcdfhhbddcghachkejeap/popup.html#/register"
    );
    await importKeplrAccount(keplrPage, KEPLR_OPTIONS);
  });

  afterAll(async () => {
    browserContext.close();
  });

  it("connect to keplr, check balance", async () => {
    const dexPage = await browserContext.newPage();
    await dexPage.goto(DEX_TARGET);
    await connectKeplrAccount(dexPage, browserContext);
    await dexPage.waitForSelector(".wallet-connect-container", {
      state: "detached",
    });
    expect(
      await dexPage.innerText(
        "#app > div > div.layout > div > div.body > div:nth-child(3) > div:nth-child(2) > div > div:nth-child(3) > div.amount"
      )
    ).toBe("10.00000");
  });
});
