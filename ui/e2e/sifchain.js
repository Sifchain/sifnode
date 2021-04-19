require("isomorphic-fetch");
import Web3 from "web3";
const web3 = new Web3("http://localhost:7545");

export async function getSifchainBalances(url, address, asset) {
  return fetch(`${url}/auth/accounts/${address}`)
    .then((res) => {
      return res.json();
    })
    .then((json) => {
      const coins = json.result.value.coins;
      let amount;
      for (const coin of coins) {
        if (coin.denom === asset) {
          console.log(coin.amount);
          amount = web3.utils.fromWei(coin.amount);
        }
      }
      return amount;
    })
    .catch((err) => {
      throw `getSifchainBalances(): ${err}`;
    });
}
