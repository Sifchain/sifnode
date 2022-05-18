import { lastValueFrom } from "rxjs"
import * as rxops from "rxjs/operators"
import { defaultSifwatchLogs, sifwatch } from "../src/watcher/watcher"
import * as hardhat from "hardhat"
import { container } from "tsyringe"
import { HardhatRuntimeEnvironmentToken } from "../src/tsyringe/injectionTokens"
import { setupDeployment } from "../src/hardhatFunctions"
import { readDevEnvObj } from "../src/tsyringe/devenvUtilities"
import { BridgeBank__factory } from "../build"

async function main() {
  container.register(HardhatRuntimeEnvironmentToken, { useValue: hardhat })

  await setupDeployment(container)

  const devenv = await readDevEnvObj("./environment.json")

  const bridgeBank = await BridgeBank__factory.connect(
    devenv.contractResults!!.contractAddresses.bridgeBank,
    hardhat.ethers.provider
  )

  const evmRelayerEvents = sifwatch(defaultSifwatchLogs(), hardhat, bridgeBank)

  evmRelayerEvents
    .pipe(
      rxops.filter((x) => {
        switch (x.kind) {
          case "SifHeartbeat":
          case "SifnodedInfoEvent":
            return false
          default:
            return true
        }
      })
    )
    .subscribe({
      next: (x) => {
        console.log(JSON.stringify(x))
      },
      error: (e) => console.log("Terminated with error: ", e),
      complete: () => console.log("Normal exit"),
    })

  const lv = await lastValueFrom(evmRelayerEvents)
}

main()
  .catch((error) => {
    console.error(error)
  })
  .finally(() => process.exit(0))
