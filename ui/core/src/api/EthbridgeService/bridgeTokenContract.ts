import Web3 from "web3";
import { AbiItem } from "web3-utils";

const json = require("../../../../../smart-contracts/build/contracts/BridgeToken.json");

export async function getBridgeTokenContract(web3: Web3, address: string) {
  return new web3.eth.Contract(json.abi as AbiItem[], address);
}
