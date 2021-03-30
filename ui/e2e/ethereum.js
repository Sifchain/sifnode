import { getEtheriumBalance } from "./../core/src/api/EthereumService/utils/ethereumUtils";
import { formatNumber } from "./../app/src/components/shared/utils.ts";

import Web3 from "web3";
const web3 = new Web3("http://localhost:7545");

export async function getEthBalance(address) {
  const balance = await getEtheriumBalance(web3, address);
  return formatNumber(balance.amount);
}
