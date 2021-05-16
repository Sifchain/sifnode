import { DEX_TARGET } from "../config";

const retry = require("retry-assert");

export class PegPage {
  constructor() {
    this.el = {
      assetAmount: (asset) => `${asset}-row-amount`,
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
    await page.click("text=×");
    await page.waitForTimeout(10000); // wait for sifnode to validate the tx
  }

  async peg(asset, amount) {
    await page.click(`[data-handle='peg-${asset.toLowerCase()}']`);
    await page.click('[data-handle="peg-input"]');
    await page.fill('[data-handle="peg-input"]', amount);
    await page.click('button:has-text("Peg")');

    const [approveSpendPopup] = await Promise.all([
      context.waitForEvent("page"),
      page.click('button:has-text("Confirm Peg")'),
    ]);

    await approveSpendPopup.click("text=View full transaction details");
    await expect(approveSpendPopup).toHaveText(
      amount + ` ${asset.toLowerCase()}`,
    );

    // TODO: abstract away confirmation flow
    const [confirmPopup2] = await Promise.all([
      context.waitForEvent("page"),
      approveSpendPopup.click('button:has-text("Confirm")'),
    ]);

    await Promise.all([
      confirmPopup2.waitForEvent("close"),
      confirmPopup2.click('button:has-text("Confirm")'),
    ]);

    await page.click("text=×");
  }

  async verifyAssetAmount(asset, expectedAmount) {
    const rowAmount = await retry()
      .fn(() => page.innerText(`[data-handle='${asset}-row-amount']`))
      .withTimeout(10000)
      .until((rowAmount) => expect(rowAmount.trim()).toBe(expectedAmount));
  }
}

export const pegPage = new PegPage();
