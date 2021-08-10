/**
 * Create files that contain only the abi for each contract
 * Expected usage (option 1): yarn peggy:generateAbi
 * Expected usage (option 2): make abi
 * 
 * When you execute one of the commands above, hardhat will compile all contracts
 * Then, this script will run and save each contract's abi to the folder 
 * cmd/ebrelayer/contract/generated/abi
 * And save each contract's binary to the folder 
 * cmd/ebrelayer/contract/generated/bin
 * And save the generated go files to the folder 
 * cmd/ebrelayer/contract/generated/bindings
 */

 const fs = require('fs');
 const path = require('path');
 const util = require('util');
 const exec = util.promisify(require('child_process').exec);
 
 // where to fetch artifacts from (hardhat should save compiled contracts in this folder)
 const HARDHAT_ARTIFACTS_DIR = './artifacts/contracts';
 
 // where to save the generated files (BASE_TARGET_DIR/TARGET_XXX_FOLDER):
 const BASE_TARGET_DIR = '../cmd/ebrelayer/contract/generated'
 
 // where to save the ABI files (BASE_TARGET_DIR/TARGET_ABI_DIR):
 const TARGET_ABI_DIR = 'abi';
 
 // where to save the BIN files (BASE_TARGET_DIR/TARGET_BIN_DIR):
 const TARGET_BIN_DIR = 'bin';
 
 // where to save the GO files (BASE_TARGET_DIR/TARGET_GO_DIR):
 const TARGET_GO_DIR = 'bindings';
 
 async function main() {
   // creates the target directory if it doesn't already exist
   createDirectories();
 
   // get only the files that we care about
   const files = getArtifacts();
 
   // For each file...
   for (let i = 0; i < files.length; i++) {
     // get the name of the file without its path
     const strippedFilename = files[i].split('/').slice(-1)[0].replace('.json', '');
     
     // get the normalized path
     const internalPath = path.dirname(files[i]).split('artifacts/contracts/')[1].replace('.sol', '').split('/')[0];
     
     console.log(`Processing ${internalPath}/${strippedFilename}...`);
 
     // read what's in the file
     const data = fs.readFileSync(files[i], 'utf8');
 
     // parse the JSON data
     const parsed = JSON.parse(data);
 
     // write the abi to a file
     const targetAbiDirectory = `${BASE_TARGET_DIR}/${TARGET_ABI_DIR}/${internalPath}`;
     createDir(targetAbiDirectory);
     const targetAbiFileName = `${targetAbiDirectory}/${strippedFilename}.abi`;
     fs.writeFileSync(targetAbiFileName, JSON.stringify(parsed.abi));
 
     // write the binary data to a file
     const targetBinDirectory = `${BASE_TARGET_DIR}/${TARGET_BIN_DIR}/${internalPath}`;
     createDir(targetBinDirectory);
     const targetBinFileName = `${targetBinDirectory}/${strippedFilename}.bin`;
     fs.writeFileSync(targetBinFileName, JSON.stringify(parsed.bytecode));
 
     // create go bindings for this contract
     const targetGoDirectory = `${BASE_TARGET_DIR}/${TARGET_GO_DIR}/${internalPath}`;
     const targetGoFileName = `${targetGoDirectory}/${strippedFilename}.go`;
     if (fs.existsSync(targetGoFileName)) {
      createDir(targetGoDirectory);
      await exec(`abigen --abi ${targetAbiFileName} --pkg ${internalPath} --type ${internalPath} --out ${targetGoFileName}`);
     }
   }
 
   printSuccess();
 }
 
 function createDirectories() {
   createDir(BASE_TARGET_DIR);
   createDir(`${BASE_TARGET_DIR}/${TARGET_ABI_DIR}`);
   createDir(`${BASE_TARGET_DIR}/${TARGET_BIN_DIR}`);
   createDir(`${BASE_TARGET_DIR}/${TARGET_GO_DIR}`);
 }
 
 function createDir(directory) {
   if (!fs.existsSync(directory)) {
     fs.mkdirSync(directory);
   }
 }
 
 // Get all files that end with .json and not .dbg.json
 function getArtifacts() {
   const allContractsArtifacts = [];
   getFilesRecursively(HARDHAT_ARTIFACTS_DIR, allContractsArtifacts);
 
   const filtered = allContractsArtifacts.filter((file) => {
     return file.endsWith('.json') && !file.endsWith('.dbg.json');
   });
 
   return filtered;
 }
 
 // Traverse through a directory and get all files recursively
 // Save the result in `targetList`
 function getFilesRecursively(directory, targetList) {
   const filesInDirectory = fs.readdirSync(directory);
 
   for (const file of filesInDirectory) {
     const absolute = path.join(directory, file);
 
     if (fs.statSync(absolute).isDirectory()) {
       getFilesRecursively(absolute, targetList);
     } else {
       targetList.push(absolute);
     }
   }
 };
 
 function printSuccess() {
   console.log(`-> Generated ABI files in ${BASE_TARGET_DIR}/${TARGET_ABI_DIR}/`);
   console.log(`-> Generated BIN files in ${BASE_TARGET_DIR}/${TARGET_BIN_DIR}/`);
   console.log(`-> Generated GO files in ${BASE_TARGET_DIR}/${TARGET_GO_DIR}/`);
 }
 
 main()
   .catch((error) => {
     console.error({ error });
   })
   .finally(() => process.exit(0))