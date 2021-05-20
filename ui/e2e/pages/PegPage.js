import { DEX_TARGET } from "../config";
import expect from "expect-playwright";
import { dexHeader } from "./DexHeader";

export class PegPage {
  constructor() {
    this.el = {
      assetAmount: (asset) => `[data-handle='${asset}-row-amount']`,
    };
  }

  async navigate() {
    await page.goto(`${DEX_TARGET}/#/peg`, { waitUntil: "domcontentloaded" });
  }

  async openTab(tab) {
    if (tab.toLowerCase() === "native") {
      await page.click("[data-handle='native-tab']");
    } else {
      await page.click("[data-handle='external-tab']");
    }
  }

  // TODO: handling popup to be done outside of this page method
  async unpeg(asset, amount) {
    await page.click(`[data-handle='unpeg-${asset.toLowerCase()}']`);
    await page.click('[data-handle="peg-input"]');
    await page.fill('[data-handle="peg-input"]', amount);
    await page.click('button:has-text("Unpeg")');

    const [confirmPopup] = await Promise.all([
      context.waitForEvent("page"),
      page.click('button:has-text("Confirm Unpeg")'),
    ]);

    await Promise.all([
      confirmPopup.waitForEvent("close"),
      confirmPopup.click('button:has-text("Approve")'),
    ]);

    await page.waitForSelector("text=Transaction Submitted");
    await this.closeSubmissionWindow();
    await page.waitForTimeout(15000); // wait for sifnode to validate the tx TODO: replace this wait with some dynamic condition
  }

  async peg(asset, amount) {
    await page.click(`[data-handle='peg-${asset.toLowerCase()}']`);
    await page.click('[data-handle="peg-input"]');
    await page.fill('[data-handle="peg-input"]', amount);
    await page.click('button:has-text("Peg")');
  }

  async clickConfirmPeg() {
    await page.click('button:has-text("Confirm Peg")'); // this opens new notification page
  }

  async verifyAssetAmount(asset, expectedAmount) {
    // waitForSelector with state 'attached' is needed because the element is resolved as not visible
    // checked DOM and it looks visible. TODO: further investigate why this happens
    await page.waitForSelector(this.el.assetAmount(asset), {
      state: "attached",
    });
    const element = await page.$(this.el.assetAmount(asset));
    await expect(element).toHaveText(expectedAmount);
  }

  async closeSubmissionWindow() {
    await page.waitForTimeout(1000);
    await page.click("text=×"); // sometimes clicking 'x' doesn't close the window (even if Playwright says it clicked). Waiting a bit helps
  }

  async verifyTransactionPending(asset) {
    await expect(page).toHaveSelector(
      `${this.el.assetAmount(asset)} [data-handle='pending-tx-marker']`,
    );
  }
}

export const pegPage = new PegPage();
