const crypto = require("crypto");
const { network } = require("hardhat");

/**
 * @dev Generates a well-formed Denom
 * @dev The return value will look like this: sif789de8f7997bd47c4a0928a001e916b5c68f1f33fef33d6588b868b93b6dcde6
 * @dev this function expects an object with the following properties
 * @param {Number} networkDescriptor : what is the token's current network? Use 1 for Ethereum mainnet
 * @param {String} tokenAddress : the address of this token in its current network
 * @param {Boolean} isERC20 : is this an EVM token (true), or an IBC token (false)?
 * @returns {String} the final denom
 */
function generateDenom({ networkDescriptor, tokenAddress, isERC20 }) {

  if (isERC20) {
    if (networkDescriptor < 0 || networkDescriptor > 9999) {
      throw("invalid ERC20 Network Descriptor")
    }
    return `sifBridge${(networkDescriptor).toString().padStart(4, '0')}${tokenAddress.toLowerCase()}`
  } else {
    const fullString = `${networkDescriptor}/${tokenAddress.toLowerCase()}`;
    const hash = crypto.createHash("sha256").update(fullString).digest("hex");
    return `ibc/${hash}`;
  }
}

const ROWAN_DENOM = generateDenom({
  networkDescriptor: 1,
  tokenAddress: "0xF44bD7e809b9EFc5328e8AfCe949fE9E2E6D45dF",
  isERC20: true,
});

const ETHER_DENOM = generateDenom({
  networkDescriptor: 1,
  tokenAddress: "0x0000000000000000000000000000000000000000",
  isERC20: true, // it's not, be we'll treat it as if it was
});

const DENOM_1 = generateDenom({
  networkDescriptor: 1,
  tokenAddress: "0xB8c77482e45F1F44dE1745F52C74426C631bDD52",
  isERC20: true,
});

const DENOM_2 = generateDenom({
  networkDescriptor: 1,
  tokenAddress: "0xdac17f958d2ee523a2206206994597c13d831ec7",
  isERC20: true,
});

const DENOM_3 = generateDenom({
  networkDescriptor: 1,
  tokenAddress: "0x2b591e99afe9f32eaa6214f7b7629768c40eeb39",
  isERC20: true,
});

const DENOM_4 = generateDenom({
  networkDescriptor: 1,
  tokenAddress: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48",
  isERC20: true,
});

const IBC_DENOM = generateDenom({
  networkDescriptor: 1,
  tokenAddress: "0x0000000000000000000000000000000000000000",
  isERC20: false,
});

module.exports = {
  generateDenom,
  ROWAN_DENOM,
  ETHER_DENOM,
  DENOM_1,
  DENOM_2,
  DENOM_3,
  DENOM_4,
  IBC_DENOM,
};
