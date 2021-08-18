import {container} from "tsyringe";
import {HardhatNodeRunner} from "../src/devenv/hardhatNode";

async function main() {
    const node = container.resolve(HardhatNodeRunner)
    const [process, resultsPromise] = node.go()
    const results = await resultsPromise
    console.log(`rsltis: ${JSON.stringify(results, undefined, 2)}`)
    await process
}

main()
    .then(() => process.exit(0))
    .catch((error) => {
        if (typeof error == "number")
            process.exit(error)
        else {
            console.error(error);
            process.exit(1)
        }
    });
