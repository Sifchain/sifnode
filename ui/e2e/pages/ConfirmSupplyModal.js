export class ConfirmSupplyModal {
  constructor() {
    this.el = {
      title: '[data-handle="confirmation-modal-title"]',
      tokenADetailsPoolRow: '[data-handle="token-a-details-panel-pool-row"]',
      tokenBDetailsPoolRow: '[data-handle="token-b-details-panel-pool-row"]',
      ratesAPerBRow: '[data-handle="real-a-per-b-row"]',
      ratesBPerARow: '[data-handle="real-b-per-a-row"]',
      ratesShareOfPool: '[data-handle="real-share-of-pool"]',
      confirmSupplyButton: 'button:has-text("Confirm Supply")',
      confirmationWaitMessage: '[data-handle="confirmation-wait-message"]',
    };
  }

  async clickConfirmSupply() {
    await page.click(this.el.confirmSupplyButton);
  }

  async getTitle() {
    return await page.innerText(this.el.title);
  }

  async getTokenADetailsPoolText() {
    return await page.innerText(this.el.tokenADetailsPoolRow);
  }

  async getTokenBDetailsPoolText() {
    return await page.innerText(this.el.tokenBDetailsPoolRow);
  }

  async getRatesAPerBRowText() {
    return await page.innerText(this.el.ratesAPerBRow);
  }

  async getRatesBPerARowText() {
    return await page.innerText(this.el.ratesBPerARow);
  }

  async getRatesShareOfPoolText() {
    return await page.innerText(this.el.ratesShareOfPool);
  }

  async getConfirmationWaitText() {
    return await page.innerText(this.el.confirmationWaitMessage);
  }
}

export const confirmSupplyModal = new ConfirmSupplyModal();
