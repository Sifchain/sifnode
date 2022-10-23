import { filter } from "rxjs/operators"
import * as readline from "readline"
import { Observable, OperatorFunction, ReplaySubject } from "rxjs"

export function readableStreamToObservable<T>(input: NodeJS.ReadableStream): Observable<string> {
  const str = new ReplaySubject<string>()

  const rl = readline.createInterface(input, undefined)
  rl.on("line", (ln) => {
    str.next(ln)
  })
  rl.on("close", () => {
    str.complete()
  })

  return str
}

export function tailFileAsObservable<T>(filename: string): Observable<string> {
  const str = new ReplaySubject<string>()

  const Tail = require("tail").Tail

  const tailedFile = new Tail(filename, { fromBeginning: false })

  tailedFile.on("line", (ln: string) => {
    str.next(ln)
  })
  tailedFile.on("close", () => {
    str.complete()
  })

  return str
}

export const jsonParseSimple = (x: string) => {
  try {
    return JSON.parse(x)
  } catch (err) {
    console.error("Error parsing json:", x)
    return {}
  }
}
export const jsonStringifySimple = (x: any) => JSON.stringify(x)

export function isNotNullOrUndefined<T>(input: null | undefined | T): input is T {
  switch (input) {
    case undefined:
    case null:
      return false
    default:
      return true
  }
}
