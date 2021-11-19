import * as hre from "hardhat"
import { EthereumAccounts, EthereumAddressAndKey, EthereumResults, ShellCommand } from "./devEnv"
import * as ChildProcess from "child_process"
import notifier from "node-notifier"

export class HardhatNodeRunner extends ShellCommand<EthereumResults> {
  private output: Promise<EthereumResults>
  private outputResolve: any
  constructor(
    readonly host = "localhost",
    readonly port = 8545,
    readonly nValidators = 1,
    readonly networkId = 1,
    readonly chainId = 1
  ) {
    super()
    this.output = new Promise<EthereumResults>((res, rej) => {
      this.outputResolve = res
    })
  }

  cmd(): [string, string[]] {
    return [
      "node_modules/.bin/hardhat",
      ["node", "--hostname", this.host, "--port", this.port.toString()],
    ]
  }

  override async run(): Promise<void> {
    const [c, args] = this.cmd()
    const childInfo = ChildProcess.spawn(c, args, {
      stdio: "inherit",
    })
    let ethereumAccounts = signerArrayToEthereumAccounts(defaultHardhatAccounts, this.nValidators)
    this.outputResolve({
      process: childInfo,
      accounts: ethereumAccounts,
      httpHost: this.host,
      httpPort: this.port,
      chainId: hre.network.config.chainId,
    })

    childInfo.on("exit", (code) => {
      notifier.notify({
        title: "HardHat Notice",
        message: `Hardhat has just exited with exit code: ${code}`,
      })
    })

    return
  }

  override async results(): Promise<EthereumResults> {
    return this.output
  }
}

// Hardhat doesn't provide a way to get the private keys of its default accounts, so
// just hardcode them for now.
const defaultHardhatAccounts: EthereumAddressAndKey[] = [
  {
    address: "0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266",
    privateKey: "0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80",
  },
  {
    address: "0x70997970c51812dc3a010c7d01b50e0d17dc79c8",
    privateKey: "0x59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d",
  },
  {
    address: "0x3c44cdddb6a900fa2b585dd299e03d12fa4293bc",
    privateKey: "0x5de4111afa1a4b94908f83103eb1f1706367c2e68ca870fc3fb9a804cdab365a",
  },
  {
    address: "0x90f79bf6eb2c4f870365e785982e1f101e93b906",
    privateKey: "0x7c852118294e51e653712a81e05800f419141751be58f605c371e15141b007a6",
  },
  {
    address: "0x15d34aaf54267db7d7c367839aaf71a00a2c6a65",
    privateKey: "0x47e179ec197488593b187f80a00eb0da91f1b9d0b13f8733639f19c30a34926a",
  },
  {
    address: "0x9965507d1a55bcc2695c58ba16fb37d819b0a4dc",
    privateKey: "0x8b3a350cf5c34c9194ca85829a2df0ec3153be0318b5e2d3348e872092edffba",
  },
  {
    address: "0x976ea74026e726554db657fa54763abd0c3a0aa9",
    privateKey: "0x92db14e403b83dfe3df233f83dfa3a0d7096f21ca9b0d6d6b8d88b2b4ec1564e",
  },
  {
    address: "0x14dc79964da2c08b23698b3d3cc7ca32193d9955",
    privateKey: "0x4bbbf85ce3377467afe5d46f804f221813b2bb87f24d81f60f1fcdbf7cbf4356",
  },
  {
    address: "0x23618e81e3f5cdf7f54c3d65f7fbc0abf5b21e8f",
    privateKey: "0xdbda1821b80551c9d65939329250298aa3472ba22feea921c0cf5d620ea67b97",
  },
  {
    address: "0xa0ee7a142d267c1f36714e4a8f75612f20a79720",
    privateKey: "0x2a871d0798f97d79848a013d4936a73bf4cc922c825d33c1cf7073dff6d409c6",
  },
]

function signerArrayToEthereumAccounts(
  accounts: EthereumAddressAndKey[],
  nValidators: number
): EthereumAccounts {
  const [operator, owner, pauser, ...moreAccounts] = accounts
  const validators = moreAccounts.slice(0, nValidators)
  const available = moreAccounts.slice(nValidators)
  return {
    proxyAdmin: operator,
    operator,
    owner,
    pauser,
    validators: validators,
    available,
  }
}
