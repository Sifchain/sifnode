export class ApproveSpendPopup {
  constructor(page) {
    this.page = page;
  }

  async clickViewFullTransactionDetails() {
    await this.page.click("text=View full transaction details");
  }

  async clickConfirm() {
    await this.page.click('button:has-text("Confirm")');
  }
}
