/**
 * Create files that contain only the abi for each contract
 * Expected usage (option 1): yarn peggy:generateAbi
 * Expected usage (option 2): make abi
 * 
 * When you execute one of the commands above, hardhat will compile all contracts
 * Then, this script will run and save each contract's abi to the folder 
 * smart-contracts/build/contracts/abi
 */

 const fs = require('fs');
 const path = require('path');
 
 // where to find built contract artifacts (they will be there before this script gets executed)
 const builtPath = './build/contracts';
 
 // where to save the abis: buildPath + abiPath. Example of the final path: './build/contracts/abi'
 const abiPath = '/abi';
 
 async function exec() {
     // creates the dir if it doesn't already exist
     if (!fs.existsSync(builtPath + abiPath)){
         fs.mkdirSync(builtPath + abiPath);
     }
 
     // fetch all files from the build folder
     const files = await fs.promises.readdir(builtPath);
 
     for (let i = 0; i < files.length; i++) {
         const filename = path.join(builtPath, files[i]);
         const strippedFilename = files[i].replace('.json', '');
         
         try {
             // if it's a folder instead of a file, ignore it
             const stat = await fs.promises.lstat(filename);
             if(!stat.isFile()) {
                 continue;
             }
 
             console.log(`Processing ${strippedFilename}...`);
 
             // read what's in the file
             const data = fs.readFileSync(filename, 'utf8');
 
             // parse the JSON data
             const parsed = JSON.parse(data);
             
             // write the abi to a file
             fs.writeFileSync(`${builtPath}${abiPath}/${strippedFilename}.abi`, JSON.stringify(parsed.abi));
         } catch(e) {
             console.log({e});
         }
     }
 }
 
 exec();