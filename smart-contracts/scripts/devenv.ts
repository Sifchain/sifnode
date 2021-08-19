import {container, registry, singleton} from "tsyringe";
import {HardhatNodeRunner} from "../src/devenv/hardhatNode";
import {GolangBuilder, GolangResultsPromise} from "../src/devenv/golangBuilder";
import {SifnodedRunner} from "../src/devenv/sifnoded";


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
    let golangResultsPromise = new GolangResultsPromise(resultsPromise);
    // container.register(GolangResultsPromise, {useValue: golangResultsPromise})
    // const sifnodeTask = sifnodedBuilder(golangResultsPromise)
    const results = await resultsPromise
    console.log(`golangBuilder: ${JSON.stringify(results, undefined, 2)}`)
    return process
}

async function sifnodedBuilder(golangResults: GolangResultsPromise) {
    const node = container.resolve(SifnodedRunner)
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
