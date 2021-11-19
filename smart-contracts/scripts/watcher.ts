import { filter, map } from "rxjs/operators"
import * as fs from "fs"
import * as readline from "readline"
import { Observable, ReplaySubject } from "rxjs"
import {
  isNotNullOrUndefined,
  jsonParseSimple,
  readableStreamToObservable,
  tailFileAsObservable,
} from "../src/watcher/utilities"
import { lastValueFrom } from "rxjs"
import { EvmStateTransition, toEvmRelayerEvent } from "../src/watcher/ebrelayer"
import { sifwatch } from "../src/watcher/watcher"

async function main() {
  const evmRelayerEvents = sifwatch("/tmp/sifnode/evmrelayer.log")

  evmRelayerEvents.subscribe({
    next: (x) => {
      console.log(x)
    },
    error: (e) => console.log("goterror: ", e),
    complete: () => console.log("alldone"),
  })

  const lv = await lastValueFrom(evmRelayerEvents)
  console.log("exitingwatcher")
}

main()
  .catch((error) => {
    console.error(error)
  })
  .finally(() => process.exit(0))
