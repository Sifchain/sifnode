export class ConfirmSupplyModal {
  constructor() {
    this.el = {
      title: '[data-handle="confirmation-modal-title"]',
      tokenInfoRow: (token) => `[data-handle="info-row-${token}"]`,
      tokenInfoAmount: (token) =>
        `[data-handle="info-row-${token}"] [data-handle="info-amount"]`,
      ratesAPerBRow: '[data-handle="real-a-per-b-row"]',
      ratesBPerARow: '[data-handle="real-b-per-a-row"]',
      ratesShareOfPool: '[data-handle="real-share-of-pool"]',
      confirmSupplyButton: 'button:has-text("Confirm Supply")',
      confirmationWaitMessage: '[data-handle="confirmation-wait-message"]',
      closeXButton: '[data-handle="modal-view-close"]',
    };
  }

  async clickConfirmSupply() {
    await page.click(this.el.confirmSupplyButton);
  }

  async getTitle() {
    return await page.innerText(this.el.title);
  }

  async getTokenInfoRowText(token) {
    return await page.innerText(this.el.tokenInfoRow(token));
  }

  async getTokenAmountText(token) {
    return await page.innerText(this.el.tokenInfoAmount(token));
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

  async closeModal() {
    await page.click(this.el.closeXButton);
  }
}

export const confirmSupplyModal = new ConfirmSupplyModal();
