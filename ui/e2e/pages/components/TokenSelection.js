import { GenericPage } from "../GenericPage";

export class TokenSelection extends GenericPage {
  constructor() {
    super();
    this.el = {
      tokenSelect: (token) => `[data-handle='${token}-select-button']`,
      tokenDropdown: (ab) => `[data-handle='token-${ab}-select-button']`,
      tokenInput: (ab) => `[data-handle="token-${ab}-input"]`,
      tokenMaxButton: (ab) => `[data-handle="token-${ab}-max-button"]`,
      tokenBalanceLabel: (token) => `[data-handle="${token}-balance-label"]`,
    };
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

  async verifyTokenBalance(token, expectedBalance) {
    expect(await page.innerText(this.el.tokenBalanceLabel(token))).toBe(
      expectedBalance,
    );
  }
}

export const tokenSelection = new TokenSelection();
