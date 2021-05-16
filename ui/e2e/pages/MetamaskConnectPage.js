export class MetamaskConnectPage {
  constructor(page) {
    this.page = page;
    this.el = {
      nextButton:
        "#app-content > div > div.main-container-wrapper > div > div.permissions-connect-choose-account > div.permissions-connect-choose-account__footer-container > div.permissions-connect-choose-account__bottom-buttons > button.button.btn-primary",
      connectButton:
        "#app-content > div > div.main-container-wrapper > div > div.page-container.permission-approval-container > div.permission-approval-container__footers > div.page-container__footer > footer > button.button.btn-primary.page-container__footer-button",
    };
  }

  async clickNext() {
    await this.page.click(this.el.nextButton);
  }

  async clickConnect() {
    await this.page.click(this.el.connectButton);
  }
}

// export const metamaskConnectPage = new MetamaskConnectPage()
