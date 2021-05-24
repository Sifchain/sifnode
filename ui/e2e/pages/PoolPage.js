import { DEX_TARGET } from "../config";
import { tokenSelection } from "./components/TokenSelection";
import { GenericPage } from "./GenericPage";

export class PoolPage extends GenericPage {
  constructor() {
    super();
    this.el = {
      detailsPriceMessage: "[data-handle='details-price-message']",
      detailsMinimumReceived: "[data-handle='details-minimum-received']",
      detailsPriceImpact: "[data-handle='details-price-impact']",
      detailsLiquidityProviderFee:
        "[data-handle='details-liquidity-provider-fee']",
      addLiquidityButton: '[data-handle="add-liquidity-button"]',
      actionsButton: '[data-handle="actions-go"]',
      poolPricesForwardNumber: '[data-handle="pool-prices-forward-number"]',
      poolPricesForwardSymbols: '[data-handle="pool-prices-forward-symbols"]',
      poolPricesBackwardNumber: '[data-handle="pool-prices-backward-number"]',
      poolPricesBackwardSymbols: '[data-handle="pool-prices-backward-symbols"]',
      poolEstimatesForwardNumber:
        '[data-handle="pool-estimates-forward-number"]',
      poolEstimatesForwardSymbols:
        '[data-handle="pool-estimates-forward-symbols"]',
      poolEstimatesBackwardNumber:
        '[data-handle="pool-estimates-backward-number"]',
      poolEstimatesBackwardSymbols:
        '[data-handle="pool-estimates-backward-symbols"]',
      poolEstimatesShareNumber: '[data-handle="pool-estimates-share-number"]',
    };
  }

  async navigate() {
    await page.goto(`${DEX_TARGET}/#/pool`, { waitUntil: "domcontentloaded" });
  }

  async clickAddLiquidity() {
    await page.click(this.el.addLiquidityButton);
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

  async getTokenAValue() {
    return await tokenSelection.getTokenAValue();
  }

  async getTokenBValue() {
    return await tokenSelection.getTokenBValue();
  }

  async clickTokenAMax() {
    await tokenSelection.clickTokenAMax();
  }

  async verifyTokenBalance(token, expectedBalance) {
    await tokenSelection.verifyTokenBalance(token, expectedBalance);
  }

  async getActionsButtonText() {
    return await page.innerText(this.el.actionsButton);
  }

  async verifyPoolPrices({
    expForwardNumber,
    expForwardSymbols,
    expBackwardNumber,
    expBackwardSymbols,
  }) {
    expect(await page.innerText(this.el.poolPricesForwardNumber)).toBe(
      expForwardNumber,
    );
    expect(await page.innerText(this.el.poolPricesForwardSymbols)).toBe(
      expForwardSymbols,
    );
    expect(await page.innerText(this.el.poolPricesBackwardNumber)).toBe(
      expBackwardNumber,
    );
    expect(await page.innerText(this.el.poolPricesBackwardSymbols)).toBe(
      expBackwardSymbols,
    );
  }

  async verifyPoolEstimates({
    expForwardNumber,
    expForwardSymbols,
    expBackwardNumber,
    expBackwardSymbols,
    expShareNumber,
  }) {
    expect(await page.innerText(this.el.poolEstimatesForwardNumber)).toBe(
      expForwardNumber,
    );
    expect(await page.innerText(this.el.poolEstimatesForwardSymbols)).toBe(
      expForwardSymbols,
    );
    expect(await page.innerText(this.el.poolEstimatesBackwardNumber)).toBe(
      expBackwardNumber,
    );
    expect(await page.innerText(this.el.poolEstimatesBackwardSymbols)).toBe(
      expBackwardSymbols,
    );
    expect(await page.innerText(this.el.poolEstimatesShareNumber)).toBe(
      expShareNumber,
    );
  }

  async clickActionsGo() {
    await page.click(this.el.actionsButton);
  }
}

export const poolPage = new PoolPage();
