export class DexHeader {
  constructor() {
    this.elements = {
      connect: "[data-handle='button-connected']",
    };
  }

  async clickConnected() {
    await page.click(this.elements.connect);
  }
}

export const dexHeader = new DexHeader();
