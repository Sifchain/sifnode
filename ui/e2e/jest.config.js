module.exports = {
  preset: "jest-playwright-preset",
  testTimeout: 1000000,
  collectCoverage: true,
  // below files use page.evaluate method which gets in trouble with jest coverage reporting. Workaround: ignore that files
  // https://github.com/facebook/jest/issues/7962
  coveragePathIgnorePatterns: ["e2e/utils.js"],
  bail: true,
  setupFilesAfterEnv: ["expect-playwright", "./setup.js"],
  // setupFilesAfterEnv: ["expect-playwright", "./setup.js", "./jest.setup.js"],
};
