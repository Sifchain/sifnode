// global-setup.js
import { globalSetup as playwrightGlobalSetup } from "jest-playwright-preset";
const { chromium } = require("playwright");
const { extractExtensionPackage } = require("./utils");
const { MM_CONFIG, KEPLR_CONFIG } = require("./config.js");
const path = require("path");
const fs = require("fs");

module.exports = async function globalSetup(globalConfig) {
  await playwrightGlobalSetup(globalConfig);

  await extractExtensionPackage(MM_CONFIG.id);
  await extractExtensionPackage(KEPLR_CONFIG.id);
  const pathToKeplrExtension = path.join(__dirname, KEPLR_CONFIG.path);
  const pathToMmExtension = path.join(__dirname, MM_CONFIG.path);
  const userDataDir = path.join(__dirname, "./playwright");
  // need to rm userDataDir or else will store extension state
  if (fs.existsSync(userDataDir)) {
    fs.rmdirSync(userDataDir, { recursive: true });
  }

  const browserServer = await chromium.launchServer();
  const wsEndpoint = browserServer.wsEndpoint();
  const browser = await chromium.connect({ wsEndpoint: wsEndpoint });
  const context = await chromium.launchPersistentContext(userDataDir, {
    // headless required with extensions. xvfb used for ci/cd
    headless: true,
    // headless: false,
    args: [
      `--disable-extensions-except=${pathToKeplrExtension},${pathToMmExtension}`,
      `--load-extension=${pathToKeplrExtension},${pathToMmExtension}`,
    ],
  });
  const page = await browser.newPage();
};
