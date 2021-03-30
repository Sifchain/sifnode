import { getExtensionPage } from "./utils.js";
class MetaMask {
  constructor(page, config) {
    this.page = page;
    this.config = config;
  }
  async setup(browserContext) {
    const mmPage = await getExtensionPage(browserContext, "home.html");
    await confirmWelcomeScreen(mmPage);
    await importAccount(mmPage, this.config);
    await addNetwork(mmPage, this.config);
  }
}

async function confirmWelcomeScreen(mmPage) {
  await mmPage.click(".welcome-page button");
}

async function importAccount(mmPage, config) {
  await mmPage.goto(
    `chrome-extension://${config.id}/home.html#initialize/create-password/import-with-seed-phrase`,
  );
  await mmPage.type(
    ".first-time-flow__seedphrase input",
    config.options.mnemonic,
  );
  await mmPage.type("#password", config.options.password);
  await mmPage.type("#confirm-password", config.options.password);
  await mmPage.click(".first-time-flow__terms");
  await mmPage.click(".first-time-flow button");
  await mmPage.click(".end-of-flow button");
  await mmPage.click(".popover-header__button");
}

async function addNetwork(mmPage, config) {
  await mmPage.goto(
    `chrome-extension://${config.id}/home.html#settings/networks`,
  );
  await mmPage.click(
    "#app-content > div > div.main-container-wrapper > div > div.settings-page__content > div.settings-page__content__modules > div > div.settings-page__sub-header > div > button",
  );
  await mmPage.type("#network-name", config.network.name);
  await mmPage.type("#rpc-url", `http://localhost:${config.network.port}`);
  await mmPage.type("#chainId", config.network.chainId);
  await mmPage.click(
    "#app-content > div > div.main-container-wrapper > div > div.settings-page__content > div.settings-page__content__modules > div > div.networks-tab__content > div.networks-tab__network-form > div.network-form__footer > button.button.btn-secondary",
  );
  await mmPage.click(
    "#app-content > div > div.main-container-wrapper > div > div.settings-page__header > div.settings-page__close-button",
  );
}

async function connectMmAccount(page, browserContext) {
  await page.click("[data-handle='button-connected']");
  await page.click("button:has-text('Connect Metamask')");
  const mmConnectPage = await getExtensionPage(
    browserContext,
    "notification.html",
  );
  await mmConnectPage.click(
    "#app-content > div > div.main-container-wrapper > div > div.permissions-connect-choose-account > div.permissions-connect-choose-account__footer-container > div.permissions-connect-choose-account__bottom-buttons > button.button.btn-primary",
  );
  await mmConnectPage.click(
    "#app-content > div > div.main-container-wrapper > div > div.page-container.permission-approval-container > div.permission-approval-container__footers > div.page-container__footer > footer > button.button.btn-primary.page-container__footer-button",
  );
  await page.click("text=×");
  return;
}

module.exports = { MetaMask, connectMmAccount };
