import { connectPopup } from "./pages/ConnectPopup";
import { dexHeader } from "./pages/components/DexHeader";
import { MetamaskConnectPage } from "./pages/MetamaskConnectPage";

export async function connectKeplrAccount() {
  await page.bringToFront();
  const newPage = await context.waitForEvent("page");
  await newPage.waitForLoadState("domcontentloaded");
  const [popup] = await Promise.all([
    context.waitForEvent("page"),
    newPage.click("text=Approve"),
  ]);
  await popup.waitForLoadState("domcontentloaded");
  await popup.click("text=Approve");
}

export async function connectMetaMaskAccount() {
  await page.bringToFront();
  await dexHeader.clickConnected();

  //   clicking 'Connect Metamask' opens page in a new tab so we need to retrieve it
  const [newPage] = await Promise.all([
    context.waitForEvent("page"),
    await connectPopup.clickConnectMetamask(),
  ]);
  const mmConnectPage = new MetamaskConnectPage(newPage);
  await mmConnectPage.clickNext();
  await mmConnectPage.clickConnect();
  await connectPopup.close();
}
