import { KEPLR_CONFIG } from "../config";
import { getExtensionPage } from "../utils";
import urls from "../data/urls.json";

export class KeplrNotificationPopup {
  constructor(config = KEPLR_CONFIG) {
    this.config = config;
  }

  async navigate(url = urls.keplr.notificationPopup.generic) {
    // const targetUrl = `chrome-extension://${this.config.id}${url}`;

    this.page = await getExtensionPage(this.config.id, url);
    // for some reason this graceful handling of case where target page url was not present on time
    // and as a result, other existing/new keplr notification page was opened, was causing troubles
    // and additional flakiness. Thus, commenting out to see if above one line solution is stable enough
    // if (!this.page) {
    //   this.page = await getExtensionPage(this.config.id);
    //   if (!this.page) {
    //     this.page = await context.newPage();
    //   }
    // }
    // if ((await this.page.url()) !== targetUrl) await this.page.goto(targetUrl);
  }

  async clickApprove() {
    this.page.click('button:has-text("Approve")');
  }
}

export const keplrNotificationPopup = new KeplrNotificationPopup();
