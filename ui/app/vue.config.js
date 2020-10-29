module.exports = {
  lintOnSave: false,
  devServer: {
    overlay: {
      warnings: true,
      errors: true,
    },
  },
  css: {
    extract: {
      filename: "app.css"
    },
    loaderOptions: {
      sass: {
        additionalData: `
          @import "normalize-scss";
          @import "@/scss/reset.scss";
          @import "@/scss/typography.scss";
          @import "@/scss/variables.scss";
          @import "@/scss/utilities.scss";
          @import "@/scss/mixins.scss";
        `
      }
    }
  }
};
