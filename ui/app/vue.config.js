module.exports = {
  lintOnSave: false,
  
  css: {
    extract: {
      filename: "app.css"
    },
    loaderOptions: {
      sass: {
        additionalData: `
          @import "normalize-scss";
          @import "@/scss/variables.scss";
          @import "@/scss/mixins.scss";
        `
      }
    }
  }
}