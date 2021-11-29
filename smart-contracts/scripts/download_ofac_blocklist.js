const parser = require("./helpers/ofacParser");
const {print} = require("./helpers/utils");
const fs = require("fs");

async function main() {
    if (process.argv.length < 3) {
        print("h_red", "please specify a filename to store parsed list")
    }
    let ofac = await parser.getList();
    let msg = {
        addresses: ofac,
    }
    let msgJSON = JSON.stringify(msg)
    await new Promise((resolve, reject) => {
        fs.realpath(process.argv[2], (e, path) => {
            print("magenta", "Saving update msg to " + path)
            try {
                fs.writeFileSync(path, msgJSON)        
            } catch (err) {
                print("h_red", err.message)
                reject()
                return
            }
            print("magenta", "File saved.")
            resolve()
        })
    })
}

main()
    .catch((error) => {
        print("h_red", error.message)
    })
    .finally(() => process.exit(0));
