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
const { poolPage } = require("./pages/PoolPage.js");
const { confirmSupplyModal } = require("./pages/ConfirmSupplyModal.js");

async function getInputValue(selector) {
  return await page.evaluate((el) => el.value, await page.$(selector));
}

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

it.only("adds liquidity", async () => {
  const tokenA = "ceth";

  await poolPage.navigate();

  await poolPage.clickAddLiquidity();

  await poolPage.selectTokenA(tokenA);
  await poolPage.fillTokenAValue("10");
  expect(await poolPage.getTokenBValue()).toBe("12048.19277");

  await poolPage.fillTokenBValue("10000");
  expect(await poolPage.getTokenAValue()).toBe("8.30000");

  await poolPage.clickTokenAMax();
  expect(await poolPage.getTokenAValue()).toBe("100.000000000000000000");
  expect(await poolPage.getTokenBValue()).toBe("120481.92771");

  expect((await poolPage.getActionsButtonText()).toUpperCase()).toBe(
    "INSUFFICIENT FUNDS",
  );

  await poolPage.clickTokenAMax();
  await poolPage.fillTokenAValue("5");

  expect(await poolPage.getTokenBValue()).toBe("6024.09639");

  expect((await poolPage.getActionsButtonText()).toUpperCase()).toBe(
    "ADD LIQUIDITY",
  );

  await poolPage.verifyPoolPrices({
    expForwardNumber: "0.000830",
    expForwardSymbols: "cETH per ROWAN",
    expBackwardNumber: "1204.819277",
    expBackwardSymbols: "ROWAN per cETH",
  });

  await poolPage.verifyPoolEstimates({
    expForwardNumber: "0.000830",
    expForwardSymbols: "CETH per ROWAN", // <-- this is a bug TODO: cETH
    expBackwardNumber: "1204.819277",
    expBackwardSymbols: "ROWAN per CETH", // <-- this is a bug TODO: cETH
    expShareNumber: "0.06%",
  });

  // click Add Liquidity
  await poolPage.clickActionsGo();

  // confirm add liquidity modal

  expect(await confirmSupplyModal.getTitle()).toBe("You are depositing");

  expect(
    prepareRowText(await confirmSupplyModal.getTokenADetailsPoolText()),
  ).toBe("cETH Deposited 5");

  expect(
    prepareRowText(await confirmSupplyModal.getTokenBDetailsPoolText()),
  ).toBe("ROWAN Deposited 6024.09639");

  expect(prepareRowText(await confirmSupplyModal.getRatesBPerARowText())).toBe(
    "Rates 1 cETH = 1204.81927711 ROWAN",
  );
  expect(prepareRowText(await confirmSupplyModal.getRatesAPerBRowText())).toBe(
    "1 ROWAN = 0.00083000 cETH",
  );
  expect(
    prepareRowText(await confirmSupplyModal.getRatesShareOfPoolText()),
  ).toBe("Share of Pool: 0.06%"); // TODO: remove ":"

  await confirmSupplyModal.clickConfirmSupply();

  expect(await confirmSupplyModal.getConfirmationWaitText()).toBe(
    "Supplying 5 ceth and 6024.09639 rowan",
  );

  // Confirm transaction popup
  await page.waitForTimeout(1000);
  await keplrNotificationPopup.navigate();
  await keplrNotificationPopup.clickApprove();
  await page.waitForTimeout(10000); // wait for blockchain to update...

  await page.pause();
  await page.click("text=×"); // TODO: put inside page object
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
  await page.waitForTimeout(1000);
  await keplrNotificationPopup.navigate();
  await keplrNotificationPopup.clickApprove();
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
