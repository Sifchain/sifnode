import { KEPLR_CONFIG } from "../config";
import { getExtensionPage } from "../utils";

export class KeplrNotificationPopup {
  constructor(config = KEPLR_CONFIG) {
    this.config = config;
    this.url = "/popup.html#/sign?interaction=true&interactionInternal=false";
  }

  async navigate() {
    const targetUrl = `chrome-extension://${this.config.id}${this.url}`;

    this.page = await getExtensionPage(this.config.id, this.url);
    if (!this.page) {
      this.page = await getExtensionPage(this.config.id);
      if (!this.page) {
        this.page = await context.newPage();
      }
    }
    if ((await this.page.url()) !== targetUrl) await this.page.goto(targetUrl);
  }

  async clickApprove() {
    this.page.click('button:has-text("Approve")');
  }
}

export const keplrNotificationPopup = new KeplrNotificationPopup();
