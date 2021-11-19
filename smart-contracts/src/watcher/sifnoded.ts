import { filter, map } from "rxjs/operators"
import * as fs from "fs"
import * as readline from "readline"
import { Observable, ReplaySubject } from "rxjs"
import { jsonParseSimple, readableStreamToObservable } from "./utilities"

export interface SifnodedInfoEvent {
  kind: "SifnodedInfoEvent"
  data: object
}

export interface SifnodedError {
  kind: "SifnodedError"
  data: object
}

export interface SifnodedPeggyEvent {
  kind: "SifnodedPeggyEvent"
  data: {
    kind: string
  }
}

export type SifnodedEvent = SifnodedInfoEvent | SifnodedError | SifnodedPeggyEvent

export function isSifnodedEvent(x: object): x is SifnodedEvent {
  switch ((x as SifnodedEvent).kind) {
    case "SifnodedError":
    case "SifnodedInfoEvent":
    case "SifnodedPeggyEvent":
      return true
    default:
      return false
  }
}

export function isNotSifnodedEvent(x: object): x is SifnodedEvent {
  return !isSifnodedEvent(x)
}

export function toSifnodedEvent(x: any): SifnodedEvent | undefined {
  if (x.message === "peggytest") return { kind: "SifnodedPeggyEvent", data: x }
  else if (x.level === "info") return { kind: "SifnodedInfoEvent", data: x }
  else if (x.level === "error") return { kind: "SifnodedError", data: x }
  else {
    return undefined
  }
}
