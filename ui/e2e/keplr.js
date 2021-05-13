export async function importKeplrAccount(page, options) {
  await page.click("text=Import existing account");
  await page.click('textarea[name="words"]');
  await page.fill('textarea[name="words"]', options.mnemonic);
  await page.click('input[name="name"]');
  await page.fill('input[name="name"]', options.name);
  await page.click('input[name="password"]');
  await page.fill('input[name="password"]', "juniper21");
  await page.click('input[name="confirmPassword"]');
  await page.fill('input[name="confirmPassword"]', "juniper21");
  await page.click("text=Next");
  await page.click("text=Done");
}

export async function connectKeplrAccount(page, browserContext) {
  const newPage = await browserContext.waitForEvent("page");
  await newPage.waitForLoadState("domcontentloaded");
  const [popup] = await Promise.all([
    browserContext.waitForEvent("page"),
    newPage.click("text=Approve"),
  ]);
  await popup.waitForLoadState("domcontentloaded");
  await popup.click("text=Approve");
}
