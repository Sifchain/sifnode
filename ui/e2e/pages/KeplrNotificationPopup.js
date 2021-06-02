import { KEPLR_CONFIG } from "../config";
import { getExtensionPage } from "../utils";
import urls from "../data/urls.json";

export class KeplrNotificationPopup {
  constructor(config = KEPLR_CONFIG) {
    this.config = config;
  }

  async navigate(url = urls.keplr.notificationPopup.generic) {
    this.page = await getExtensionPage(this.config.id, url);
  }

  async clickApprove() {
    this.page.click('button:has-text("Approve")');
  }
}

export const keplrNotificationPopup = new KeplrNotificationPopup();
