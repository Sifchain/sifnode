module.exports = {
  preset: "jest-playwright-preset",
  testTimeout: 1000000,
  collectCoverage: true,
  bail: true,
  setupFilesAfterEnv: ["expect-playwright", "./setup.js", "./jest.setup.js"],
};
