import { ISifService } from "../../api/SifService";
import { sleep } from "./sleep";

export const createWaitForBalance = (sifService: ISifService) => {
  return async function checkBalance(
    symbol: string,
    expectedAmount: string,
    account: string,
    maxTries = 100
  ) {
    for (let i = 0; i < maxTries; i++) {
      await sleep(1000);

      const newBalance = (await sifService.getBalance(account)).find(
        (bal) => bal.asset.symbol === symbol
      );

      if (newBalance?.amount.toString() === expectedAmount) {
        return newBalance;
      }
    }
    throw new Error(`Balance of ${expectedAmount} was never realised`);
  };
};
