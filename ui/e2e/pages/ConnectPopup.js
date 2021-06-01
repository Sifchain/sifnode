import { assertWaitedText, assertWaitedValue } from "../utils";

export class ConnectPopup {
  constructor() {
    this.elements = {
      connectMetamask: '[data-handle="metamask-connect-button"]',
      connectKeplr: '[data-handle="keplr-connect-button"]',
    };
  }

  async clickConnectMetamask() {
    await page.click(this.elements.connectMetamask);
  }

  async clickConnectKeplr() {
    await page.click(this.elements.connectKeplr);
  }

  async close() {
    await page.click('[data-handle="modal-view-close"]');
  }

  async isKeplrConnected() {
    return await page.$("text='Keplr Connected'");
  }

  async isMetamaskConnected() {
    return await page.$("text='Metamask Connected'");
  }

  async verifyMetamaskConnected() {
    await assertWaitedText(
      this.elements.connectMetamask,
      "Metamask Connected",
      5000,
    );
  }

  async verifyKeplrConnected() {
    await assertWaitedText(this.elements.connectKeplr, "Keplr Connected", 5000);
  }
}

export const connectPopup = new ConnectPopup();
