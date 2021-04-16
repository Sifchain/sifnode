const path = require("path");
const rootPath = path.resolve(__dirname, "../app");

const sassAdditionalData = `
  @import "normalize-scss";
  @import "${rootPath}/src/scss/typography.scss";
  @import "${rootPath}/src/scss/variables.scss";
  @import "${rootPath}/src/scss/reset.scss";
  @import "${rootPath}/src/scss/utilities.scss";
  @import "${rootPath}/src/scss/mixins.scss";
`;

const sassLoaderRule = {
  test: /\.s[ac]ss$/i,
  use: [
    {
      loader: "sass-loader",
      options: {
        additionalData: sassAdditionalData,
      },
    },
  ],
};

module.exports = {
  stories: ["../app/**/*.stories.mdx", "../app/**/*.stories.@(js|jsx|ts|tsx)"],
  addons: [
    "@storybook/addon-links",
    "@storybook/addon-essentials",
    "storybook-addon-designs",
    {
      name: "@storybook/preset-scss",
      options: {
        sassLoaderOptions: {
          additionalData: sassAdditionalData,
        },
      },
    },
  ],
  webpackFinal: async (config, { configType }) => {
    config.module.rules.push({
      test: /\.scss$/,
      use: ["style-loader", "css-loader", "sass-loader"],
      include: path.resolve(__dirname, "../"),
    });
    config.module.rules.push(sassLoaderRule);
    return config;
  },
};
