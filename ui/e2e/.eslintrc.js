module.exports = {
  root: true,
  plugins: ["prettier", "jest-playwright"],
  extends: ["prettier", "plugin:jest-playwright/recommended"],
  parserOptions: {
    ecmaVersion: 2018,
    sourceType: "module",
  },
  env: {
    node: true,
    jest: true,
    browser: true,
  },
  globals: {
    page: true,
    browser: true,
    context: true,
  },
};
