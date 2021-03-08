const {chromium} = require("playwright");
const expect = require("expect");

// todo some kind of config file env mapping
const targetURL = "https://dex.sifchain.finance"
const sifAddress = "sif1k3ldrhrtkr6ncgjlu2ltyls9rxqtu9643lknrk"
const ethAddress = "0xc3b2058C07A865c50eA456131c739E305b000e7D"

let browser;
let page;

describe("connect to page", () => {
  beforeAll(async () => {
    browser = await chromium.launch();
  });
  afterAll(async () => {
    await browser.close();
  });
  beforeEach(async () => {
    page = await browser.newPage();
  });
  afterEach(async () => {
    await page.close();
  });
  
  it("should work", async () => {
    await page.goto(targetURL);
    expect(await page.title()).toBe("Peg Listing - Sifchain");
  });


})
