import {container} from "tsyringe";
import {HardhatNodeRunner} from "../src/devenv/hardhatNode";
import {GolangBuilder} from "../src/devenv/golangBuilder";

async function startHardhat() {
    const node = container.resolve(HardhatNodeRunner)
    const [process, resultsPromise] = node.go()
    const results = await resultsPromise
    console.log(`rsltis: ${JSON.stringify(results, undefined, 2)}`)
    return process
}

async function golangBuilder() {
    const node = container.resolve(GolangBuilder)
    const [process, resultsPromise] = node.go()
    const results = await resultsPromise
    console.log(`golangBuilder: ${JSON.stringify(results, undefined, 2)}`)
    return process
}

async function main() {
    await Promise.all([startHardhat(), golangBuilder()])
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
