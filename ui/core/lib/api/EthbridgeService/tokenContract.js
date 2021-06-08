"use strict";
var __awaiter = (this && this.__awaiter) || function (thisArg, _arguments, P, generator) {
    function adopt(value) { return value instanceof P ? value : new P(function (resolve) { resolve(value); }); }
    return new (P || (P = Promise))(function (resolve, reject) {
        function fulfilled(value) { try { step(generator.next(value)); } catch (e) { reject(e); } }
        function rejected(value) { try { step(generator["throw"](value)); } catch (e) { reject(e); } }
        function step(result) { result.done ? resolve(result.value) : adopt(result.value).then(fulfilled, rejected); }
        step((generator = generator.apply(thisArg, _arguments || [])).next());
    });
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.getTokenContract = void 0;
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
function getTokenContract(web3, address) {
    return __awaiter(this, void 0, void 0, function* () {
        return new web3.eth.Contract(abi, address);
    });
}
exports.getTokenContract = getTokenContract;
//# sourceMappingURL=tokenContract.js.map