import { ETH } from "../../constants";
import { Asset, Balance, Token } from "../../entities";
import Web3 from "web3";
import { AbiItem } from "web3-utils";

type Address = string;
type Balances = Balance[];

function isToken(value?: Asset | Token): value is Token {
  return !value || Object.keys(value).includes("address");
}

const generalTokenAbi: AbiItem[] = [
  // balanceOf
  {
    constant: true,
    inputs: [{ name: "_owner", type: "address" }],
    name: "balanceOf",
    outputs: [{ name: "balance", type: "uint256" }],
    type: "function",
  },
  // decimals
  {
    constant: true,
    inputs: [],
    name: "decimals",
    outputs: [{ name: "", type: "uint8" }],
    type: "function",
  },
];

async function getTokenBalance(web3: Web3, address: Address, asset: Token) {
  const contract = new web3.eth.Contract(generalTokenAbi, asset.address);
  const tokenBalance = await contract.methods.balanceOf(address).call();
  return Balance.create(asset, tokenBalance);
}

async function getEtheriumBalance(web3: Web3, address: Address) {
  const ethBalance = await web3.eth.getBalance(address);
  return Balance.create(ETH, ethBalance);
}

export class EthProvider {
  constructor(
    private address: Address,
    private web3: Web3,
    private supportedTokens: Token[]
  ) {}

  getAddress(): Address {
    return this.address;
  }

  async getBalance(
    address?: Address,
    asset?: Asset | Token
  ): Promise<Balances> {
    const addr = address || this.getAddress();

    if (asset) {
      if (!isToken(asset)) {
        // Asset must be eth
        const ethBalance = await getEtheriumBalance(this.web3, addr);
        return [ethBalance];
      }

      // Asset must be ERC-20
      const tokenBalance = await getTokenBalance(this.web3, addr, asset);
      return [tokenBalance];
    }

    // No address no asset get everything
    const balances = await Promise.all([
      getEtheriumBalance(this.web3, addr),
      ...this.supportedTokens.map((token: Token) => {
        return getTokenBalance(this.web3, addr, token);
      }),
    ]);

    return balances;
  }
}
