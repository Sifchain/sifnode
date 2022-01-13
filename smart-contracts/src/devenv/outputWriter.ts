import { GolangResults } from "./golangBuilder"
import { SifnodedResults } from "./sifnoded"
import { SmartContractDeployResult } from "./smartcontractDeployer"
import { EthereumResults } from "./devEnv"
import path from "path"
import fs from "fs"
import hb from "handlebars"
import { renderIntellijFiles } from "./transform_vscode_run_scripts_to_intellij"
interface EnvOutput {
  Computed: {
    BASEDIR: string
    CHAINDIR?: string
  }
  Dev: DevEnvObject
  Env?: string
}
export interface DevEnvObject {
  ethResults?: EthereumResults
  goResults?: GolangResults
  sifResults?: SifnodedResults
  contractResults?: SmartContractDeployResult
}

/**
 * Takes a Handle Bars Template file, a object of arguments to replace in the template file, and
 * then compiles and saves the rendered document to the save location
 * @param templateLocation Location of the handlebars template *.hbs
 * @param saveLocation Where the rendered document should be saved
 * @param args The variables to be replaced in the template during render
 */
function RenderTemplateToFile(
  templateLocation: string,
  saveLocation: string,
  args: unknown
): string {
  hb.registerHelper(
    "subString",
    function (inputString: string, startIndex: number, endIndex?: number) {
      /**
       * This if statement is needed because handlebar passes in a hash as it's
       * last param. This causes issue because endIndex is optional
       */
      if (!endIndex || typeof endIndex != "number") {
        endIndex = undefined
      }
      let trimmedString: string = inputString.substring(startIndex, endIndex)
      return new hb.SafeString(trimmedString)
    }
  )

  const template = fs.readFileSync(templateLocation, "utf8")
  const compiledTemplate = hb.compile(template)
  const renderedTemplate = compiledTemplate(args)
  // Make sure the .vscode directory exists
  const dirPath = path.dirname(saveLocation)
  fs.mkdirSync(dirPath, { recursive: true })
  // Save the file in the .vscode directory
  fs.writeFileSync(saveLocation, renderedTemplate)
  return renderedTemplate
}

export function EnvJSONWriter(args: DevEnvObject) {
  const baseDir = path.resolve(__dirname, "../../..")
  const output: EnvOutput = {
    Computed: {
      BASEDIR: baseDir,
    },
    Dev: args,
  }
  if (args.sifResults != undefined) {
    const sif = args.sifResults
    const val = sif.validatorValues[0]
    output.Computed.CHAINDIR = path.resolve(
      "/tmp/sifnodedNetwork/validators",
      val.chain_id,
      val.moniker
    )
  }
  try {
    RenderTemplateToFile(
      path.resolve(__dirname, "templates", "env.hbs"),
      path.resolve(__dirname, "../../", ".env"),
      output
    )
    fs.writeFileSync(path.resolve(__dirname, "../../", "environment.json"), JSON.stringify(args))
    const envJSON = RenderTemplateToFile(
      path.resolve(__dirname, "templates", "env.json.hbs"),
      path.resolve(__dirname, "../../", "env.json"),
      output
    )
    output.Env = envJSON
    RenderTemplateToFile(
      path.resolve(__dirname, "templates", "launch.json.hbs"),
      path.resolve(__dirname, "../../../", ".vscode", "launch.json"),
      output
    )
    renderIntellijFiles(path.join(__dirname, "../../.."))
    console.log("Wrote environment and JSON values to disk. PATH: ", path.resolve(__dirname))
  } catch (error) {
    console.error("Failed to write environment/json values to disk, ERROR: ", error)
  }
}
