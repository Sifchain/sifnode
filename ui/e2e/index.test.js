/**
 * THIS IS PROTOTYPAL
 *
 * TODO
 * ==============
 * Playwright object models https://playwright.dev/docs/pom for maniplating windows
 * - MetaMask class should represent the metamask popup
 * - Keplr class should represent Keplr popup
 * Clients for inspecting the blockchain
 * - sifchainBlockchainAccount - class should represent sifchain blockain
 * - ethereumBlockchainAccount - class should represent ethereumBlockchain
 */
require("@babel/polyfill");
const path = require("path");
const fs = require("fs");
const { chromium } = require("playwright");

// configs
const { DEX_TARGET, MM_CONFIG, KEPLR_CONFIG } = require("./config.js");
const keplrConfig = require("../core/src/config.localnet.json");

// extension
const {
  MetaMask,
  connectMmAccount,
  confirmTransaction,
} = require("./metamask.js");
const { importKeplrAccount, connectKeplrAccount } = require("./keplr");

// services
const { getSifchainBalances } = require("./sifchain.js");
const { getEthBalance, advanceEthBlocks } = require("./ethereum.js");
const { extractFile, getExtensionPage } = require("./utils");

async function getInputValue(page, selector) {
  return await page.evaluate((el) => el.value, await page.$(selector));
}

let browserContext;
let dexPage;

describe("connect to page", () => {
  beforeAll(async () => {
    // extract extension zips
    await extractExtensionPackages();
    const pathToKeplrExtension = path.join(__dirname, KEPLR_CONFIG.path);
    const pathToMmExtension = path.join(__dirname, MM_CONFIG.path);
    const userDataDir = path.join(__dirname, "./playwright");
    // need to rm userDataDir or else will store extension state
    if (fs.existsSync(userDataDir)) {
      fs.rmdirSync(userDataDir, { recursive: true });
    }

    browserContext = await chromium.launchPersistentContext(userDataDir, {
      // headless required with extensions. xvfb used for ci/cd
      headless: false,
      args: [
        `--disable-extensions-except=${pathToKeplrExtension},${pathToMmExtension}`,
        `--load-extension=${pathToKeplrExtension},${pathToMmExtension}`,
      ],
      // devtools: true,
    });

    // setup metamask
    const MM = new MetaMask(browserContext, MM_CONFIG);
    await MM.setup(browserContext);

    // setup keplr account
    const keplrPage = await browserContext.newPage();
    await keplrPage.goto(
      "chrome-extension://dmkamcknogkgcdfhhbddcghachkejeap/popup.html#/register",
    );
    await importKeplrAccount(keplrPage, KEPLR_CONFIG.options);
    // goto dex page
    dexPage = await browserContext.newPage();
    dexPage.setDefaultTimeout(60000);
    await dexPage.goto(DEX_TARGET, { waitUntil: "domcontentloaded" });
  });

  afterAll(async () => {
    browserContext.close();
  });

  it("connect to keplr, check balance", async () => {
    const cEthBalance = await getSifchainBalances(
      keplrConfig.sifApiUrl,
      KEPLR_CONFIG.options.address,
      "ceth",
    );

    await connectKeplrAccount(dexPage, browserContext);
    await dexPage.waitForTimeout(1000); // todo capture out extension page close event
    expect(
      (await dexPage.innerText("[data-handle='ceth-row-amount']")).trim(),
    ).toBe(Number(cEthBalance).toFixed(6));
  });

  it("connects to metamask, check balance", async () => {
    const mmEthBalance = await getEthBalance(MM_CONFIG.options.address);
    // connect wallet
    await connectMmAccount(dexPage, browserContext, MM_CONFIG.id);
    await dexPage.waitForTimeout(1000); // todo capture out extension page close event
    // click external tokens tab
    await dexPage.click("text=External Tokens");
    // expect
    expect(await dexPage.innerText("[data-handle='eth-row-amount']")).toBe(
      Number(mmEthBalance).toFixed(6),
    );
  });

  it("pegs", async () => {
    // XXX: This currently reuses the page from a previous test - this might be ok for now but we will probably want to provide that state some other way
    // assumes wallets connected
    const mmEthBalance = await getEthBalance(MM_CONFIG.options.address);
    const cEthBalance = await getSifchainBalances(
      keplrConfig.sifApiUrl,
      KEPLR_CONFIG.options.address,
      "ceth",
    );

    const pegAmount = "1";

    await dexPage.click("[data-handle='external-tab']");
    await dexPage.click("[data-handle='peg-eth']");
    await dexPage.click('[data-handle="peg-input"]');
    await dexPage.fill('[data-handle="peg-input"]', pegAmount);
    await dexPage.click('button:has-text("Peg")');
    await dexPage.click('button:has-text("Confirm Peg")');

    await confirmTransaction(dexPage, browserContext, pegAmount, MM_CONFIG.id);
    // move chain forward
    await advanceEthBlocks(52);

    const rowAmount = await dexPage.innerText(
      "[data-handle='ceth-row-amount']",
    );

    const expected = (Number(cEthBalance) + Number(pegAmount)).toFixed(6);

    expect(rowAmount.trim()).toBe(expected);
  });

  it("swaps", async () => {
    // Navigate to swap page
    await dexPage.goto(DEX_TARGET, {
      waitUntil: "domcontentloaded",
    });

    await dexPage.waitForTimeout(1000); // slowing down to avoid tokens not updating

    await dexPage.click("[data-handle='swap-page-button']");

    await dexPage.waitForTimeout(1000); // slowing down to avoid tokens not updating

    // Get values of token A and token B in account
    // Select Token A
    await dexPage.click("[data-handle='token-a-select-button']");
    await dexPage.click("[data-handle='cusdc-select-button']");
    // Select Token B
    await dexPage.click("[data-handle='token-b-select-button']");
    await dexPage.click("[data-handle='rowan-select-button']");
    // Input amount A
    await dexPage.click('[data-handle="token-a-input"]');
    await dexPage.fill('[data-handle="token-a-input"]', "100");

    expect(await getInputValue(dexPage, '[data-handle="token-b-input"]')).toBe(
      "99.99800003",
    );

    // Check expected output (XXX: hmmm - might have to pull in formulae from core??)

    // Input amount B
    await dexPage.click('[data-handle="token-b-input"]');
    await dexPage.fill('[data-handle="token-b-input"]', "100");

    expect(await getInputValue(dexPage, '[data-handle="token-a-input"]')).toBe(
      "100.00200005",
    );

    // Click max
    await dexPage.click("[data-handle='token-a-max-button']");

    // Check expected estimated values
    expect(await getInputValue(dexPage, '[data-handle="token-a-input"]')).toBe(
      "10000.0", // TODO: trim mantissa
    );
    expect(await getInputValue(dexPage, '[data-handle="token-b-input"]')).toBe(
      "9980.0299600499",
    );
    expect(
      await dexPage.innerText("[data-handle='details-price-message']"),
    ).toBe("0.998003 ROWAN per cUSDC");
    expect(
      await dexPage.innerText("[data-handle='details-minimum-received']"),
    ).toBe("9880.229660 ROWAN");
    expect(
      await dexPage.innerText("[data-handle='details-price-impact']"),
    ).toBe("0.10%");
    expect(
      await dexPage.innerText("[data-handle='details-liquidity-provider-fee']"),
    ).toBe("9.9800 ROWAN");

    // Input Amount A
    await dexPage.click('[data-handle="token-a-input"]');
    await dexPage.fill('[data-handle="token-a-input"]', "50");

    expect(await getInputValue(dexPage, '[data-handle="token-b-input"]')).toBe(
      "49.9995000037",
    );
    expect(
      await dexPage.innerText("[data-handle='details-price-message']"),
    ).toBe("0.999990 ROWAN per cUSDC");
    expect(
      await dexPage.innerText("[data-handle='details-minimum-received']"),
    ).toBe("49.499505 ROWAN");
    expect(
      await dexPage.innerText("[data-handle='details-price-impact']"),
    ).toBe("< 0.01%");
    expect(
      await dexPage.innerText("[data-handle='details-liquidity-provider-fee']"),
    ).toBe("0.00025 ROWAN");

    // Click Swap Button
    await dexPage.click('button:has-text("Swap")');

    // Confirm dialog shows the expected values
    expect(
      await dexPage.innerText(
        "[data-handle='confirm-swap-modal'] [data-handle='details-price-message']",
      ),
    ).toBe("0.999990 ROWAN per cUSDC");
    expect(
      await dexPage.innerText(
        "[data-handle='confirm-swap-modal'] [data-handle='details-minimum-received']",
      ),
    ).toBe("49.499505 ROWAN");
    expect(
      await dexPage.innerText(
        "[data-handle='confirm-swap-modal'] [data-handle='details-price-impact']",
      ),
    ).toBe("< 0.01%");
    expect(
      await dexPage.innerText(
        "[data-handle='confirm-swap-modal'] [data-handle='details-liquidity-provider-fee']",
      ),
    ).toBe("0.00025 ROWAN");

    await dexPage.click('button:has-text("Confirm Swap")');

    // Confirm transactioni popup

    const keplrPage = await getExtensionPage(browserContext, KEPLR_CONFIG.id);

    await keplrPage.waitForLoadState();
    await keplrPage.click("text=Approve");
    await keplrPage.waitForLoadState();
    await dexPage.waitForTimeout(10000); // wait for blockchain to update...

    // Wait for balances to be the amounts expected
    expect(await dexPage.innerText('[data-handle="swap-message"]')).toBe(
      "Swapped 50 cusdc for 49.9995000037 rowan",
    );

    await dexPage.click("[data-handle='modal-view-close']");

    expect(await dexPage.innerText('[data-handle="cusdc-balance-label"]')).toBe(
      "Balance: 9,950.00 cUSDC",
    );

    expect(await dexPage.innerText('[data-handle="rowan-balance-label"]')).toBe(
      "Balance: 10,050.00 ROWAN",
    );
  });

  it("adds liquidity", async () => {
    // Click pool page
    await dexPage.click('[data-handle="pool-page-button"]');

    // Click add liquidity button
    await dexPage.click('[data-handle="add-liquidity-button"]');

    // Select ceth
    await dexPage.click("[data-handle='token-a-select-button']");
    await dexPage.click("[data-handle='ceth-select-button']");

    await dexPage.click('[data-handle="token-a-input"]');
    await dexPage.fill('[data-handle="token-a-input"]', "10");

    expect(await getInputValue(dexPage, '[data-handle="token-b-input"]')).toBe(
      "12048.19277",
    );

    await dexPage.click('[data-handle="token-b-input"]');
    await dexPage.fill('[data-handle="token-b-input"]', "10000");

    await dexPage.click('[data-handle="token-a-input"]');

    expect(await getInputValue(dexPage, '[data-handle="token-a-input"]')).toBe(
      "8.30000",
    );

    await dexPage.click('[data-handle="token-a-max-button"]');

    expect(await getInputValue(dexPage, '[data-handle="token-a-input"]')).toBe(
      "101.000000000000000000",
    );

    await dexPage.click('[data-handle="token-b-input"]');

    expect(await getInputValue(dexPage, '[data-handle="token-b-input"]')).toBe(
      "121686.74699",
    );

    expect(
      (await dexPage.innerText('[data-handle="actions-go"]')).toUpperCase(),
    ).toBe("INSUFFICIENT FUNDS");

    await dexPage.click('[data-handle="token-a-max-button"]');
    await dexPage.fill('[data-handle="token-a-input"]', "5");

    expect(await getInputValue(dexPage, '[data-handle="token-b-input"]')).toBe(
      "6024.09639",
    );

    expect(
      (await dexPage.innerText('[data-handle="actions-go"]')).toUpperCase(),
    ).toBe("ADD LIQUIDITY");

    expect(
      await dexPage.innerText('[data-handle="pool-prices-forward-number"]'),
    ).toBe("0.000830");
    expect(
      await dexPage.innerText('[data-handle="pool-prices-forward-symbols"]'),
    ).toBe("cETH per ROWAN");

    expect(
      await dexPage.innerText('[data-handle="pool-prices-backward-number"]'),
    ).toBe("1204.819277");
    expect(
      await dexPage.innerText('[data-handle="pool-prices-backward-symbols"]'),
    ).toBe("ROWAN per cETH");

    expect(
      await dexPage.innerText('[data-handle="pool-estimates-forwards-number"]'),
    ).toBe("0.000830");
    expect(
      await dexPage.innerText(
        '[data-handle="pool-estimates-forwards-symbols"]',
      ),
    ).toBe("CETH per ROWAN"); // <-- this is a bug TODO: cETH

    expect(
      await dexPage.innerText(
        '[data-handle="pool-estimates-backwards-number"]',
      ),
    ).toBe("1204.819277");
    expect(
      await dexPage.innerText(
        '[data-handle="pool-estimates-backwards-symbols"]',
      ),
    ).toBe("ROWAN per CETH"); // <-- this is a bug TODO: cETH
    expect(
      await dexPage.innerText('[data-handle="pool-estimates-share-number"]'),
    ).toBe("0.06%");

    await dexPage.click('[data-handle="actions-go"]');

    expect(
      await dexPage.innerText('[data-handle="confirmation-modal-title"]'),
    ).toBe("You are depositing");

    expect(
      prepareRowText(
        await dexPage.innerText(
          '[data-handle="token-a-details-panel-pool-row"]',
        ),
      ),
    ).toBe("cETH Deposited 5");

    expect(
      prepareRowText(
        await dexPage.innerText(
          '[data-handle="token-b-details-panel-pool-row"]',
        ),
      ),
    ).toBe("ROWAN Deposited 6024.09639");

    expect(
      prepareRowText(
        await dexPage.innerText('[data-handle="real-b-per-a-row"]'),
      ),
    ).toBe("Rates 1 cETH = 1204.81927711 ROWAN");
    expect(
      prepareRowText(
        await dexPage.innerText('[data-handle="real-a-per-b-row"]'),
      ),
    ).toBe("1 ROWAN = 0.00083000 cETH");
    expect(
      prepareRowText(
        await dexPage.innerText('[data-handle="real-share-of-pool"]'),
      ),
    ).toBe("Share of Pool: 0.06%"); // TODO: remove ":"

    await dexPage.click("button:has-text('CONFIRM SUPPLY')");

    expect(
      await dexPage.innerText('[data-handle="confirmation-wait-message"]'),
    ).toBe("Supplying 5 ceth and 6024.09639 rowan");

    // Confirm transaction popup

    const keplrPage = await getExtensionPage(browserContext, KEPLR_CONFIG.id);

    await keplrPage.waitForLoadState();
    await keplrPage.click("text=Approve");
    await keplrPage.waitForLoadState();
    await dexPage.waitForTimeout(10000); // wait for blockchain to update...

    await dexPage.click("text=×");
    await dexPage.click('[data-handle="ceth-rowan-pool-list-item"]');

    expect(
      prepareRowText(
        await dexPage.innerText('[data-handle="total-pooled-ceth"]'),
      ),
    ).toBe("Total Pooled cETH: 8305.00000");

    expect(
      prepareRowText(
        await dexPage.innerText('[data-handle="total-pooled-rowan"]'),
      ),
    ).toBe("Total Pooled ROWAN: 10006024.09639");

    expect(
      prepareRowText(
        await dexPage.innerText('[data-handle="total-pool-share"]'),
      ),
    ).toBe("Your pool share: 0.0602 %");
  });
});

function prepareRowText(row) {
  return row
    .split("\n")
    .map((s) => s.trim())
    .filter(Boolean)
    .join(" ");
}

async function extractExtensionPackages() {
  await extractFile(`downloads/${KEPLR_CONFIG.id}.zip`, "./extensions");
  await extractFile(`downloads/${MM_CONFIG.id}.zip`, "./extensions");
  return;
}
