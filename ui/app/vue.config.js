const path = require('path');

module.exports = {
  lintOnSave: false,
  configureWebpack: {
    resolve: {
      extensions: ['.js', '.vue', '.json'],
      alias: {
        core: path.resolve(__dirname,'./../core/lib/index.js'),
      }
    }
  }

}

