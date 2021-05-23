export class GenericPage {
  constructor() {}
  async getInputValue(selector) {
    return await page.$eval(selector, (el) => el.value);
    // return await page.evaluate((el) => el.value, await page.$(selector));
  }
}
