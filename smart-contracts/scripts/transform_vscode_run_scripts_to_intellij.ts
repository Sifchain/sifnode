import * as fs from "fs";
import path from "path";
import hb from "handlebars";

const fileContents = fs.readFileSync(__dirname + "/../../.vscode/launch.json", {
  encoding: "utf-8",
});
const goodContents = fileContents.replace(/\$\{workspaceFolder\}\//g, "");
const cjson = JSON.parse(goodContents);
for (const x of cjson.configurations) {
  if (x.name.startsWith("Debug Relayer")) {
    console.log(`cj is: ${JSON.stringify(x)}`);
    renderEbrelayerConfig(x);
  }
  if (x.name.startsWith("Debug Sifnoded")) {
    console.log(`cj is: ${JSON.stringify(x)}`);
    renderSifnodedConfig(x);
  }
}

function renderEbrelayerConfig(x: any) {
  const templatePath = path.resolve(__dirname, "../src/devenv/templates", "ebrelayer.run.xml.hbs");
  const templateContents = fs.readFileSync(templatePath, { encoding: "utf-8" });
  const template = hb.compile(templateContents);
  const templateOutput = template({ ...x, joinedArgs: x["args"].join(" ") });
  fs.writeFileSync(path.resolve(__dirname, "../../.run", "ebrelayer.run.xml"), templateOutput);
}

function renderSifnodedConfig(x: any) {
  const templatePath = path.resolve(__dirname, "../src/devenv/templates", "sifnoded.run.xml.hbs");
  const templateContents = fs.readFileSync(templatePath, { encoding: "utf-8" });
  const template = hb.compile(templateContents);
  const templateOutput = template({
    ...x,
    joinedArgs: x["args"].join(" "),
    ethPrivateKey: x["ETHEREUM_PRIVATE_KEY"],
  });
  fs.writeFileSync(path.resolve(__dirname, "../../.run", "sifnoded.run.xml"), templateOutput);
}
