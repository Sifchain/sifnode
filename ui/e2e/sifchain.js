require("isomorphic-fetch");

export async function getSifchainBalances(url, address, asset) {
  return fetch(`${url}/auth/accounts/${address}`)
    .then(res => {
      return res.json();
    })
    .then(json => {
      const coins = json.result.value.coins;
      let amount;
      for (const coin of coins) {
        if (coin.denom === asset) {
          console.log(coin.amount);
          // TODO - Decimals and Amount API formatting
          amount = Number(coin.amount).toFixed(3);
        }
      }
      return amount;
    })
    .catch(err => {
      throw `getSifchainBalances(): ${err}`;
    });
}
