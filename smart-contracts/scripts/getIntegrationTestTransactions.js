module.exports = async () => {
    require("dotenv").config();
    const Web3 = require("web3");
    const HDWalletProvider = require("@truffle/hdwallet-provider");
    try {
        let provider = new Web3.providers.HttpProvider(process.env.LOCAL_PROVIDER);
        const web3 = new Web3(provider);
        let logs = await web3.eth.getPastLogs({fromBlock: 0})
        return console.log("result:", JSON.stringify(logs, undefined, 0));
    } catch (error) {
        console.error({error})
    }
};
