import Web3 from "web3";
import { AbiItem } from "web3-utils";

// We should add other ABIs as they are required here
// standard ERC-20 approveFn see https://docs.openzeppelin.com/contracts/2.x/api/token/erc20#IERC20-approve-address-uint256-
// We also may want to create library of these standard ABI calls as well as custom ones we write
const approveFn = {
  constant: false,
  inputs: [
    {
      internalType: "address",
      name: "spender",
      type: "address",
    },
    {
      internalType: "uint256",
      name: "amount",
      type: "uint256",
    },
  ],
  name: "approve",
  outputs: [
    {
      internalType: "bool",
      name: "",
      type: "bool",
    },
  ],
  payable: false,
  stateMutability: "nonpayable",
  type: "function",
};

// https://ethereumdev.io/abi-for-erc20-contract-on-ethereum/
const allowanceFn = {
  constant: true,
  inputs: [
    {
      name: "_owner",
      type: "address",
    },
    {
      name: "_spender",
      type: "address",
    },
  ],
  name: "allowance",
  outputs: [
    {
      name: "",
      type: "uint256",
    },
  ],
  payable: false,
  stateMutability: "view",
  type: "function",
};

// todo allowance function
const abi = [approveFn, allowanceFn];

export async function getTokenContract(web3: Web3, address: string) {
  return new web3.eth.Contract(abi as AbiItem[], address);
}
