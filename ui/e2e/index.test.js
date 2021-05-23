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

// configs
const { DEX_TARGET, MM_CONFIG, KEPLR_CONFIG } = require("./config.js");
const keplrConfig = require("../core/src/config.localnet.json");

// extension
const { metamaskPage } = require("./pages/MetaMaskPage");

// services
const { getSifchainBalances } = require("./sifchain.js");
const { getEthBalance, advanceEthBlocks } = require("./ethereum.js");
const { getExtensionPage, extractExtensionPackage } = require("./utils");
const { useStack } = require("../test/stack");
const { keplrPage } = require("./pages/KeplrPage.js");
const { connectMetaMaskAccount, connectKeplrAccount } = require("./helpers.js");
const { pegPage } = require("./pages/PegPage.js");
const {
  metamaskNotificationPopup,
} = require("./pages/MetamaskNotificationPage.js");
const { swapPage } = require("./pages/SwapPage.js");
const { keplrNotificationPopup } = require("./pages/KeplrNotificationPopup.js");
const { confirmSwapModal } = require("./pages/ConfirmSwapModal.js");

// async function getInputValue(selector) {
//   return await page.evaluate((el) => el.value, await page.$(selector));
// }

let dexPage;

useStack("every-test");

beforeAll(async () => {
  // extract extension zips
  await extractExtensionPackage(MM_CONFIG.id);
  await extractExtensionPackage(KEPLR_CONFIG.id);

  await metamaskPage.navigate();
  await metamaskPage.setup();

  await keplrPage.navigate();
  await keplrPage.setup();

  // goto dex page
  pegPage.navigate();
  page.waitForTimeout(4000); // wait a second before keplr is finished being setup

  // Keplr will automatically connect and cause the add chain popup to come up
  await connectKeplrAccount();
  await connectMetaMaskAccount();
});

afterAll(async () => {
  context.close();
});

beforeEach(async () => {
  // metamaskPage.reset();
});

it("pegs rowan", async () => {
  const assetNative = "rowan";
  const unpegAmount = "500";
  const assetExternal = "erowan";
  const pegAmount = "100";

  // First we need to unpeg rowan in order to have erowan on the bridgebank contract
  await pegPage.navigate();

  await pegPage.openTab("native");
  await pegPage.unpeg(assetNative, unpegAmount);

  await pegPage.openTab("external");
  await pegPage.verifyAssetAmount(assetExternal, "600.000000");

  // Now lets peg erowan
  await pegPage.peg(assetExternal, pegAmount);

  await pegPage.clickConfirmPeg();

  await page.waitForTimeout(500);
  await metamaskNotificationPopup.navigate();

  await page.waitForTimeout(1000);
  await metamaskNotificationPopup.clickViewFullTransactionDetails();
  await metamaskNotificationPopup.verifyTransactionDetails(
    `${pegAmount} ${assetExternal}`,
  );

  await metamaskNotificationPopup.clickConfirm();
  await page.waitForTimeout(1000);

  await metamaskNotificationPopup.navigate(); // this call is needed to reload this.page with a new popup page
  await metamaskNotificationPopup.clickConfirm();

  await page.waitForTimeout(1000);
  await pegPage.closeSubmissionWindow();
  // Check that tx marker for the tx is there
  await pegPage.verifyTransactionPending(assetNative);

  // move chain forward
  await advanceEthBlocks(52);

  await page.waitForSelector("text=has succeded", { timeout: 60 * 1000 });

  await pegPage.openTab("native");
  await pegPage.verifyAssetAmount(assetNative, "9600.000000");
});

it("pegs ether", async () => {
  const pegAmount = "1";
  const assetExternal = "eth";
  const assetNative = "ceth";

  // Navigate to peg page
  await pegPage.navigate();

  const cEthBalance = await getSifchainBalances(
    keplrConfig.sifApiUrl,
    KEPLR_CONFIG.options.address,
    assetNative,
  );

  await pegPage.openTab("external");
  await pegPage.peg(assetExternal, pegAmount);

  await pegPage.clickConfirmPeg();
  await page.waitForTimeout(1000);

  await metamaskNotificationPopup.navigate();
  await metamaskNotificationPopup.clickConfirm();

  await page.waitForTimeout(1000);
  await pegPage.closeSubmissionWindow();

  // Check that tx marker for the tx is there
  await pegPage.verifyTransactionPending(assetNative);

  // move chain forward
  await advanceEthBlocks(52);

  await page.waitForSelector("text=has succeded", { timeout: 60 * 1000 });

  const expectedAmount = (Number(cEthBalance) + Number(pegAmount)).toFixed(6);
  await pegPage.verifyAssetAmount(assetNative, expectedAmount);
});

it("pegs tokens", async () => {
  const pegAmount = "1";
  const pegAsset = "usdc";

  // Navigate to peg page
  await pegPage.navigate();

  const cBalance = await getSifchainBalances(
    keplrConfig.sifApiUrl,
    KEPLR_CONFIG.options.address,
    "cusdc",
  );

  await pegPage.openTab("external");
  await pegPage.peg(pegAsset, pegAmount);

  await pegPage.clickConfirmPeg();
  await page.waitForTimeout(1000);

  await metamaskNotificationPopup.navigate();
  await metamaskNotificationPopup.clickViewFullTransactionDetails();
  await metamaskNotificationPopup.verifyTransactionDetails(
    `${pegAmount} ${pegAsset}`,
  );
  await metamaskNotificationPopup.clickConfirm();
  await page.waitForTimeout(1000);
  await metamaskNotificationPopup.navigate(); // this call is needed to reload this.page with a new popup page
  await metamaskNotificationPopup.clickConfirm();

  await pegPage.closeSubmissionWindow();

  await advanceEthBlocks(52);

  await page.waitForSelector("text=has succeded", { timeout: 60 * 1000 });

  const expectedAmount = (Number(cBalance) + Number(pegAmount)).toFixed(6);
  await pegPage.verifyAssetAmount("cusdc", expectedAmount);
});

it("swaps", async () => {
  const tokenA = "cusdc";
  const tokenB = "rowan";
  // Navigate to swap page
  await swapPage.navigate();

  await swapPage.selectTokenA(tokenA);
  await page.waitForTimeout(1000); // slowing down to avoid tokens not updating
  await swapPage.selectTokenB(tokenB);

  await swapPage.fillTokenAValue("100");
  expect(await swapPage.getTokenBValue()).toBe("99.99800003");

  // Check expected output (XXX: hmmm - might have to pull in formulae from core??)

  await swapPage.fillTokenBValue("100");
  expect(await swapPage.getTokenAValue()).toBe("100.00200005");

  await swapPage.clickTokenAMax();

  expect(await swapPage.getTokenAValue()).toBe(
    "10000.0", // TODO: trim mantissa
  );
  expect(await swapPage.getTokenBValue()).toBe("9980.0299600499");
  await swapPage.verifyDetails({
    expPriceMessage: "0.998003 ROWAN per cUSDC",
    expMinimumReceived: "9880.229660 ROWAN",
    expPriceImpact: "0.10%",
    expLiquidityProviderFee: "9.9800 ROWAN",
  });

  // Input Amount A
  await swapPage.fillTokenAValue("50");
  await swapPage.verifyDetails({
    expPriceMessage: "0.999990 ROWAN per cUSDC",
    expMinimumReceived: "49.499505 ROWAN",
    expPriceImpact: "< 0.01%",
    expLiquidityProviderFee: "0.00025 ROWAN",
  });
  expect(await swapPage.getTokenBValue()).toBe("49.9995000037");

  await swapPage.clickSwap();

  // Confirm dialog shows the expected values
  await confirmSwapModal.verifyDetails({
    expPriceMessage: "0.999990 ROWAN per cUSDC",
    expMinimumReceived: "49.499505 ROWAN",
    expPriceImpact: "< 0.01%",
    expLiquidityProviderFee: "0.00025 ROWAN",
  });

  await confirmSwapModal.clickConfirmSwap();

  // Confirm transactioni popup
  await page.waitForTimeout(1000);
  await keplrNotificationPopup.navigate();
  await keplrNotificationPopup.clickApprove();
  await page.waitForTimeout(10000); // wait for blockchain to update...

  // Wait for balances to be the amounts expected
  await confirmSwapModal.verifySwapMessage(
    "Swapped 50 cusdc for 49.9995000037 rowan",
  );

  await confirmSwapModal.clickClose();

  await swapPage.verifyTokenBalance(tokenA, "Balance: 9,950.00 cUSDC");
  await swapPage.verifyTokenBalance(tokenB, "Balance: 10,050.00 ROWAN");
});

it("fails to swap when it can't pay gas with rowan", async () => {
  const tokenA = "rowan";
  const tokenB = "cusdc";
  // Navigate to swap page
  await swapPage.navigate();

  // Get values of token A and token B in account
  await swapPage.selectTokenA(tokenA);
  await page.waitForTimeout(1000); // slowing down to avoid tokens not updating
  await swapPage.selectTokenB(tokenB);

  await swapPage.fillTokenAValue("10000");
  await swapPage.clickSwap();
  await confirmSwapModal.clickConfirmSwap();

  // Confirm transaction popup
  await page.waitForTimeout(1000);
  await keplrNotificationPopup.navigate();
  await keplrNotificationPopup.clickApprove();

  await page.waitForTimeout(10000); // wait for blockchain to update...

  await expect(page).toHaveText("Transaction Failed");
  await expect(page).toHaveText("Not enough ROWAN to cover the gas fees.");

  await confirmSwapModal.closeModal();
});

it("adds liquidity", async () => {
  // Navigate to swap page
  await page.goto(DEX_TARGET, {
    waitUntil: "domcontentloaded",
  });
  // Click pool page
  await page.click('[data-handle="pool-page-button"]');

  // Click add liquidity button
  await page.click('[data-handle="add-liquidity-button"]');

  // Select ceth
  await page.click("[data-handle='token-a-select-button']");
  await page.click("[data-handle='ceth-select-button']");

  await page.click('[data-handle="token-a-input"]');
  await page.fill('[data-handle="token-a-input"]', "10");

  expect(await getInputValue(page, '[data-handle="token-b-input"]')).toBe(
    "12048.19277",
  );

  await page.click('[data-handle="token-b-input"]');
  await page.fill('[data-handle="token-b-input"]', "10000");

  await page.click('[data-handle="token-a-input"]');

  expect(await getInputValue(page, '[data-handle="token-a-input"]')).toBe(
    "8.30000",
  );

  await page.click('[data-handle="token-a-max-button"]');

  expect(await getInputValue(page, '[data-handle="token-a-input"]')).toBe(
    "100.000000000000000000",
  );

  await page.click('[data-handle="token-b-input"]');

  expect(await getInputValue(page, '[data-handle="token-b-input"]')).toBe(
    "120481.92771",
  );

  expect(
    (await page.innerText('[data-handle="actions-go"]')).toUpperCase(),
  ).toBe("INSUFFICIENT FUNDS");

  await page.click('[data-handle="token-a-max-button"]');
  await page.fill('[data-handle="token-a-input"]', "5");

  expect(await getInputValue(page, '[data-handle="token-b-input"]')).toBe(
    "6024.09639",
  );

  expect(
    (await page.innerText('[data-handle="actions-go"]')).toUpperCase(),
  ).toBe("ADD LIQUIDITY");

  expect(
    await page.innerText('[data-handle="pool-prices-forward-number"]'),
  ).toBe("0.000830");
  expect(
    await page.innerText('[data-handle="pool-prices-forward-symbols"]'),
  ).toBe("cETH per ROWAN");

  expect(
    await page.innerText('[data-handle="pool-prices-backward-number"]'),
  ).toBe("1204.819277");
  expect(
    await page.innerText('[data-handle="pool-prices-backward-symbols"]'),
  ).toBe("ROWAN per cETH");

  expect(
    await page.innerText('[data-handle="pool-estimates-forwards-number"]'),
  ).toBe("0.000830");
  expect(
    await page.innerText('[data-handle="pool-estimates-forwards-symbols"]'),
  ).toBe("CETH per ROWAN"); // <-- this is a bug TODO: cETH

  expect(
    await page.innerText('[data-handle="pool-estimates-backwards-number"]'),
  ).toBe("1204.819277");
  expect(
    await page.innerText('[data-handle="pool-estimates-backwards-symbols"]'),
  ).toBe("ROWAN per CETH"); // <-- this is a bug TODO: cETH
  expect(
    await page.innerText('[data-handle="pool-estimates-share-number"]'),
  ).toBe("0.06%");

  await page.click('[data-handle="actions-go"]');

  expect(await page.innerText('[data-handle="confirmation-modal-title"]')).toBe(
    "You are depositing",
  );

  expect(
    prepareRowText(
      await page.innerText('[data-handle="token-a-details-panel-pool-row"]'),
    ),
  ).toBe("cETH Deposited 5");

  expect(
    prepareRowText(
      await page.innerText('[data-handle="token-b-details-panel-pool-row"]'),
    ),
  ).toBe("ROWAN Deposited 6024.09639");

  expect(
    prepareRowText(await page.innerText('[data-handle="real-b-per-a-row"]')),
  ).toBe("Rates 1 cETH = 1204.81927711 ROWAN");
  expect(
    prepareRowText(await page.innerText('[data-handle="real-a-per-b-row"]')),
  ).toBe("1 ROWAN = 0.00083000 cETH");
  expect(
    prepareRowText(await page.innerText('[data-handle="real-share-of-pool"]')),
  ).toBe("Share of Pool: 0.06%"); // TODO: remove ":"

  await page.click("button:has-text('CONFIRM SUPPLY')");

  expect(
    await page.innerText('[data-handle="confirmation-wait-message"]'),
  ).toBe("Supplying 5 ceth and 6024.09639 rowan");

  // Confirm transaction popup

  const keplrPage = await getExtensionPage(browserContext, KEPLR_CONFIG.id);

  await keplrPage.waitForLoadState();
  await keplrPage.click("text=Approve");
  await keplrPage.waitForLoadState();
  await page.waitForTimeout(10000); // wait for blockchain to update...

  await page.click("text=×");
  await page.click('[data-handle="ceth-rowan-pool-list-item"]');

  expect(
    prepareRowText(await page.innerText('[data-handle="total-pooled-ceth"]')),
  ).toBe("Total Pooled cETH: 8305.00000");

  expect(
    prepareRowText(await page.innerText('[data-handle="total-pooled-rowan"]')),
  ).toBe("Total Pooled ROWAN: 10006024.09639");

  expect(
    prepareRowText(await page.innerText('[data-handle="total-pool-share"]')),
  ).toBe("Your pool share: 0.0602 %");
});

it("fails to add liquidity when can't pay gas with rowan", async () => {
  // Navigate to swap page
  await page.goto(DEX_TARGET, {
    waitUntil: "domcontentloaded",
  });
  // Click pool page
  await page.click('[data-handle="pool-page-button"]');

  // Click add liquidity button
  await page.click('[data-handle="add-liquidity-button"]');

  // Select cusdc
  await page.click("[data-handle='token-a-select-button']");
  await page.click("[data-handle='cusdc-select-button']");

  await page.click('[data-handle="token-b-input"]');
  await page.fill('[data-handle="token-b-input"]', "10000");

  await page.click('[data-handle="actions-go"]');

  await page.click("button:has-text('CONFIRM SUPPLY')");

  // Confirm transaction popup

  const keplrPage = await getExtensionPage(browserContext, KEPLR_CONFIG.id);

  await keplrPage.waitForLoadState();
  await keplrPage.click("text=Approve");
  await keplrPage.waitForLoadState();
  await page.waitForTimeout(10000); // wait for blockchain to update...

  await expect(page).toHaveText("Transaction Failed");
  await expect(page).toHaveText("Not enough ROWAN to cover the gas fees");

  await page.click("[data-handle='modal-view-close']");
});

function prepareRowText(row) {
  return row
    .split("\n")
    .map((s) => s.trim())
    .filter(Boolean)
    .join(" ");
}
