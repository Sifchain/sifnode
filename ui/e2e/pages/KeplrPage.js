import { KEPLR_CONFIG } from "../config";

export class KeplrPage {
  constructor(config = KEPLR_CONFIG) {
    this.config = config;
  }

  async setup() {
    this.page = await getExtensionPage(this.config.id);
    await this.importAccount();
  }

  async importAccount() {
    await this.page.goto(
      `chrome-extension://${this.config.id}/popup.html#/register`,
    );
  }
}
