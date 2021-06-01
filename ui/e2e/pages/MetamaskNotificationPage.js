import { MM_CONFIG } from "../config";
import { getExtensionPage } from "../utils";
import expect from "expect-playwright";

export class MetamaskNotificationPopup {
  constructor(config = MM_CONFIG) {
    this.config = config;
    this.url = "/popup.html";
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
    else await this.page.reload();
  }

  async clickViewFullTransactionDetails() {
    await this.page.click(
      ".confirm-approve-content__view-full-tx-button-wrapper :text('View full transaction details')",
    );
  }

  async clickConfirm() {
    await this.page.click('button:has-text("Confirm")');
  }

  async verifyTransactionDetails(expectedText) {
    await expect(this.page).toHaveText(
      ".confirm-approve-content__permission .confirm-approve-content__medium-text",
      expectedText,
    );
  }
}

export const metamaskNotificationPopup = new MetamaskNotificationPopup();
