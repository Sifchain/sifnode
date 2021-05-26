import { DEX_TARGET } from "../config";
import { assertWaitedValue } from "../utils";
import { tokenSelection } from "./components/TokenSelection";
import { GenericPage } from "./GenericPage";

export class SwapPage extends GenericPage {
  constructor() {
    super();
    this.el = {
      detailsPriceMessage: "[data-handle='details-price-message']",
      detailsMinimumReceived: "[data-handle='details-minimum-received']",
      detailsPriceImpact: "[data-handle='details-price-impact']",
      detailsLiquidityProviderFee:
        "[data-handle='details-liquidity-provider-fee']",
      swapButton: 'button:has-text("Swap")',
    };
  }

  async navigate() {
    await page.goto(`${DEX_TARGET}/#/swap`, { waitUntil: "domcontentloaded" });
  }

  async selectTokenA(token) {
    await tokenSelection.selectTokenA(token);
  }

  async selectTokenB(token) {
    await tokenSelection.selectTokenB(token);
  }

  async fillTokenAValue(value) {
    await tokenSelection.fillTokenAValue(value);
  }

  async fillTokenBValue(value) {
    await tokenSelection.fillTokenBValue(value);
  }

  async verifyTokenAValue(expectedValue) {
    await page.click(tokenSelection.el.tokenInput("a"));
    await assertWaitedValue(tokenSelection.el.tokenInput("a"), expectedValue);
  }

  async verifyTokenBValue(expectedValue) {
    await page.click(tokenSelection.el.tokenInput("b"));
    await assertWaitedValue(tokenSelection.el.tokenInput("b"), expectedValue);
  }

  async clickTokenAMax() {
    await tokenSelection.clickTokenAMax();
  }

  async verifyTokenBalance(token, expectedBalance) {
    await tokenSelection.verifyTokenBalance(token, expectedBalance);
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
}

export const swapPage = new SwapPage();
