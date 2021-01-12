import { ISifService } from "../../api/SifService";
import { sleep } from "./sleep";

export const createWaitForBalance = (sifService: ISifService) => {
  return async function checkBalance(
    symbol: string,
    expectedAmount: string,
    account: string,
    maxTries = 20 // 20 seconds
  ) {
    let latestBalance: string = "-no balance-";
    for (let i = 0; i < maxTries; i++) {
      await sleep(1000);

      const newBalance = (await sifService.getBalance(account)).find(
        (bal) => bal.asset.symbol === symbol
      );
      latestBalance = newBalance?.amount.toString() || "-no balance-";
      if (latestBalance === expectedAmount) {
        return newBalance;
      }
    }
    throw new Error(
      `Balance of ${expectedAmount} (${symbol}) was never realised. Last recorded balance was ${latestBalance} (${symbol})`
    );
  };
};
