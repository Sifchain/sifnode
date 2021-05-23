import { DEX_TARGET } from "../config";
import { GenericPage } from "./GenericPage";

export class SwapPage extends GenericPage {
  constructor() {
    super();
    this.el = {
      tokenSelect: (token) => `[data-handle='${token}-select-button']`,
      tokenDropdown: (ab) => `[data-handle='token-${ab}-select-button']`,
      tokenInput: (ab) => `[data-handle="token-${ab}-input"]`,
      tokenMaxButton: (ab) => `[data-handle="token-${ab}-max-button"]`,
      detailsPriceMessage: "[data-handle='details-price-message']",
      detailsMinimumReceived: "[data-handle='details-minimum-received']",
      detailsPriceImpact: "[data-handle='details-price-impact']",
      detailsLiquidityProviderFee:
        "[data-handle='details-liquidity-provider-fee']",
      swapButton: 'button:has-text("Swap")',
      tokenBalanceLabel: (token) => `[data-handle="${token}-balance-label"]`,
    };
  }

  async navigate() {
    await page.goto(`${DEX_TARGET}/#/swap`, { waitUntil: "domcontentloaded" });
  }

  async selectTokenA(token) {
    await page.click(this.el.tokenDropdown("a"));
    await page.click(this.el.tokenSelect(token));
  }

  async selectTokenB(token) {
    await page.click(this.el.tokenDropdown("b"));
    await page.click(this.el.tokenSelect(token));
  }

  async fillTokenAValue(value) {
    await page.fill(this.el.tokenInput("a"), value);
  }

  async fillTokenBValue(value) {
    await page.fill(this.el.tokenInput("b"), value);
  }

  async getTokenAValue() {
    return await this.getInputValue(this.el.tokenInput("a"));
  }

  async getTokenBValue() {
    return await this.getInputValue(this.el.tokenInput("b"));
  }

  async clickTokenAMax() {
    await page.click(this.el.tokenMaxButton("a"));
  }

  async verifyDetails({
    expPriceMessage,
    expMinimumReceived,
    expPriceImpact,
    expLiquidityProviderFee,
  }) {
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

  async clickSwap() {
    await page.click(this.el.swapButton);
  }

  async verifyTokenBalance(token, expectedBalance) {
    expect(await page.innerText(this.el.tokenBalanceLabel(token))).toBe(
      expectedBalance,
    );
  }
}

export const swapPage = new SwapPage();
