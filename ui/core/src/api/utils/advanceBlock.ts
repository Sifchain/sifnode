// No TS defs yet provided https://github.com/OpenZeppelin/openzeppelin-test-helpers/pull/141
const { time } = require("@openzeppelin/test-helpers");

export async function advanceBlock(count: number) {
  for (let i = 0; i < count; i++) {
    await time.advanceBlock();
    // Need to provide time between advances or it doesn't work
    await new Promise((resolve) => setTimeout(resolve, 500));
  }
}
