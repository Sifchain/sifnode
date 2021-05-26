/**
 * THIS IS PROTOTYPAL
 *
 * TODO
 * ==============
 * Clients for inspecting the blockchain
 * - sifchainBlockchainAccount - class should represent sifchain blockain
 * - ethereumBlockchainAccount - class should represent ethereumBlockchain
 */
require("@babel/polyfill");

// configs
const { MM_CONFIG, KEPLR_CONFIG } = require("./config.js");
const keplrConfig = require("../core/src/config.localnet.json");

// extension
const { metamaskPage } = require("./pages/MetaMaskPage");
const { keplrPage } = require("./pages/KeplrPage.js");
const { keplrNotificationPopup } = require("./pages/KeplrNotificationPopup.js");
const {
  metamaskNotificationPopup,
} = require("./pages/MetamaskNotificationPage.js");

// services
const { getSifchainBalances } = require("./sifchain.js");
const { advanceEthBlocks } = require("./ethereum.js");
const { extractExtensionPackage } = require("./utils");
const { useStack } = require("../test/stack");

// utils
const { connectMetaMaskAccount, connectKeplrAccount } = require("./helpers.js");

// dex pages
const { balancesPage } = require("./pages/BalancesPage.js");
const { swapPage } = require("./pages/SwapPage.js");
const { confirmSwapModal } = require("./pages/ConfirmSwapModal.js");
const { poolPage } = require("./pages/PoolPage.js");
const { confirmSupplyModal } = require("./pages/ConfirmSupplyModal.js");

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
  balancesPage.navigate();
  page.waitForTimeout(4000); // wait a second before keplr is finished being setup

  // Keplr will automatically connect and cause the add chain popup to come up
  await connectKeplrAccount();
  await connectMetaMaskAccount();
});

afterAll(async () => {
  await context.close();
});

beforeEach(async () => {
  await metamaskPage.reset();
});

it("imports rowan", async () => {
  const assetNative = "rowan";
  const exportAmount = "500";
  const assetExternal = "erowan";
  const importAmount = "100";

  // First we need to export rowan in order to have erowan on the bridgebank contract
  await balancesPage.navigate();

  await balancesPage.openTab("native");
  await balancesPage.export(assetNative, exportAmount);

  await balancesPage.openTab("external");
  // await page.pause();
  await balancesPage.verifyAssetAmount(assetExternal, "600.000000");

  // Now lets import erowan
  await balancesPage.import(assetExternal, importAmount);

  await balancesPage.clickConfirmImport();

  await page.waitForTimeout(500);
  await metamaskNotificationPopup.navigate();

  await page.waitForTimeout(1000);
  await metamaskNotificationPopup.clickViewFullTransactionDetails();
  await metamaskNotificationPopup.verifyTransactionDetails(
    `${importAmount} ${assetExternal}`,
  );

  await metamaskNotificationPopup.clickConfirm();
  await page.waitForTimeout(1000);

  await metamaskNotificationPopup.navigate(); // this call is needed to reload this.page with a new popup page
  await metamaskNotificationPopup.clickConfirm();

  await page.waitForTimeout(1000);
  await balancesPage.closeSubmissionWindow();
  // Check that tx marker for the tx is there
  await balancesPage.verifyTransactionPending(assetNative);

  // move chain forward
  await advanceEthBlocks(52);

  await page.waitForSelector("text=has succeded", { timeout: 60 * 1000 });

  await balancesPage.openTab("native");
  await balancesPage.verifyAssetAmount(assetNative, "9600.000000");
  // await balancesPage.verifyAssetAmount(assetNative, "96000.000000");
});

it("imports ether", async () => {
  const importAmount = "1";
  const assetExternal = "eth";
  const assetNative = "ceth";

  await balancesPage.navigate();

  const cEthBalance = await getSifchainBalances(
    keplrConfig.sifApiUrl,
    KEPLR_CONFIG.options.address,
    assetNative,
  );

  await balancesPage.openTab("external");
  await balancesPage.import(assetExternal, importAmount);

  await balancesPage.clickConfirmImport();
  await page.waitForTimeout(1000);

  await metamaskNotificationPopup.navigate();
  await metamaskNotificationPopup.clickConfirm();

  await page.waitForTimeout(1000);
  await balancesPage.closeSubmissionWindow();

  // Check that tx marker for the tx is there
  await balancesPage.verifyTransactionPending(assetNative);

  // move chain forward
  await advanceEthBlocks(52);

  await page.waitForSelector("text=has succeded", { timeout: 60 * 1000 });

  const expectedAmount = (Number(cEthBalance) + Number(importAmount)).toFixed(
    6,
  );
  await balancesPage.verifyAssetAmount(assetNative, expectedAmount);
});

it("imports tokens", async () => {
  const importAmount = "1";
  const importAsset = "usdc";

  await balancesPage.navigate();

  const cBalance = await getSifchainBalances(
    keplrConfig.sifApiUrl,
    KEPLR_CONFIG.options.address,
    "cusdc",
  );

  await balancesPage.openTab("external");
  await balancesPage.import(importAsset, importAmount);

  await balancesPage.clickConfirmImport();
  await page.waitForTimeout(1000);

  await metamaskNotificationPopup.navigate();
  await metamaskNotificationPopup.clickViewFullTransactionDetails();
  await metamaskNotificationPopup.verifyTransactionDetails(
    `${importAmount} ${importAsset}`,
  );
  await metamaskNotificationPopup.clickConfirm();
  await page.waitForTimeout(1000);
  await metamaskNotificationPopup.navigate(); // this call is needed to reload this.page with a new popup page
  await metamaskNotificationPopup.clickConfirm();

  await balancesPage.closeSubmissionWindow();

  await advanceEthBlocks(52);

  await page.waitForSelector("text=has succeded", { timeout: 60 * 1000 });

  const expectedAmount = (Number(cBalance) + Number(importAmount)).toFixed(6);
  await balancesPage.verifyAssetAmount("cusdc", expectedAmount);
});

it("swaps", async () => {
  const tokenA = "cusdc";
  const tokenB = "rowan";

  await swapPage.navigate();

  await swapPage.selectTokenA(tokenA);
  await page.waitForTimeout(1000); // slowing down to avoid tokens not updating
  await swapPage.selectTokenB(tokenB);

  await swapPage.fillTokenAValue("100");
  await swapPage.verifyTokenBValue("99.99800003");

  // Check expected output (XXX: hmmm - might have to pull in formulae from core??)

  await swapPage.fillTokenBValue("100");
  await swapPage.verifyTokenAValue("100.00200005");

  await swapPage.clickTokenAMax();
  await swapPage.verifyTokenAValue("10000.0"); // TODO: trim mantissa
  await swapPage.verifyTokenBValue("9980.0299600499");
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
  await swapPage.verifyTokenBValue("49.9995000037");

  await swapPage.clickSwap();

  // Confirm dialog shows the expected values
  // expect(
  //   await dexPage.innerText(
  //     "[data-handle='info-row-cusdc'] [data-handle='info-amount']",
  //   ),
  // ).toBe("50.000000");
  // expect(
  //   await dexPage.innerText(
  //     "[data-handle='info-row-rowan'] [data-handle='info-amount']",
  //   ),
  // ).toBe("49.999500");
  await confirmSwapModal.verifyDetails({
    expPriceMessage: "0.999990 ROWAN per cUSDC",
    expMinimumReceived: "49.499505 ROWAN",
    expPriceImpact: "< 0.01%",
    expLiquidityProviderFee: "0.00025 ROWAN",
  });

  await confirmSwapModal.clickConfirmSwap();

  // Confirm transaction popup
  await page.waitForTimeout(1000);
  await keplrNotificationPopup.navigate();
  await keplrNotificationPopup.clickApprove();
  await page.waitForTimeout(10000); // wait for blockchain to update...

  // Wait for balances to be the amounts expected
  await confirmSwapModal.verifySwapMessage(
    "Swapped 50 cUSDC for 49.9995000037 ROWAN",
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
  const tokenA = "ceth";
  const tokenB = "rowan";

  await poolPage.navigate();

  await poolPage.clickAddLiquidity();

  await poolPage.selectTokenA(tokenA);
  await poolPage.fillTokenAValue("10");
  await poolPage.verifyTokenBValue("12048.19277");

  await poolPage.fillTokenBValue("10000");
  await poolPage.verifyTokenAValue("8.30000");

  await poolPage.clickTokenAMax();
  await poolPage.verifyTokenAValue("100.000000000000000000");
  await poolPage.verifyTokenBValue("120481.92771");

  expect((await poolPage.getActionsButtonText()).toUpperCase()).toBe(
    "INSUFFICIENT FUNDS",
  );

  await poolPage.clickTokenAMax();
  await poolPage.fillTokenAValue("5");
  await poolPage.verifyTokenBValue("6024.09639");

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

  expect(await confirmSupplyModal.getTitle()).toBe("You are depositing");

  expect(
    prepareRowText(await confirmSupplyModal.getTokenADetailsPoolText()),
  ).toBe("cETH Deposited 5.00000");

  expect(
    prepareRowText(await confirmSupplyModal.getTokenBDetailsPoolText()),
  ).toBe("ROWAN Deposited 6024.096390");

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
    "Supplying 5.00000 ceth and 6024.09639 rowan",
  );

  // Confirm transaction popup
  await page.waitForTimeout(1000);
  await keplrNotificationPopup.navigate();
  await keplrNotificationPopup.clickApprove();
  await page.waitForTimeout(10000); // wait for blockchain to update...

  await confirmSupplyModal.closeModal();

  await poolPage.clickManagePool(tokenA, tokenB);
  expect(prepareRowText(await poolPage.getTotalPooledText(tokenA))).toBe(
    "Total Pooled cETH: 8305.00000",
  );
  expect(prepareRowText(await poolPage.getTotalPooledText(tokenB))).toBe(
    "Total Pooled ROWAN: 10006024.09639",
  );
  expect(prepareRowText(await poolPage.getTotalPoolShareText())).toBe(
    "Your Pool Share (%): 0.0602",
  );
});

it("fails to add liquidity when can't pay gas with rowan", async () => {
  const tokenA = "cusdc";

  await poolPage.navigate();
  await poolPage.clickAddLiquidity();
  await poolPage.selectTokenA(tokenA);
  await poolPage.fillTokenBValue("10000");
  await poolPage.clickActionsGo();
  await confirmSupplyModal.clickConfirmSupply();

  // Confirm transaction popup
  await page.waitForTimeout(1000);
  await keplrNotificationPopup.navigate();
  await keplrNotificationPopup.clickApprove();
  await page.waitForTimeout(10000); // wait for blockchain to update...

  await expect(page).toHaveText("Transaction Failed");
  await expect(page).toHaveText("Not enough ROWAN to cover the gas fees");

  await confirmSupplyModal.closeModal();
});

it("formats long amounts in confirmation screen", async () => {
  // Navigate to swap page
  await dexPage.goto(DEX_TARGET, {
    waitUntil: "domcontentloaded",
  });
  // Click pool page
  await dexPage.click('[data-handle="pool-page-button"]');

  // Click add liquidity button
  await dexPage.click('[data-handle="add-liquidity-button"]');

  // Select ceth
  await dexPage.click("[data-handle='token-a-select-button']");
  await dexPage.click("[data-handle='ceth-select-button']");

  await dexPage.click('[data-handle="token-a-input"]');
  await dexPage.fill(
    '[data-handle="token-a-input"]',
    "1.00000000000000000000000000000",
  );

  await dexPage.click('[data-handle="actions-go"]');

  expect(
    await dexPage.innerText(
      '[data-handle="token-a-row"] [data-handle="info-amount"]',
    ),
  ).toEqual("1.000000");
});

function prepareRowText(row) {
  return row
    .split("\n")
    .map((s) => s.trim())
    .filter(Boolean)
    .join(" ");
}
