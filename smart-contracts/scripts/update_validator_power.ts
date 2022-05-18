import * as hardhat from "hardhat"
import { container } from "tsyringe"
import { HardhatRuntimeEnvironmentToken } from "../src/tsyringe/injectionTokens"
import { Valset, Valset__factory } from "../build"
import { SifchainAccountsPromise } from "../src/tsyringe/sifchainAccounts"

// Usage
//
// COSMOSBRIDGE=0x890 VALIDATORS=0x123,0x456 POWERS=20,20 npx hardhat run scripts/update_validator_power.ts
//
// Normally this is only used by siftool

async function main() {
  container.register(HardhatRuntimeEnvironmentToken, { useValue: hardhat })
  const ownerAccount = (await (await container.resolve(SifchainAccountsPromise)).accounts)
    .operatorAccount
  const cosmosBridge = Valset__factory.connect(process.env["COSMOSBRIDGE"]!!, ownerAccount)
  let validators = process.env["VALIDATORS"]!!.split(",")
  let powers = process.env["POWERS"]!!.split(",")
  await (cosmosBridge as Valset).updateValset(validators, powers)
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error)
    process.exit(1)
  })
