import { getExtensionPage } from "./utils.js";
export class MetaMask {
  constructor(page, config) {
    this.page = page;
    this.config = config;
  }
  async setup(browserContext) {
    const mmPage = await getExtensionPage(browserContext, this.config.id);
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

export async function connectMmAccount(page, browserContext, extensionId) {
  await page.click("[data-handle='button-connected']");
  await page.click("button:has-text('Connect Metamask')");
  const mmConnectPage = await getExtensionPage(browserContext, extensionId);
  await mmConnectPage.click(
    "#app-content > div > div.main-container-wrapper > div > div.permissions-connect-choose-account > div.permissions-connect-choose-account__footer-container > div.permissions-connect-choose-account__bottom-buttons > button.button.btn-primary",
  );
  await mmConnectPage.click(
    "#app-content > div > div.main-container-wrapper > div > div.page-container.permission-approval-container > div.permission-approval-container__footers > div.page-container__footer > footer > button.button.btn-primary.page-container__footer-button",
  );
  await page.click("text=×");
  return;
}

export async function confirmTransaction(
  page,
  browserContext,
  amount,
  extensionId,
) {
  // extension popup
  const mmConnectPage = await getExtensionPage(browserContext, extensionId);
  await mmConnectPage.click("text=Confirm");
  // haven't yet figured out how to capture close popup event
  await page.waitForTimeout(1000);
  await page.click("text=×");
}

export async function confirmApproval(
  page,
  browserContext,
  amount,
  extensionId,
) {
  const mmConnectPage = await getExtensionPage(browserContext, extensionId);
  await mmConnectPage.click("text=View full transaction details");
  await expect(mmConnectPage).toHaveText(amount);
  await mmConnectPage.click("text=Confirm");
  await page.waitForTimeout(1000);
}

export async function resetAccount(browserContext, extensionId) {
  const page = await browserContext.newPage();

  await page.goto(
    `chrome-extension://${extensionId}/home.html#settings/advanced`,
    {
      waitUntil: "domcontentloaded",
    },
  );
  await page.waitForTimeout(1000);
  await page.click('[data-testid="advanced-setting-reset-account"] button');
  await page.waitForTimeout(1000);
  await page.click('.modal-container button:has-text("Reset")');
}
