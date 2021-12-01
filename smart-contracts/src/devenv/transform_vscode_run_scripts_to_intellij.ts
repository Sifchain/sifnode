import * as fs from "fs"
import path from "path"
import hb from "handlebars"

function renderEbrelayerConfig(x: any, rootDir: string) {
  const templatePath = path.join(
    rootDir,
    "smart-contracts/src/devenv/templates",
    "ebrelayer.run.xml.hbs"
  )
  const templateContents = fs.readFileSync(templatePath, { encoding: "utf-8" })
  const template = hb.compile(templateContents)
  const templateOutput = template({ ...x, joinedArgs: x["args"].join(" ") })
  fs.writeFileSync(path.join(rootDir, ".run", "ebrelayer.run.xml"), templateOutput)
}

function renderSifnodedConfig(x: any, rootDir: string) {
  const templatePath = path.join(
    rootDir,
    "smart-contracts/src/devenv/templates",
    "sifnoded.run.xml.hbs"
  )
  const templateContents = fs.readFileSync(templatePath, { encoding: "utf-8" })
  const template = hb.compile(templateContents)
  const templateOutput = template({
    ...x,
    joinedArgs: x["args"].join(" "),
    ethPrivateKey: x["ETHEREUM_PRIVATE_KEY"],
  })
  fs.writeFileSync(path.resolve(rootDir, ".run", "sifnoded.run.xml"), templateOutput)
}

export function renderIntellijFiles(rootDir: string) {
  fs.mkdirSync(path.join(rootDir, ".run"), { recursive: true })
  const fileContents = fs.readFileSync(path.join(rootDir, ".vscode/launch.json"), {
    encoding: "utf-8",
  })
  const goodContents = fileContents.replace(/\$\{workspaceFolder\}\//g, "")
  const cjson = JSON.parse(goodContents)
  for (const x of cjson.configurations) {
    if (x.name.startsWith("Debug Relayer")) {
      renderEbrelayerConfig(x, rootDir)
    }
    if (x.name.startsWith("Debug Sifnoded")) {
      renderSifnodedConfig(x, rootDir)
    }
  }
}
