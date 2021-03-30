// models/metamask.js

class MetaMask {
  constructor(page) {
    this.page = page;
  }
  async setup(browserContext, options) {
    const page = await getHomeScreen(browserContext)
    await this.confirmWelcomeScreen(page)
    await importAccount(page, options)
  }
  async confirmWelcomeScreen(metamaskPage) {
    await metamaskPage.click('.welcome-page button')
  }

}
module.exports = { MetaMask };

async function getHomeScreen(browserContext) {
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
  
  async function importAccount(metamaskPage, options) {
      console.log(metamaskPage)

    await metamaskPage.goto('chrome-extension://nkbihfbeogaeaoehlefnkodbefgpgknn/home.html#initialize/create-password/import-with-seed-phrase')
  
    // await metamaskPage.click('.metametrics-opt-in button.btn-primary');
  
    await metamaskPage.pause()

    const showSeedPhraseInput = await metamaskPage.waitForSelector('#ftf-chk1-label');
    await showSeedPhraseInput.click();
  
    const seedPhraseInput = await metamaskPage.waitForSelector('.first-time-flow textarea');
    await seedPhraseInput.type(options.mnemonic);
  
    const passwordInput = await metamaskPage.waitForSelector('#password');
    await passwordInput.type(options.password);
  
    const passwordConfirmInput = await metamaskPage.waitForSelector('#confirm-password');
    await passwordConfirmInput.type(options.password);
  
    const acceptTerms = await metamaskPage.waitForSelector('.first-time-flow__terms');
    await acceptTerms.click();
  
    const restoreButton = await metamaskPage.waitForSelector('.first-time-flow button');
    await restoreButton.click();
  
    const doneButton = await metamaskPage.waitForSelector('.end-of-flow button');
    await doneButton.click();
  
    const popupButton = await metamaskPage.waitForSelector('.popover-header__button');
    await popupButton.click();
  }