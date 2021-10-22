import pEvent from "p-event";

async function fakeListener() {
  const proc = $`zx ./test/fakeRunner.mjs`;

  const event = pEvent.iterator(proc.stdout, "data", {
    resolutionEvents: ["finish"],
  });

  for await (let chunk of event) {
    if (chunk.includes("count: 1")) break;
  }
  console.log(`STOP`);
  await sleep(5000);
  console.log(`now run something different`);
}

fakeListener();
