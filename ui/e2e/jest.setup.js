const mkdirp = require("mkdirp");

global.it = async (name, func) => {
  return await test(name, async () => {
    try {
      await func();
    } catch (e) {
      const date = new Date();
      const year = date.getFullYear();
      const month = date.getUTCMonth() + 1;
      const dateOfMonth = date.getUTCDate();
      const hour = date.getUTCHours();
      const minute = date.getUTCMinutes();
      const sec = date.getUTCSeconds();
      const dateString = `${year}-${month}-${dateOfMonth}-${hour}-${minute}-${sec}`;

      const errorScreenshotPath = `screenshots/${browserName}-${dateString}-${name.replace(
        / /g,
        "_",
      )}`;

      await mkdirp("screenshots");

      const pages = await context.pages();
      await pages.forEach(async (page) => {
        const title = await page.title();
        await page.screenshot({
          path: `${errorScreenshotPath}_${title}.png`,
        });
      });

      throw e;
    }
  });
};
