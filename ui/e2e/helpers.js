import { connectPopup } from "./pages/ConnectPopup";
import { dexHeader } from "./pages/components/DexHeader";
import {
  metamaskConnectPage,
  MetamaskConnectPage,
} from "./pages/MetamaskConnectPage";
import { keplrNotificationPopup } from "./pages/KeplrNotificationPopup";
import urls from "./data/urls.json";

export async function connectKeplrAccount() {
  // it's not necessary to invoke connectPopup.clickConnectKeplr()
  // since keplrPage.setup() returns connect notification popup implicitly
  await keplrNotificationPopup.navigate(urls.keplr.notificationPopup.connect);
  await keplrNotificationPopup.clickApprove();
  // new page opens
  await page.waitForTimeout(1000);
  await keplrNotificationPopup.navigate(
    urls.keplr.notificationPopup.signinApprove,
  );
  await keplrNotificationPopup.clickApprove();
}

export async function reconnectKeplrAccount() {
  await dexHeader.clickConnected();
  await page.waitForTimeout(1000);
  const isConnected = await connectPopup.isKeplrConnected();
  if (!isConnected) await connectPopup.clickConnectKeplr();
}

export async function connectMetaMaskAccount() {
  await page.bringToFront();
  await dexHeader.clickConnected();
  await connectPopup.clickConnectMetamask();
  // opens new page so we need to retrieve it
  await page.waitForTimeout(4000); // TODO: replace explicit wait with dynamic waiting for page with given url
  await metamaskConnectPage.navigate();
  await metamaskConnectPage.clickNext();
  await metamaskConnectPage.clickConnect();
  await connectPopup.close();
}
