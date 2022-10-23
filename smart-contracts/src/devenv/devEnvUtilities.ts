import events from "events"
import { ReplaySubject } from "rxjs"

export class ErrorEvent {
  constructor(readonly errorObject: any) {}
}

export function eventEmitterToObservable(
  eventEmitter: events.EventEmitter,
  sourceName: string = "no name given"
) {
  const subject = new ReplaySubject<"exit" | ErrorEvent>(1)
  eventEmitter.on("error", (e) => {
    subject.error(new ErrorEvent(e))
  })
  eventEmitter.on("exit", (e) => {
    console.log("in eventEmitter")
    switch (e) {
      case 0:
        subject.next("exit")
        subject.complete()
        break
      default:
        subject.error(new ErrorEvent(e))
        break
    }
  })
  return subject.asObservable()
}

export async function sleep(milliseconds: number) {
  await new Promise((resolve) => setTimeout(resolve, milliseconds))
}

export const sleepForever = Promise.race([])
