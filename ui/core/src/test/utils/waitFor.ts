import { sleep } from "./sleep";

type ToStringable = { toString: () => string };

export async function waitFor(
  getter: () => Promise<ToStringable>,
  expected: ToStringable,
  name: string,
) {
  console.log(
    `Starting wait: "${name}" for value to be ${expected.toString()}`,
  );
  let value: any;
  for (let i = 0; i < 100; i++) {
    await sleep(1000);
    value = await getter();
    if (value.toString() === expected.toString()) {
      return;
    }
  }
  throw new Error(
    `${value.toString()} never was ${expected.toString()} in wait: ${name}`,
  );
}
