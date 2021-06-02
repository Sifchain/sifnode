export class GenericPage {
  constructor() {}
  async getInputValue(selector) {
    return await page.$eval(selector, (el) => el.value);
  }
}
