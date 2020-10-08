import { Context } from ".";
import { TokenAmount } from "../entities";

export default ({ api, store }: Context) => ({
  async setQuantityOfToken(tokenAmount: TokenAmount) {
    // IF LocalStorage[Transaction][SellToken]
    // USERINPUT: Number of Shares to Selll	-> LocalStorage: Transaction: Object: [sellToken: String, quantity: Number]
  },
});
