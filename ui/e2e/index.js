const { chromium } = require("playwright");
const path = require("path");

const KEPLR_PATH = "./extensions/dmkamcknogkgcdfhhbddcghachkejeap/0.8.1_0";
const METAMASK_PATH = "./extensions/nkbihfbeogaeaoehlefnkodbefgpgknn/9.1.1_0";
const pathToMetamaskExtension = path.join(__dirname, METAMASK_PATH);
const pathToKeplrExtension = path.join(__dirname, KEPLR_PATH);

const userDataDir = path.join(__dirname, "./playwright");

(async () => {
  const browserContext = await chromium.launchPersistentContext(userDataDir, {
    headless: false,
    args: [
      `--disable-extensions-except=${pathToKeplrExtension},${pathToMetamaskExtension}`,
      `--load-extension=${pathToKeplrExtension},${pathToMetamaskExtension}`,
    ],
  });

  const keplrPage = await browserContext.newPage();
  await keplrPage.goto(
    "chrome-extension://dmkamcknogkgcdfhhbddcghachkejeap/popup.html#/register"
  );

  // Metamask opens new tab on its own (hence "close home page")
  const metamaskPage = await browserContext.newPage();
  await metamaskPage.goto(
    "chrome-extension://nkbihfbeogaeaoehlefnkodbefgpgknn/home.html"
  );
})();
