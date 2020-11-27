const Artifactor = require("@truffle/artifactor");
const artifactor = new Artifactor(__dirname + "/../build-deploy/");
const fs = require('fs');

/*
 *
 *  This script will save contracts in the build directory to another directory without all of the nonsense like file paths
 * which are system dependant. This way, we can keep track of all needed fields like smart contract addresses without
 * having to keep build folder in git.
 * 
 *  The framework we are using, truffle artifactor, can detect changes in the address field for certain networks, so if you
 * make a change to one contract, i.e. deploy it again to the same network, it will just change that address field on that 
 * network and preserve all other network address information.
 *
 */

// Read all files in from a directory, call 2nd cb passed in with each file to process it
function readFiles(dirname, onFileContent, onError) {
  fs.readdir(dirname, function(err, filenames) {
    if (err) {
      onError(err);
      return;
    }
    filenames.forEach(function(filename) {
      fs.readFile(dirname + filename, 'utf-8', function(err, content) {
        if (err) {
          onError(filename, err);
          return;
        }
        onFileContent(filename, content);
      });
    });
  });
}

// See truffle-schema for more info: https://github.com/trufflesuite/truffle/tree/develop/packages/contract-schema
function handleFileContents(filename, content) {
    content = JSON.parse(content)
    if (!content.networks) {
        console.log("No network config found for: ", filename)
        return
    }
    const networkArray = Object.keys(content.networks)
    for (let i = 0; i < networkArray.length; i++) {
        const networkName = networkArray[i];
        const contractData = {
            contractName: content.contractName,// + networkArray[i],
            abi: content.abi,
            compiler: content.compiler,
            bytecode: content.bytecode,
            deployedBytecode: content.deployedBytecode,
            address: content.networks[networkName].address,
            transactionHash: content.networks[networkName].transactionHash,
            networks: {
                [networkName]: content.networks[networkName]
            }
        };
        artifactor.save(contractData);
        console.log("network: " + networkName + " filename: ", filename);
    }
}

function handleError(filename, error) {
    console.log("Error reading file: " + filename + " because " + error)
}

readFiles("build/contracts/", handleFileContents, handleError)
