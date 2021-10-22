async function fakeRunner() {
  let count = 0;

  while (true) {
    console.log(`count: ${count}`);
    await sleep(1000);
    count++;
  }
}

fakeRunner();
