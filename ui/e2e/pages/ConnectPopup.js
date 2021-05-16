export class ConnectPopup {
  constructor() {
    this.elements = {
      connectMetamask: "button:has-text('Connect Metamask')",
      connectKeplr: "button:has-text('Connect Keplr')",
    };
  }

  async clickConnectMetamask() {
    await page.click(this.elements.connectMetamask);
  }

  async clickConnectKeplr() {
    await page.click(this.elements.connectKeplr);
  }

  async close() {
    await page.click("text=×");
  }
}

export const connectPopup = new ConnectPopup();
