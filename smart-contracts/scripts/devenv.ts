import {container} from "tsyringe";
import {HardhatNodeRunner} from "../src/devenv/devEnv"

async function main() {
    const node = container.resolve(HardhatNodeRunner)
    const process = node.run()
    const r = await node.results()
    console.log(`rsltis: ${JSON.stringify(r, undefined, 2)}`)
    await process
}

main()
    .then(() => process.exit(0))
    .catch((error) => {
        console.error(error);
        process.exit(1);
    });
