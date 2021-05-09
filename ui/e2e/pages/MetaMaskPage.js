import { MM_CONFIG } from "../config";
import { getExtensionPage } from "../utils";
export class MetaMaskPage {
  //   constructor(context, config = MM_CONFIG) {
  constructor(config = MM_CONFIG) {
    this.config = config;
    // this.context = context
  }
  async setup() {
    this.page = await getExtensionPage(this.config.id);
    // this.page = await getExtensionPage(this.context, this.config.id);
    await this.confirmWelcomeScreen();
    await this.importAccount();
    await this.addNetwork();
    await this.page.close();
  }

  async confirmWelcomeScreen() {
    await this.page.click(".welcome-page button");
  }

  async importAccount() {
    await this.page.goto(
      `chrome-extension://${this.config.id}/home.html#initialize/create-password/import-with-seed-phrase`,
    );
    await this.page.type(
      ".first-time-flow__seedphrase input",
      this.config.options.mnemonic,
    );
    await this.page.type("#password", this.config.options.password);
    await this.page.type("#confirm-password", this.config.options.password);
    await this.page.click(".first-time-flow__terms");
    await this.page.click(".first-time-flow button");
    await this.page.click(".end-of-flow button");
    await this.page.click(".popover-header__button");
  }

  async addNetwork() {
    await this.page.goto(
      `chrome-extension://${this.config.id}/home.html#settings/networks`,
    );
    await this.page.click(
      "#app-content > div > div.main-container-wrapper > div > div.settings-page__content > div.settings-page__content__modules > div > div.settings-page__sub-header > div > button",
    );
    await this.page.type("#network-name", this.config.network.name);
    await this.page.type(
      "#rpc-url",
      `http://localhost:${this.config.network.port}`,
    );
    await this.page.type("#chainId", this.config.network.chainId);
    await this.page.click(
      "#app-content > div > div.main-container-wrapper > div > div.settings-page__content > div.settings-page__content__modules > div > div.networks-tab__content > div.networks-tab__network-form > div.network-form__footer > button.button.btn-secondary",
    );
    await this.page.click(
      "#app-content > div > div.main-container-wrapper > div > div.settings-page__header > div.settings-page__close-button",
    );
  }

  async connectAccount(extensionId) {
    await this.page.click("[data-handle='button-connected']");
    const [mmConnectPage] = await Promise.all([
      this.context.waitForEvent("page"),
      this.page.click("button:has-text('Connect Metamask')"),
    ]);
    await mmConnectPage.click(
      "#app-content > div > div.main-container-wrapper > div > div.permissions-connect-choose-account > div.permissions-connect-choose-account__footer-container > div.permissions-connect-choose-account__bottom-buttons > button.button.btn-primary",
    );
    await mmConnectPage.click(
      "#app-content > div > div.main-container-wrapper > div > div.page-container.permission-approval-container > div.permission-approval-container__footers > div.page-container__footer > footer > button.button.btn-primary.page-container__footer-button",
    );
    await this.page.click("text=×");
    return;
  }
}
