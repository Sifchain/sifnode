const expectPw = require("expect-playwright");

export class ConfirmSwapModal {
  constructor() {
    this.rootSelector = '[data-handle="confirm-swap-modal"]';
    this.el = {
      tokenInfoRow: (token) =>
        `${this.rootSelector} [data-handle="info-row-${token}"]`,
      tokenInfoAmount: (token) =>
        `${this.rootSelector} [data-handle="info-row-${token}"] [data-handle="info-amount"]`,
      detailsPriceMessage: `${this.rootSelector} [data-handle='details-price-message']`,
      detailsMinimumReceived: `${this.rootSelector} [data-handle='details-minimum-received']`,
      detailsPriceImpact: `${this.rootSelector} [data-handle='details-price-impact']`,
      detailsLiquidityProviderFee: `${this.rootSelector} [data-handle='details-liquidity-provider-fee']`,
      confirmSwapButton: `${this.rootSelector} button:has-text("Confirm Swap")`,
      closeButton: 'button:has-text("Close")', // TODO: these 2 selectors actually belong to a different modal but functionally they fit here, consider moving to a separate class
      swapMessage: '[data-handle="swap-message"]',
      closeXButton: '[data-handle="modal-view-close"]',
    };
  }

  async verifyDetails({
    tokenA,
    tokenB,
    expTokenAAmount,
    expTokenBAmount,
    expPriceMessage,
    expMinimumReceived,
    expPriceImpact,
    expLiquidityProviderFee,
  }) {
    expect(await page.innerText(this.el.tokenInfoAmount(tokenA))).toBe(
      expTokenAAmount,
    );
    expect(await page.innerText(this.el.tokenInfoAmount(tokenB))).toBe(
      expTokenBAmount,
    );
    expect(await page.innerText(this.el.detailsPriceMessage)).toBe(
      expPriceMessage,
    );
    expect(await page.innerText(this.el.detailsMinimumReceived)).toBe(
      expMinimumReceived,
    );
    expect(await page.innerText(this.el.detailsPriceImpact)).toBe(
      expPriceImpact,
    );
    expect(await page.innerText(this.el.detailsLiquidityProviderFee)).toBe(
      expLiquidityProviderFee,
    );
  }

  async clickConfirmSwap() {
    await page.click(this.el.confirmSwapButton);
  }

  async clickClose() {
    await page.click(this.el.closeButton);
  }

  async closeModal() {
    await page.click(this.el.closeXButton);
  }

  async getSwapMessage() {
    return await page.innerText(this.el.swapMessage);
  }

  async getTokenInfoRowText(token) {
    return await page.innerText(this.el.tokenInfoRow(token));
  }

  async getTokenAmountText(token) {
    return await page.innerText(this.el.tokenInfoAmount(token));
  }

  async verifySwapMessage(expectedMessage) {
    await expectPw(page).toHaveText(this.el.swapMessage, expectedMessage);
  }
}

export const confirmSwapModal = new ConfirmSwapModal();
