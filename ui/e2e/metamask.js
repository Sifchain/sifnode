// models/metamask.js

class MetaMask {
  constructor(page, config) {
    this.page = page;
    this.config = config;
  }
  async setup(browserContext) {
    const mmPage = await getHomePage(browserContext)
    await confirmWelcomeScreen(mmPage)
    await importAccount(mmPage, this.config)
    await addNetwork(mmPage, this.config)
  }
}

async function getHomePage(browserContext) {
  return new Promise((resolve, reject) => {
    browserContext.waitForEvent('page', async (page) => {
      if (page.url().match('chrome-extension://[a-z]+/home.html')) {
        try {
          resolve(page);
        } catch (e) {
          reject(e);
        }
      }
    });
  });
}

async function confirmWelcomeScreen(mmPage) {
  await mmPage.click('.welcome-page button')
}

async function importAccount(mmPage, config) {
  await mmPage.goto(`chrome-extension://${config.id}/home.html#initialize/create-password/import-with-seed-phrase`)
  await mmPage.type('.first-time-flow__seedphrase input', config.options.mnemonic);
  await mmPage.type('#password', config.options.password);
  await mmPage.type('#confirm-password', config.options.password);
  await mmPage.click('.first-time-flow__terms');
  await mmPage.click('.first-time-flow button');
  await mmPage.click('.end-of-flow button');
  await mmPage.click('.popover-header__button');
}
async function addNetwork(mmPage, config) {
  await mmPage.goto(`chrome-extension://${config.id}/home.html#settings/networks`)
  await mmPage.click("#app-content > div > div.main-container-wrapper > div > div.settings-page__content > div.settings-page__content__modules > div > div.settings-page__sub-header > div > button")
  await mmPage.type('#network-name', config.network.name);
  await mmPage.type("#rpc-url", `http://localhost:${config.network.port}`)
  await mmPage.type("#chainId", config.network.chainId)
  await mmPage.click("#app-content > div > div.main-container-wrapper > div > div.settings-page__content > div.settings-page__content__modules > div > div.networks-tab__content > div.networks-tab__network-form > div.network-form__footer > button.button.btn-secondary")
  await mmPage.click("#app-content > div > div.main-container-wrapper > div > div.settings-page__header > div.settings-page__close-button")
}

 async function connectMmAccount(page, browserContext) {
 // Click button:has-text("Not connected")
 await page.click('button:has-text("Not connected")');
 // Click button:has-text("Connect Metamask")
 await page.click('button:has-text("Connect Metamask")');
 await page.pause()

 // Open new page
 const page4 = await browserContext.newPage();
 page4.goto('chrome-extension://nkbihfbeogaeaoehlefnkodbefgpgknn/notification.html');
 // Go to chrome-extension://nkbihfbeogaeaoehlefnkodbefgpgknn/notification.html#connect/7e4be24d-a014-4d43-9622-c8387fdebfc9
 await page4.goto('chrome-extension://nkbihfbeogaeaoehlefnkodbefgpgknn/notification.html#connect/7e4be24d-a014-4d43-9622-c8387fdebfc9');
 // Go to chrome-extension://nkbihfbeogaeaoehlefnkodbefgpgknn/notification.html#connect/7e4be24d-a014-4d43-9622-c8387fdebfc9/confirm-permissions
 await page4.goto('chrome-extension://nkbihfbeogaeaoehlefnkodbefgpgknn/notification.html#connect/7e4be24d-a014-4d43-9622-c8387fdebfc9/confirm-permissions');
 // Go to chrome-extension://nkbihfbeogaeaoehlefnkodbefgpgknn/notification.html#
 await page4.goto('chrome-extension://nkbihfbeogaeaoehlefnkodbefgpgknn/notification.html#');
 // Close page
 await page4.close();
 // Click text=×
 await page.click('text=×');
 await page.pause()
}

module.exports = { MetaMask, connectMmAccount };
