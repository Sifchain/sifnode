import { Context } from ".";
import { Token } from "../entities";

export default ({ api, store }: Context) => ({
  async selectToken(token: Token) {
    // "Get <LIST>"
    // USERINPUT: Select Token (From <LIST>)
    // PUT: Selected Sell Token	-> LocalStorage: Transaction: Object: [sellToken: String]
    // RENDER: TransactionWindow.vue
  },
});
