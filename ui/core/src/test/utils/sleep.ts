export const sleep = (ms: number) =>
  new Promise((done) => setTimeout(done, ms));
